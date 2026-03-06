package clipboard

import (
	"errors"
	"os"
	"os/exec"
	"strings"

	"github.com/atotto/clipboard"
)

// Read reads plain text from the system clipboard.
// On Linux with Wayland, it tries wl-paste first (if available) before
// falling back to xclip/xsel via atotto/clipboard.
func Read() (string, error) {
	text, err := read()
	if err != nil {
		return "", err
	}
	text = strings.TrimSpace(text)
	if text == "" {
		return "", errors.New("clipboard is empty")
	}
	return text, nil
}

func read() (string, error) {
	// On Linux with Wayland: prefer wl-paste to avoid the X11/Wayland
	// clipboard split where xclip reads from the X11 clipboard (empty) while
	// content was copied in a native Wayland application.
	if os.Getenv("WAYLAND_DISPLAY") != "" {
		if text, err := wlPaste(); err == nil {
			return text, nil
		}
		// wl-paste not installed or failed — fall through to xclip/xsel.
	}

	text, err := clipboard.ReadAll()
	if err != nil {
		return "", errors.New(
			"clipboard unavailable: " + err.Error() +
				"\nOn Wayland, install wl-clipboard: sudo apt install wl-clipboard" +
				"\nOn X11, install xclip: sudo apt install xclip",
		)
	}
	return text, nil
}

// wlPaste runs wl-paste --no-newline and returns its output.
// Returns an error if wl-paste is not installed or exits non-zero.
func wlPaste() (string, error) {
	path, err := exec.LookPath("wl-paste")
	if err != nil {
		return "", errors.New("wl-paste not found")
	}
	out, err := exec.Command(path, "--no-newline").Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}
