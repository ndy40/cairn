package updater

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// makeDesktopTarGz builds a Linux desktop release archive holding the given
// binary and icon under the archive's top-level directory.
func makeDesktopTarGz(t *testing.T, dir, binary, icon string) []byte {
	t.Helper()
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gz)
	for name, body := range map[string]string{desktopName: binary, desktopName + ".png": icon} {
		if err := tw.WriteHeader(&tar.Header{Name: dir + "/" + name, Mode: 0o755, Size: int64(len(body))}); err != nil {
			t.Fatal(err)
		}
		if _, err := tw.Write([]byte(body)); err != nil {
			t.Fatal(err)
		}
	}
	_ = tw.Close()
	_ = gz.Close()
	return buf.Bytes()
}

// desktopProcess records the stubbed process calls made during a test.
type desktopProcess struct {
	running bool
	stopped bool
	started string
}

// withLinuxDesktop points the desktop updater at a temp data dir, stubs
// process control and returns the path of an installed old binary.
func withLinuxDesktop(t *testing.T, proc *desktopProcess) (binPath, dataDir string) {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("desktop app is not updated on Windows")
	}
	binPath = filepath.Join(t.TempDir(), desktopName)
	if err := os.WriteFile(binPath, []byte("old"), 0o755); err != nil {
		t.Fatal(err)
	}
	dataDir = t.TempDir()

	origGOOS, origData := desktopGOOS, desktopDataDirFunc
	origRunning, origStop, origStart := desktopRunning, stopDesktop, startDesktop
	desktopGOOS = "linux"
	desktopDataDirFunc = func() string { return dataDir }
	desktopRunning = func(string) []int {
		if proc.running {
			return []int{1}
		}
		return nil
	}
	stopDesktop = func(string) { proc.stopped = true }
	startDesktop = func(path string) error { proc.started = path; return nil }
	t.Cleanup(func() {
		desktopGOOS, desktopDataDirFunc = origGOOS, origData
		desktopRunning, stopDesktop, startDesktop = origRunning, origStop, origStart
	})
	return binPath, dataDir
}

func linuxDesktopArchive() (dir, name string) {
	dir = fmt.Sprintf("%s-linux-%s", desktopName, runtime.GOARCH)
	return dir, dir + ".tar.gz"
}

func TestUpdateDesktop_Linux_ReplacesAndRestarts(t *testing.T) {
	proc := &desktopProcess{running: true}
	binPath, dataDir := withLinuxDesktop(t, proc)

	dir, archiveName := linuxDesktopArchive()
	archive := makeDesktopTarGz(t, dir, "new", "icon")
	srv := mockReleaseServer(t, "v0.2.0", map[string][]byte{
		archiveName:     archive,
		"checksums.txt": makeChecksums(map[string][]byte{archiveName: archive}),
	})
	withMockServer(t, srv)

	if err := UpdateDesktop(binPath, "v0.2.0"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, _ := os.ReadFile(binPath)
	if string(got) != "new" {
		t.Errorf("binary content %q, want %q", got, "new")
	}
	if info, err := os.Stat(binPath); err != nil || info.Mode().Perm()&0o100 == 0 {
		t.Errorf("binary is not executable: %v %v", info.Mode(), err)
	}
	if icon, _ := os.ReadFile(filepath.Join(dataDir, desktopName+".png")); string(icon) != "icon" {
		t.Errorf("icon content %q, want %q", icon, "icon")
	}
	if v, _ := os.ReadFile(filepath.Join(dataDir, desktopVersionFile)); strings.TrimSpace(string(v)) != "v0.2.0" {
		t.Errorf("version file %q, want v0.2.0", v)
	}
	if !proc.stopped || proc.started != binPath {
		t.Errorf("running app not restarted: stopped=%v started=%q", proc.stopped, proc.started)
	}
}

func TestUpdateDesktop_NotRunning_NotStarted(t *testing.T) {
	proc := &desktopProcess{}
	binPath, _ := withLinuxDesktop(t, proc)

	dir, archiveName := linuxDesktopArchive()
	archive := makeDesktopTarGz(t, dir, "new", "icon")
	srv := mockReleaseServer(t, "v0.2.0", map[string][]byte{
		archiveName:     archive,
		"checksums.txt": makeChecksums(map[string][]byte{archiveName: archive}),
	})
	withMockServer(t, srv)

	if err := UpdateDesktop(binPath, "v0.2.0"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if proc.stopped || proc.started != "" {
		t.Errorf("app should not be touched when not running: stopped=%v started=%q", proc.stopped, proc.started)
	}
}

func TestUpdateDesktop_ChecksumMismatch_LeavesAppUntouched(t *testing.T) {
	proc := &desktopProcess{running: true}
	binPath, dataDir := withLinuxDesktop(t, proc)

	dir, archiveName := linuxDesktopArchive()
	archive := makeDesktopTarGz(t, dir, "new", "icon")
	srv := mockReleaseServer(t, "v0.2.0", map[string][]byte{
		archiveName:     archive,
		"checksums.txt": []byte(strings.Repeat("a", 64) + "  " + archiveName + "\n"),
	})
	withMockServer(t, srv)

	err := UpdateDesktop(binPath, "v0.2.0")
	if !errors.Is(err, ErrChecksumMismatch) {
		t.Fatalf("expected ErrChecksumMismatch, got %v", err)
	}
	if got, _ := os.ReadFile(binPath); string(got) != "old" {
		t.Errorf("binary changed to %q after checksum mismatch", got)
	}
	if _, err := os.Stat(filepath.Join(dataDir, desktopVersionFile)); err == nil {
		t.Error("version file written after checksum mismatch")
	}
	if proc.stopped {
		t.Error("app stopped after checksum mismatch")
	}
}

func TestUpdateDesktop_MissingChecksums(t *testing.T) {
	binPath, _ := withLinuxDesktop(t, &desktopProcess{})

	dir, archiveName := linuxDesktopArchive()
	srv := mockReleaseServer(t, "v0.2.0", map[string][]byte{
		archiveName: makeDesktopTarGz(t, dir, "new", "icon"),
	})
	withMockServer(t, srv)

	if err := UpdateDesktop(binPath, "v0.2.0"); err == nil {
		t.Fatal("expected an error when checksums.txt is missing")
	}
	if got, _ := os.ReadFile(binPath); string(got) != "old" {
		t.Errorf("binary changed to %q without a checksum", got)
	}
}

func TestCheckDesktopVersion(t *testing.T) {
	_, dataDir := withLinuxDesktop(t, &desktopProcess{})
	srv := mockReleaseServer(t, "v0.2.0", map[string][]byte{})
	withMockServer(t, srv)

	current, latest, available, err := CheckDesktopVersion()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if current != "unknown" || latest != "v0.2.0" || !available {
		t.Errorf("without version file got (%q, %q, %v), want (unknown, v0.2.0, true)", current, latest, available)
	}

	if err := os.WriteFile(filepath.Join(dataDir, desktopVersionFile), []byte("v0.2.0\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	current, _, available, err = CheckDesktopVersion()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if current != "v0.2.0" || available {
		t.Errorf("with current version got (%q, %v), want (v0.2.0, false)", current, available)
	}
}

func TestDetectDesktop(t *testing.T) {
	orig := desktopPathFunc
	t.Cleanup(func() { desktopPathFunc = orig })

	desktopPathFunc = func() string { return "" }
	if _, installed := DetectDesktop(); installed {
		t.Error("empty path reported as installed")
	}

	path := filepath.Join(t.TempDir(), desktopName)
	desktopPathFunc = func() string { return path }
	if _, installed := DetectDesktop(); installed {
		t.Error("missing app reported as installed")
	}
	if err := os.WriteFile(path, nil, 0o755); err != nil {
		t.Fatal(err)
	}
	if got, installed := DetectDesktop(); !installed || got != path {
		t.Errorf("DetectDesktop() = (%q, %v), want (%q, true)", got, installed, path)
	}
}

func TestIsDesktopRunning_IgnoresOtherInstalls(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("reads /proc")
	}
	// Whatever cairn-desktop this machine runs, it is not at this path.
	if pids := isDesktopRunning(filepath.Join(t.TempDir(), desktopName)); len(pids) != 0 {
		t.Errorf("matched processes %v for an app that is not installed there", pids)
	}
}
