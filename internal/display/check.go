package display

import (
	"os"
	"os/exec"
)

// DisplayType represents the detected display environment.
type DisplayType int

const (
	Unknown DisplayType = iota
	Wayland
	X11
)

// CheckResult holds the outcome of the prerequisite check.
type CheckResult struct {
	DisplayType DisplayType
	ToolFound   bool
	MissingTool string
	InstallHint string
	ShouldBlock bool
}

// CheckPrerequisites detects the display environment and verifies the
// required clipboard tool is available. Wayland takes precedence when
// both WAYLAND_DISPLAY and DISPLAY are set.
func CheckPrerequisites() CheckResult {
	if os.Getenv("WAYLAND_DISPLAY") != "" {
		return checkWayland()
	}
	if os.Getenv("DISPLAY") != "" {
		return checkX11()
	}
	return CheckResult{
		DisplayType: Unknown,
		ToolFound:   false,
		MissingTool: "",
		InstallHint: "Warning: display environment not detected. Clipboard paste (Ctrl+P) may not be available.",
		ShouldBlock: false,
	}
}

func checkWayland() CheckResult {
	if _, err := exec.LookPath("wl-paste"); err == nil {
		return CheckResult{DisplayType: Wayland, ToolFound: true}
	}
	return CheckResult{
		DisplayType: Wayland,
		ToolFound:   false,
		MissingTool: "wl-paste",
		InstallHint: "Wayland detected. Please install wl-clipboard:\n" +
			"  sudo apt install wl-clipboard    # Debian/Ubuntu\n" +
			"  sudo pacman -S wl-clipboard      # Arch\n" +
			"  sudo dnf install wl-clipboard    # Fedora",
		ShouldBlock: true,
	}
}

func checkX11() CheckResult {
	if _, err := exec.LookPath("xclip"); err == nil {
		return CheckResult{DisplayType: X11, ToolFound: true}
	}
	if _, err := exec.LookPath("xsel"); err == nil {
		return CheckResult{DisplayType: X11, ToolFound: true}
	}
	return CheckResult{
		DisplayType: X11,
		ToolFound:   false,
		MissingTool: "xclip/xsel",
		InstallHint: "X11 detected. Please install xclip or xsel:\n" +
			"  sudo apt install xclip           # Debian/Ubuntu (recommended)\n" +
			"  sudo pacman -S xclip             # Arch\n" +
			"  sudo dnf install xclip           # Fedora",
		ShouldBlock: true,
	}
}
