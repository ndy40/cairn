package updater

import (
	"archive/tar"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const (
	defaultAPIURL = "https://api.github.com/repos/ndy40/cairn/releases/latest"
	httpTimeout   = 8 * time.Second
)

// ErrChecksumMismatch is returned when the downloaded file's SHA256 does not
// match the expected value from checksums.txt.
var ErrChecksumMismatch = errors.New("checksum mismatch")

// ErrPermission is returned when the install directory is not writable by the
// current user.
var ErrPermission = errors.New("permission denied")

// Package-level vars allow tests to inject a custom server URL and HTTP client.
var (
	releaseAPIURL    = defaultAPIURL
	updateClient     = &http.Client{Timeout: httpTimeout}
	extensionDirFunc = extensionDir // overridable in tests
)

// releaseInfo holds the fields we need from the GitHub Releases API response.
type releaseInfo struct {
	TagName string `json:"tag_name"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

// CheckLatestVersion queries the releases API for the latest published version.
// It returns the latest version tag, whether an update is available, and any
// error. A currentVersion of "dev" is always considered out of date.
func CheckLatestVersion(currentVersion string) (latest string, available bool, err error) {
	info, err := fetchReleaseInfo(currentVersion)
	if err != nil {
		return "", false, err
	}
	available = currentVersion == "dev" || currentVersion != info.TagName
	return info.TagName, available, nil
}

// UpdateBinary downloads the latest cairn binary for the current platform,
// verifies its SHA256 checksum, backs up the existing binary, and replaces it
// atomically. The existing binary is restored on failure.
func UpdateBinary(currentVersion, latestVersion string) error {
	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("resolve executable path: %w", err)
	}
	exe, err = filepath.EvalSymlinks(exe)
	if err != nil {
		return fmt.Errorf("resolve symlinks: %w", err)
	}
	return updateBinaryAt(exe, currentVersion, latestVersion)
}

// updateBinaryAt is the internal, testable implementation of UpdateBinary.
func updateBinaryAt(targetPath, currentVersion, latestVersion string) error {
	if runtime.GOOS == "windows" {
		fmt.Println("cairn update: Windows in-process update is not supported; re-run the install script to update")
		return nil
	}

	installDir := filepath.Dir(targetPath)
	if !isWritable(installDir) {
		return ErrPermission
	}

	binaryName := fmt.Sprintf("cairn-%s-%s", runtime.GOOS, runtime.GOARCH)

	info, err := fetchReleaseInfo(currentVersion)
	if err != nil {
		return err
	}

	binaryURL, checksumURL := findAssetURLs(info, binaryName)
	if binaryURL == "" {
		return fmt.Errorf("no release asset found for %s", binaryName)
	}

	// Download checksums.txt to a temp file.
	checksumPath, err := downloadToTemp(installDir, checksumURL, ".cairn-checksums-*")
	if err != nil {
		return fmt.Errorf("download checksums: %w", err)
	}
	defer func() { _ = os.Remove(checksumPath) }()

	checksumContent, err := os.ReadFile(checksumPath)
	if err != nil {
		return fmt.Errorf("read checksums: %w", err)
	}
	expectedHash, err := parseChecksum(string(checksumContent), binaryName)
	if err != nil {
		return fmt.Errorf("find checksum for %s: %w", binaryName, err)
	}

	// Download binary to a temp file in the same directory as the target so
	// that os.Rename is guaranteed to be on the same filesystem.
	fmt.Printf("cairn: downloading %s...\n", binaryName)
	tmpPath, err := downloadToTemp(installDir, binaryURL, ".cairn-update-*")
	if err != nil {
		return fmt.Errorf("download binary: %w", err)
	}
	defer func() { _ = os.Remove(tmpPath) }()

	fmt.Println("cairn: verifying checksum...")
	if err := verifyChecksum(tmpPath, expectedHash); err != nil {
		return err
	}

	if err := os.Chmod(tmpPath, 0o755); err != nil {
		return fmt.Errorf("chmod binary: %w", err)
	}

	// Backup the existing binary before replacing it.
	bakPath := targetPath + ".bak"
	if err := copyFile(targetPath, bakPath); err != nil {
		return fmt.Errorf("backup existing binary: %w", err)
	}

	// Atomic rename — restores backup on failure.
	if err := os.Rename(tmpPath, targetPath); err != nil {
		_ = os.Rename(bakPath, targetPath)
		return fmt.Errorf("replace binary: %w", err)
	}
	_ = os.Remove(bakPath)

	fmt.Printf("cairn: updated to %s\n", latestVersion)
	return nil
}

// DetectExtension returns the platform-specific extension directory and whether
// it is currently installed (i.e., the directory exists).
func DetectExtension() (dir string, installed bool) {
	dir = extensionDirFunc()
	_, err := os.Stat(dir)
	return dir, err == nil
}

// CheckExtensionVersion reads the installed extension version from
// dir/version.txt and queries the releases API for the latest version. It
// returns both versions and whether an update is available.
func CheckExtensionVersion(dir string) (current, latest string, available bool, err error) {
	data, readErr := os.ReadFile(filepath.Join(dir, "version.txt"))
	if readErr != nil {
		current = "unknown"
	} else {
		current = strings.TrimSpace(string(data))
	}

	latest, _, err = CheckLatestVersion(current)
	if err != nil {
		return current, "", false, err
	}

	available = current == "unknown" || current != latest
	return current, latest, available, nil
}

// UpdateExtension downloads the latest Vicinae extension archive for the
// current release, verifies its checksum, extracts it to dir, and writes
// dir/version.txt.
func UpdateExtension(dir, latestVersion string) error {
	if runtime.GOOS == "windows" {
		fmt.Println("cairn update: Windows in-process update is not supported; re-run the install script to update")
		return nil
	}

	archiveName := fmt.Sprintf("vicinae-extension-%s.tar.gz", latestVersion)

	info, err := fetchReleaseInfo(latestVersion)
	if err != nil {
		return err
	}

	archiveURL, checksumURL := findAssetURLs(info, archiveName)
	if archiveURL == "" {
		return fmt.Errorf("no release asset found for %s", archiveName)
	}

	tmpDir, err := os.MkdirTemp("", "cairn-ext-update-*")
	if err != nil {
		return fmt.Errorf("create temp dir: %w", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	// Download and verify checksum if available.
	var expectedHash string
	if checksumURL != "" {
		checksumPath, err := downloadToTemp(tmpDir, checksumURL, ".cairn-checksums-*")
		if err != nil {
			return fmt.Errorf("download checksums: %w", err)
		}
		checksumContent, err := os.ReadFile(checksumPath)
		if err != nil {
			return fmt.Errorf("read checksums: %w", err)
		}
		expectedHash, err = parseChecksum(string(checksumContent), archiveName)
		if err != nil {
			return fmt.Errorf("find checksum for %s: %w", archiveName, err)
		}
	}

	fmt.Printf("cairn: downloading extension %s...\n", latestVersion)
	archivePath := filepath.Join(tmpDir, archiveName)
	if err := downloadFile(archiveURL, archivePath); err != nil {
		return fmt.Errorf("download extension archive: %w", err)
	}

	if expectedHash != "" {
		fmt.Println("cairn: verifying checksum...")
		if err := verifyChecksum(archivePath, expectedHash); err != nil {
			return err
		}
	}

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create extension directory: %w", err)
	}

	if err := extractTarGz(archivePath, dir); err != nil {
		return fmt.Errorf("extract extension: %w", err)
	}

	if err := os.WriteFile(filepath.Join(dir, "version.txt"), []byte(latestVersion+"\n"), 0o644); err != nil {
		return fmt.Errorf("write version.txt: %w", err)
	}

	fmt.Printf("cairn: extension updated to %s\n", latestVersion)
	return nil
}

// ── internal helpers ─────────────────────────────────────────────────────────

func fetchReleaseInfo(currentVersion string) (*releaseInfo, error) {
	req, err := http.NewRequest(http.MethodGet, releaseAPIURL, nil)
	if err != nil {
		return nil, fmt.Errorf("build release request: %w", err)
	}
	req.Header.Set("User-Agent", "cairn/"+currentVersion)
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := updateClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch release info: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetch release info: server returned %s", resp.Status)
	}

	var info releaseInfo
	if err := json.NewDecoder(io.LimitReader(resp.Body, 1<<20)).Decode(&info); err != nil {
		return nil, fmt.Errorf("parse release info: %w", err)
	}
	if info.TagName == "" {
		return nil, fmt.Errorf("release info missing tag_name")
	}
	return &info, nil
}

func findAssetURLs(info *releaseInfo, targetName string) (assetURL, checksumURL string) {
	for _, asset := range info.Assets {
		switch asset.Name {
		case targetName:
			assetURL = asset.BrowserDownloadURL
		case "checksums.txt":
			checksumURL = asset.BrowserDownloadURL
		}
	}
	return
}

// downloadToTemp downloads url to a new temp file in dir and returns its path.
func downloadToTemp(dir, url, pattern string) (string, error) {
	f, err := os.CreateTemp(dir, pattern)
	if err != nil {
		return "", err
	}
	path := f.Name()
	_ = f.Close()
	if err := downloadFile(url, path); err != nil {
		_ = os.Remove(path)
		return "", err
	}
	return path, nil
}

func downloadFile(url, destPath string) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "cairn/updater")

	resp, err := updateClient.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download %s: server returned %s", url, resp.Status)
	}

	f, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()

	_, err = io.Copy(f, resp.Body)
	return err
}

func verifyChecksum(filePath, expectedHash string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return err
	}
	got := hex.EncodeToString(h.Sum(nil))
	if got != strings.ToLower(expectedHash) {
		return ErrChecksumMismatch
	}
	return nil
}

// parseChecksum parses a "sha256sum"-style file and returns the hex hash for
// the named file. Each line has the format: "<hash>  <filename>".
func parseChecksum(content, filename string) (string, error) {
	for _, line := range strings.Split(content, "\n") {
		fields := strings.Fields(line)
		if len(fields) >= 2 && fields[1] == filename {
			return fields[0], nil
		}
	}
	return "", fmt.Errorf("%s not found in checksums", filename)
}

func isWritable(dir string) bool {
	f, err := os.CreateTemp(dir, ".cairn-perm-check-*")
	if err != nil {
		return false
	}
	_ = f.Close()
	_ = os.Remove(f.Name())
	return true
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func() { _ = in.Close() }()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() { _ = out.Close() }()

	_, err = io.Copy(out, in)
	return err
}

func extensionDir() string {
	switch runtime.GOOS {
	case "darwin":
		home, _ := os.UserHomeDir()
		return filepath.Join(home, "Library", "Application Support", "vicinae", "extensions", "cairn")
	default:
		dataHome := os.Getenv("XDG_DATA_HOME")
		if dataHome == "" {
			home, _ := os.UserHomeDir()
			dataHome = filepath.Join(home, ".local", "share")
		}
		return filepath.Join(dataHome, "vicinae", "extensions", "cairn")
	}
}

func extractTarGz(archivePath, destDir string) error {
	f, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()

	gz, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer func() { _ = gz.Close() }()

	tr := tar.NewReader(gz)
	for {
		hdr, err := tr.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return err
		}

		// Guard against path traversal.
		target := filepath.Join(destDir, filepath.Clean(hdr.Name))
		if !strings.HasPrefix(target+string(filepath.Separator), destDir+string(filepath.Separator)) {
			return fmt.Errorf("archive path traversal rejected: %s", hdr.Name)
		}

		switch hdr.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0o755); err != nil {
				return err
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
				return err
			}
			out, err := os.Create(target)
			if err != nil {
				return err
			}
			_, copyErr := io.Copy(out, tr)
			_ = out.Close()
			if copyErr != nil {
				return copyErr
			}
		}
	}
	return nil
}
