package updater

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// desktopName is the desktop app's executable and process name.
const desktopName = "cairn-desktop"

// desktopVersionFile records the installed desktop app version. The app does
// not report its own version, so the installer and UpdateDesktop write it.
const desktopVersionFile = "desktop-version.txt"

// Package-level vars allow tests to redirect paths and stub process control.
var (
	desktopPathFunc    = desktopPath
	desktopDataDirFunc = desktopDataDir
	desktopRunning     = isDesktopRunning
	stopDesktop        = stopDesktopProcess
	startDesktop       = startDesktopProcess
	desktopGOOS        = runtime.GOOS
)

// DetectDesktop returns the installed desktop app's path (the binary on
// Linux, the .app bundle on macOS) and whether it is installed.
func DetectDesktop() (path string, installed bool) {
	path = desktopPathFunc()
	if path == "" {
		return "", false
	}
	_, err := os.Stat(path)
	return path, err == nil
}

// CheckDesktopVersion reads the installed desktop app version and queries the
// releases API for the latest version. A missing version file (an install
// made before versions were recorded) reports "unknown" and an update.
func CheckDesktopVersion() (current, latest string, available bool, err error) {
	data, readErr := os.ReadFile(filepath.Join(desktopDataDirFunc(), desktopVersionFile))
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

// UpdateDesktop downloads the latest desktop app archive, verifies its
// checksum, replaces the app at path, records the new version and, if the app
// was running, restarts it.
func UpdateDesktop(path, latestVersion string) error {
	archiveName := desktopArchiveName()
	if archiveName == "" {
		fmt.Println("cairn update: the desktop app is not available on this platform")
		return nil
	}
	if !isWritable(filepath.Dir(path)) {
		return ErrPermission
	}

	info, err := fetchReleaseInfo(latestVersion)
	if err != nil {
		return err
	}
	archiveURL, checksumURL := findAssetURLs(info, archiveName)
	if archiveURL == "" {
		return fmt.Errorf("no release asset found for %s", archiveName)
	}
	if checksumURL == "" {
		return fmt.Errorf("no checksums.txt in release %s", info.TagName)
	}

	tmpDir, err := os.MkdirTemp("", "cairn-desktop-update-*")
	if err != nil {
		return fmt.Errorf("create temp dir: %w", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	checksumPath, err := downloadToTemp(tmpDir, checksumURL, ".cairn-checksums-*")
	if err != nil {
		return fmt.Errorf("download checksums: %w", err)
	}
	checksumContent, err := os.ReadFile(checksumPath)
	if err != nil {
		return fmt.Errorf("read checksums: %w", err)
	}
	expectedHash, err := parseChecksum(string(checksumContent), archiveName)
	if err != nil {
		return fmt.Errorf("find checksum for %s: %w", archiveName, err)
	}

	fmt.Printf("cairn: downloading desktop app %s...\n", latestVersion)
	archivePath := filepath.Join(tmpDir, archiveName)
	if err := downloadFile(archiveURL, archivePath); err != nil {
		return fmt.Errorf("download desktop archive: %w", err)
	}
	fmt.Println("cairn: verifying checksum...")
	if err := verifyChecksum(archivePath, expectedHash); err != nil {
		return err
	}

	wasRunning := len(desktopRunning(path)) > 0
	if desktopGOOS == "darwin" {
		err = replaceDesktopDarwin(path, archivePath)
	} else {
		err = replaceDesktopLinux(path, archivePath, tmpDir)
	}
	if err != nil {
		return err
	}

	dataDir := desktopDataDirFunc()
	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		return fmt.Errorf("create data directory: %w", err)
	}
	if err := os.WriteFile(filepath.Join(dataDir, desktopVersionFile), []byte(latestVersion+"\n"), 0o644); err != nil {
		return fmt.Errorf("write %s: %w", desktopVersionFile, err)
	}
	fmt.Printf("cairn: desktop app updated to %s\n", latestVersion)

	// The old version keeps running until it exits, so restart it to pick up
	// the update.
	if wasRunning {
		fmt.Println("cairn: restarting the desktop app...")
		stopDesktop(path)
		if err := startDesktop(path); err != nil {
			fmt.Printf("cairn: could not restart the desktop app (%v); start it with: %s --hidden\n", err, path)
		}
	}
	return nil
}

// replaceDesktopLinux swaps in the binary and icon from the release tarball.
// The binary is staged next to the target so the rename is atomic.
func replaceDesktopLinux(path, archivePath, tmpDir string) error {
	extractDir := filepath.Join(tmpDir, "extract")
	if err := os.MkdirAll(extractDir, 0o755); err != nil {
		return err
	}
	if err := extractTarGz(archivePath, extractDir); err != nil {
		return fmt.Errorf("extract desktop archive: %w", err)
	}
	unpacked := filepath.Join(extractDir, strings.TrimSuffix(filepath.Base(archivePath), ".tar.gz"))

	staged, err := os.CreateTemp(filepath.Dir(path), ".cairn-desktop-update-*")
	if err != nil {
		return fmt.Errorf("stage desktop binary: %w", err)
	}
	stagedPath := staged.Name()
	_ = staged.Close()
	defer func() { _ = os.Remove(stagedPath) }()

	if err := copyFile(filepath.Join(unpacked, desktopName), stagedPath); err != nil {
		return fmt.Errorf("stage desktop binary: %w", err)
	}
	if err := os.Chmod(stagedPath, 0o755); err != nil {
		return fmt.Errorf("chmod desktop binary: %w", err)
	}
	if err := os.Rename(stagedPath, path); err != nil {
		return fmt.Errorf("replace desktop binary: %w", err)
	}

	dataDir := desktopDataDirFunc()
	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		return fmt.Errorf("create data directory: %w", err)
	}
	if err := copyFile(filepath.Join(unpacked, desktopName+".png"), filepath.Join(dataDir, desktopName+".png")); err != nil {
		return fmt.Errorf("install desktop icon: %w", err)
	}
	return nil
}

// replaceDesktopDarwin swaps in the .app bundle from the release zip. ditto
// keeps the bundle's symlinks and code signature intact. The app is stopped
// first and the old bundle is restored if the swap fails.
func replaceDesktopDarwin(path, archivePath string) error {
	stage, err := os.MkdirTemp(filepath.Dir(path), ".cairn-desktop-update-*")
	if err != nil {
		return fmt.Errorf("stage desktop app: %w", err)
	}
	defer func() { _ = os.RemoveAll(stage) }()

	if out, err := exec.Command("ditto", "-x", "-k", archivePath, stage).CombinedOutput(); err != nil {
		return fmt.Errorf("unpack desktop app: %v: %s", err, strings.TrimSpace(string(out)))
	}
	newApp := filepath.Join(stage, desktopName+".app")
	if _, err := os.Stat(newApp); err != nil {
		return fmt.Errorf("unpack desktop app: %w", err)
	}

	stopDesktop(path)
	oldApp := filepath.Join(stage, "previous.app")
	if err := os.Rename(path, oldApp); err != nil {
		return fmt.Errorf("replace desktop app: %w", err)
	}
	if err := os.Rename(newApp, path); err != nil {
		_ = os.Rename(oldApp, path)
		return fmt.Errorf("replace desktop app: %w", err)
	}
	return nil
}

func desktopArchiveName() string {
	switch desktopGOOS {
	case "linux":
		return fmt.Sprintf("%s-linux-%s.tar.gz", desktopName, runtime.GOARCH)
	case "darwin":
		return desktopName + "-darwin-universal.zip"
	}
	return ""
}

// desktopPath returns where the installer puts the desktop app: on Linux the
// binary beside the cairn executable (falling back to PATH), on macOS the
// bundle in ~/Applications.
func desktopPath() string {
	switch runtime.GOOS {
	case "darwin":
		home, _ := os.UserHomeDir()
		return filepath.Join(home, "Applications", desktopName+".app")
	case "linux":
		if exe, err := os.Executable(); err == nil {
			if exe, err = filepath.EvalSymlinks(exe); err == nil {
				sibling := filepath.Join(filepath.Dir(exe), desktopName)
				if _, err := os.Stat(sibling); err == nil {
					return sibling
				}
			}
		}
		if p, err := exec.LookPath(desktopName); err == nil {
			return p
		}
	}
	return ""
}

// desktopDataDir holds the desktop app's icon and version file, matching the
// installer.
func desktopDataDir() string {
	if runtime.GOOS == "darwin" {
		home, _ := os.UserHomeDir()
		return filepath.Join(home, "Library", "Application Support", "cairn")
	}
	return filepath.Join(xdgDataHome(), "cairn")
}

// isDesktopRunning returns this user's cairn-desktop processes running the
// app at path. Copies installed elsewhere are left alone.
func isDesktopRunning(path string) []int {
	out, err := exec.Command("pgrep", "-x", "-u", strconv.Itoa(os.Getuid()), desktopName).Output()
	if err != nil {
		return nil
	}
	want := path
	if desktopGOOS == "darwin" {
		want = filepath.Join(path, "Contents", "MacOS", desktopName)
	}
	if resolved, err := filepath.EvalSymlinks(want); err == nil {
		want = resolved
	}
	var pids []int
	for _, field := range strings.Fields(string(out)) {
		pid, err := strconv.Atoi(field)
		if err == nil && processExe(pid) == want {
			pids = append(pids, pid)
		}
	}
	return pids
}

// processExe returns the executable path of pid, or "" if it is unknown.
func processExe(pid int) string {
	if desktopGOOS == "darwin" {
		out, err := exec.Command("ps", "-o", "comm=", "-p", strconv.Itoa(pid)).Output()
		if err != nil {
			return ""
		}
		return strings.TrimSpace(string(out))
	}
	exe, err := os.Readlink(fmt.Sprintf("/proc/%d/exe", pid))
	if err != nil {
		return ""
	}
	// A binary replaced on disk shows as "<path> (deleted)".
	return strings.TrimSuffix(exe, " (deleted)")
}

// stopDesktopProcess quits the app running from path and waits up to 5
// seconds for it to exit.
func stopDesktopProcess(path string) {
	for _, pid := range isDesktopRunning(path) {
		if p, err := os.FindProcess(pid); err == nil {
			_ = p.Signal(syscall.SIGTERM)
		}
	}
	for i := 0; i < 5 && len(isDesktopRunning(path)) > 0; i++ {
		time.Sleep(time.Second)
	}
}

// startDesktopProcess launches the app in the background, detached from this
// terminal so it outlives it.
func startDesktopProcess(path string) error {
	if runtime.GOOS == "darwin" {
		return exec.Command("open", "-g", path, "--args", "--hidden").Run()
	}
	if os.Getenv("WAYLAND_DISPLAY") == "" && os.Getenv("DISPLAY") == "" {
		return fmt.Errorf("no graphical session found")
	}
	cmd := exec.Command(path, "--hidden")
	detach(cmd)
	if err := cmd.Start(); err != nil {
		return err
	}
	return cmd.Process.Release()
}
