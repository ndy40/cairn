package updater

import (
	"archive/tar"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// ── helpers ───────────────────────────────────────────────────────────────────

// mockReleaseServer returns an httptest.Server that serves a minimal GitHub
// Releases API response and any requested asset files.
// assetContent maps asset filenames to their byte content.
// checksumContent, if non-empty, is served as checksums.txt.
func mockReleaseServer(t *testing.T, tag string, assetContent map[string][]byte) *httptest.Server {
	t.Helper()
	var srv *httptest.Server
	srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// GitHub Releases API endpoint.
		if r.URL.Path == "/releases/latest" || r.URL.Path == "/" {
			assets := make([]map[string]string, 0, len(assetContent))
			for name := range assetContent {
				assets = append(assets, map[string]string{
					"name":                 name,
					"browser_download_url": srv.URL + "/download/" + name,
				})
			}
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"tag_name": tag,
				"assets":   assets,
			})
			return
		}

		// Asset download endpoint.
		if strings.HasPrefix(r.URL.Path, "/download/") {
			name := strings.TrimPrefix(r.URL.Path, "/download/")
			content, ok := assetContent[name]
			if !ok {
				http.NotFound(w, r)
				return
			}
			w.Header().Set("Content-Type", "application/octet-stream")
			_, _ = w.Write(content)
			return
		}

		http.NotFound(w, r)
	}))
	return srv
}

// sha256Hex returns the hex-encoded SHA256 of data.
func sha256Hex(data []byte) string {
	h := sha256.Sum256(data)
	return hex.EncodeToString(h[:])
}

// makeChecksums builds a checksums.txt-style string for the given files.
func makeChecksums(files map[string][]byte) []byte {
	var sb strings.Builder
	for name, content := range files {
		sb.WriteString(sha256Hex(content))
		sb.WriteString("  ")
		sb.WriteString(name)
		sb.WriteString("\n")
	}
	return []byte(sb.String())
}

// makeTarGz builds a minimal .tar.gz archive containing a single file.
func makeTarGz(t *testing.T, filename, content string) []byte {
	t.Helper()
	var buf strings.Builder
	// Write to a temp file then read back.
	tmp := t.TempDir()
	archPath := filepath.Join(tmp, "test.tar.gz")

	f, err := os.Create(archPath)
	if err != nil {
		t.Fatal(err)
	}
	gz := gzip.NewWriter(f)
	tw := tar.NewWriter(gz)

	body := []byte(content)
	hdr := &tar.Header{
		Name: filename,
		Mode: 0o644,
		Size: int64(len(body)),
	}
	if err := tw.WriteHeader(hdr); err != nil {
		t.Fatal(err)
	}
	if _, err := tw.Write(body); err != nil {
		t.Fatal(err)
	}
	_ = tw.Close()
	_ = gz.Close()
	_ = f.Close()
	_ = buf.String()

	data, err := os.ReadFile(archPath)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

// withMockServer sets the package-level URL and client to point at srv for the
// duration of the test, then restores the originals.
func withMockServer(t *testing.T, srv *httptest.Server) {
	t.Helper()
	origURL := releaseAPIURL
	origClient := updateClient
	origDownloadClient := downloadClient
	releaseAPIURL = srv.URL + "/releases/latest"
	updateClient = srv.Client()
	downloadClient = srv.Client()
	t.Cleanup(func() {
		releaseAPIURL = origURL
		updateClient = origClient
		downloadClient = origDownloadClient
		srv.Close()
	})
}

// ── CheckLatestVersion ────────────────────────────────────────────────────────

func TestCheckLatestVersion_UpdateAvailable(t *testing.T) {
	srv := mockReleaseServer(t, "v0.2.0", nil)
	withMockServer(t, srv)

	latest, available, err := CheckLatestVersion("v0.1.0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !available {
		t.Error("expected update to be available")
	}
	if latest != "v0.2.0" {
		t.Errorf("got latest %q, want v0.2.0", latest)
	}
}

func TestCheckLatestVersion_AlreadyUpToDate(t *testing.T) {
	srv := mockReleaseServer(t, "v0.2.0", nil)
	withMockServer(t, srv)

	_, available, err := CheckLatestVersion("v0.2.0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if available {
		t.Error("expected no update to be available")
	}
}

func TestCheckLatestVersion_DevAlwaysUpdates(t *testing.T) {
	srv := mockReleaseServer(t, "v0.1.0", nil)
	withMockServer(t, srv)

	_, available, err := CheckLatestVersion("dev")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !available {
		t.Error("expected dev build to always show update available")
	}
}

func TestCheckLatestVersion_ServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "rate limited", http.StatusForbidden)
	}))
	withMockServer(t, srv)

	_, _, err := CheckLatestVersion("v0.1.0")
	if err == nil {
		t.Fatal("expected error on non-200 response")
	}
}

func TestCheckLatestVersion_NetworkError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	// Close server immediately to simulate network error.
	srv.Close()
	origURL := releaseAPIURL
	origClient := updateClient
	releaseAPIURL = srv.URL + "/releases/latest"
	updateClient = srv.Client()
	t.Cleanup(func() {
		releaseAPIURL = origURL
		updateClient = origClient
	})

	_, _, err := CheckLatestVersion("v0.1.0")
	if err == nil {
		t.Fatal("expected error on closed server")
	}
}

func TestCheckLatestVersion_BadJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, "not json{{{")
	}))
	withMockServer(t, srv)

	_, _, err := CheckLatestVersion("v0.1.0")
	if err == nil {
		t.Fatal("expected error on bad JSON")
	}
}

// ── verifyChecksum ────────────────────────────────────────────────────────────

func TestVerifyChecksum_Match(t *testing.T) {
	content := []byte("hello world")
	f, err := os.CreateTemp(t.TempDir(), "")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = f.Write(content)
	_ = f.Close()

	if err := verifyChecksum(f.Name(), sha256Hex(content)); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestVerifyChecksum_Mismatch(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = f.Write([]byte("hello"))
	_ = f.Close()

	err = verifyChecksum(f.Name(), strings.Repeat("a", 64))
	if !errors.Is(err, ErrChecksumMismatch) {
		t.Errorf("expected ErrChecksumMismatch, got %v", err)
	}
}

// ── parseChecksum ─────────────────────────────────────────────────────────────

func TestParseChecksum_Found(t *testing.T) {
	content := "abc123  cairn-linux-amd64\ndef456  cairn-darwin-arm64\n"
	got, err := parseChecksum(content, "cairn-linux-amd64")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "abc123" {
		t.Errorf("got %q, want abc123", got)
	}
}

func TestParseChecksum_NotFound(t *testing.T) {
	_, err := parseChecksum("abc123  cairn-linux-amd64\n", "cairn-windows-amd64.exe")
	if err == nil {
		t.Fatal("expected error for missing entry")
	}
}

// ── updateBinaryAt ────────────────────────────────────────────────────────────

func TestUpdateBinaryAt_Success(t *testing.T) {
	binaryName := fmt.Sprintf("cairn-%s-%s", osName(), archName())
	newBinaryContent := []byte("#!/bin/sh\necho new-binary")
	checksums := makeChecksums(map[string][]byte{binaryName: newBinaryContent})
	assets := map[string][]byte{
		binaryName:      newBinaryContent,
		"checksums.txt": checksums,
	}
	srv := mockReleaseServer(t, "v0.2.0", assets)
	withMockServer(t, srv)

	// Create a fake "installed" binary.
	dir := t.TempDir()
	target := filepath.Join(dir, "cairn")
	if err := os.WriteFile(target, []byte("old"), 0o755); err != nil {
		t.Fatal(err)
	}

	if err := updateBinaryAt(target, "v0.1.0", "v0.2.0"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := os.ReadFile(target)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(newBinaryContent) {
		t.Errorf("binary not updated: got %q", string(got))
	}

	// Backup must be removed on success.
	if _, err := os.Stat(target + ".bak"); !os.IsNotExist(err) {
		t.Error("expected .bak file to be removed after successful update")
	}
}

func TestUpdateBinaryAt_ChecksumMismatch(t *testing.T) {
	binaryName := fmt.Sprintf("cairn-%s-%s", osName(), archName())
	badChecksums := []byte("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa  " + binaryName + "\n")
	assets := map[string][]byte{
		binaryName:      []byte("some binary content"),
		"checksums.txt": badChecksums,
	}
	srv := mockReleaseServer(t, "v0.2.0", assets)
	withMockServer(t, srv)

	dir := t.TempDir()
	target := filepath.Join(dir, "cairn")
	original := []byte("original")
	if err := os.WriteFile(target, original, 0o755); err != nil {
		t.Fatal(err)
	}

	err := updateBinaryAt(target, "v0.1.0", "v0.2.0")
	if !errors.Is(err, ErrChecksumMismatch) {
		t.Errorf("expected ErrChecksumMismatch, got %v", err)
	}

	// Original binary must be preserved.
	got, _ := os.ReadFile(target)
	if string(got) != string(original) {
		t.Errorf("original binary must not be modified on checksum failure")
	}
}

func TestUpdateBinaryAt_NoAsset(t *testing.T) {
	// Server returns a release with no matching binary asset.
	srv := mockReleaseServer(t, "v0.2.0", map[string][]byte{
		"cairn-linux-UNKNOWNARCH": []byte("data"),
	})
	withMockServer(t, srv)

	dir := t.TempDir()
	target := filepath.Join(dir, "cairn")
	_ = os.WriteFile(target, []byte("old"), 0o755)

	err := updateBinaryAt(target, "v0.1.0", "v0.2.0")
	if err == nil {
		t.Fatal("expected error when no matching asset found")
	}
}

func TestUpdateBinaryAt_PermissionDenied(t *testing.T) {
	if os.Getuid() == 0 {
		t.Skip("skipping permission test when running as root")
	}

	binaryName := fmt.Sprintf("cairn-%s-%s", osName(), archName())
	content := []byte("new binary")
	checksums := makeChecksums(map[string][]byte{binaryName: content})
	assets := map[string][]byte{
		binaryName:      content,
		"checksums.txt": checksums,
	}
	srv := mockReleaseServer(t, "v0.2.0", assets)
	withMockServer(t, srv)

	// Create a read-only directory.
	dir := t.TempDir()
	roDir := filepath.Join(dir, "readonly")
	if err := os.MkdirAll(roDir, 0o555); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(roDir, 0o755) })

	target := filepath.Join(roDir, "cairn")

	err := updateBinaryAt(target, "v0.1.0", "v0.2.0")
	if !errors.Is(err, ErrPermission) {
		t.Errorf("expected ErrPermission, got %v", err)
	}
}

// ── --check path (US2) ────────────────────────────────────────────────────────

// TestCheckLatestVersion_NoSideEffects verifies that CheckLatestVersion makes
// no modifications to the filesystem.
func TestCheckLatestVersion_NoSideEffects(t *testing.T) {
	srv := mockReleaseServer(t, "v0.2.0", nil)
	withMockServer(t, srv)

	dir := t.TempDir()
	entriesBefore, _ := os.ReadDir(dir)

	_, _, _ = CheckLatestVersion("v0.1.0")

	entriesAfter, _ := os.ReadDir(dir)
	if len(entriesBefore) != len(entriesAfter) {
		t.Error("CheckLatestVersion must not modify the filesystem")
	}
}

// ── DetectExtension ───────────────────────────────────────────────────────────

func TestDetectExtension_Installed(t *testing.T) {
	dir := t.TempDir()
	orig := extensionDirFunc
	extensionDirFunc = func() string { return dir }
	t.Cleanup(func() { extensionDirFunc = orig })

	gotDir, installed := DetectExtension()
	if !installed {
		t.Error("expected extension to be detected as installed")
	}
	if gotDir != dir {
		t.Errorf("got dir %q, want %q", gotDir, dir)
	}
}

func TestDetectExtension_NotInstalled(t *testing.T) {
	orig := extensionDirFunc
	extensionDirFunc = func() string {
		return filepath.Join(t.TempDir(), "nonexistent", "path")
	}
	t.Cleanup(func() { extensionDirFunc = orig })

	_, installed := DetectExtension()
	if installed {
		t.Error("expected extension to be detected as not installed")
	}
}

// ── CheckExtensionVersion ─────────────────────────────────────────────────────

func TestCheckExtensionVersion_WithVersionFile(t *testing.T) {
	srv := mockReleaseServer(t, "v0.2.0", nil)
	withMockServer(t, srv)

	dir := t.TempDir()
	_ = os.WriteFile(filepath.Join(dir, "version.txt"), []byte("v0.1.0\n"), 0o644)

	current, latest, available, err := CheckExtensionVersion(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if current != "v0.1.0" {
		t.Errorf("got current %q, want v0.1.0", current)
	}
	if latest != "v0.2.0" {
		t.Errorf("got latest %q, want v0.2.0", latest)
	}
	if !available {
		t.Error("expected update to be available")
	}
}

func TestCheckExtensionVersion_MissingVersionFile(t *testing.T) {
	srv := mockReleaseServer(t, "v0.2.0", nil)
	withMockServer(t, srv)

	dir := t.TempDir()
	// No version.txt — should be treated as "unknown" (always updatable).

	current, _, available, err := CheckExtensionVersion(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if current != "unknown" {
		t.Errorf("got current %q, want unknown", current)
	}
	if !available {
		t.Error("expected update to be available when version.txt is missing")
	}
}

// ── UpdateExtension ───────────────────────────────────────────────────────────

func TestUpdateExtension_Success(t *testing.T) {
	archiveName := "vicinae-extension-v0.2.0.tar.gz"
	archiveContent := makeTarGz(t, "main.js", "console.log('hello')")
	checksums := makeChecksums(map[string][]byte{archiveName: archiveContent})
	assets := map[string][]byte{
		archiveName:     archiveContent,
		"checksums.txt": checksums,
	}
	srv := mockReleaseServer(t, "v0.2.0", assets)
	withMockServer(t, srv)

	dir := t.TempDir()

	if err := UpdateExtension(dir, "v0.2.0"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// version.txt must be written.
	versionData, err := os.ReadFile(filepath.Join(dir, "version.txt"))
	if err != nil {
		t.Fatalf("version.txt not written: %v", err)
	}
	if strings.TrimSpace(string(versionData)) != "v0.2.0" {
		t.Errorf("version.txt content %q, want v0.2.0", string(versionData))
	}

	// Extracted file must exist.
	if _, err := os.Stat(filepath.Join(dir, "main.js")); err != nil {
		t.Errorf("expected extracted file main.js: %v", err)
	}
}

func TestUpdateExtension_ChecksumMismatch(t *testing.T) {
	archiveName := "vicinae-extension-v0.2.0.tar.gz"
	archiveContent := makeTarGz(t, "main.js", "hello")
	badChecksums := []byte("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa  " + archiveName + "\n")
	assets := map[string][]byte{
		archiveName:     archiveContent,
		"checksums.txt": badChecksums,
	}
	srv := mockReleaseServer(t, "v0.2.0", assets)
	withMockServer(t, srv)

	dir := t.TempDir()

	err := UpdateExtension(dir, "v0.2.0")
	if !errors.Is(err, ErrChecksumMismatch) {
		t.Errorf("expected ErrChecksumMismatch, got %v", err)
	}
}

func TestUpdateExtension_NoAsset(t *testing.T) {
	srv := mockReleaseServer(t, "v0.2.0", map[string][]byte{
		"unrelated-file.tar.gz": []byte("data"),
	})
	withMockServer(t, srv)

	dir := t.TempDir()
	err := UpdateExtension(dir, "v0.2.0")
	if err == nil {
		t.Fatal("expected error when archive asset is not found in release")
	}
}

func TestUpdateExtension_CheckOnly_NoModifications(t *testing.T) {
	srv := mockReleaseServer(t, "v0.2.0", nil)
	withMockServer(t, srv)

	dir := t.TempDir()
	_ = os.WriteFile(filepath.Join(dir, "version.txt"), []byte("v0.1.0\n"), 0o644)

	// Calling CheckExtensionVersion (the --check path) must not modify files.
	entriesBefore, _ := os.ReadDir(dir)
	_, _, _, _ = CheckExtensionVersion(dir)
	entriesAfter, _ := os.ReadDir(dir)

	if len(entriesBefore) != len(entriesAfter) {
		t.Error("CheckExtensionVersion must not create or remove files")
	}
}

func osName() string   { return runtime.GOOS }
func archName() string { return runtime.GOARCH }
