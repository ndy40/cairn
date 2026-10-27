package main

// action is what the launcher does in response to the shortcut, --toggle or
// opening a link.
type action int

const (
	// actionShow brings the window forward with search focused.
	actionShow action = iota
	// actionHideApp hides the whole app so focus returns to the previous one (macOS).
	actionHideApp
	// actionMinimise minimises the window (Linux).
	actionMinimise
	// actionHideWindow hides the window, leaving the app running (macOS, or
	// Linux background mode with a tray).
	actionHideWindow
	// actionQuit quits the app (closing the window on Linux without a tray
	// to return to).
	actionQuit
)

// shortcutAction decides what the global shortcut does: bring Cairn forward
// unless its window is already in front, in which case put it away the
// platform's way. background is the effective Linux "Run in background" mode
// (on, with a tray available).
func shortcutAction(goos string, visible, focused, background bool) action {
	if !visible || !focused {
		return actionShow
	}
	return putAwayAction(goos, background)
}

// putAwayAction is how Cairn gets out of the way: hide the app on macOS;
// on Linux hide to the tray in background mode, otherwise minimise so the
// window stays reachable from the taskbar.
func putAwayAction(goos string, background bool) action {
	switch {
	case goos == "darwin":
		return actionHideApp
	case background:
		return actionHideWindow
	default:
		return actionMinimise
	}
}

// closeAction decides what closing the window does: on macOS the app keeps
// running (Dock and shortcut bring it back); on Linux it keeps running only
// when "Run in background" is on and a tray can show it, and quits otherwise
// so it is never left running with no way back.
func closeAction(goos string, runInBackground, trayAvailable bool) action {
	if goos == "darwin" || (runInBackground && trayAvailable) {
		return actionHideWindow
	}
	return actionQuit
}

// Keys in SingleInstanceOptions.AdditionalData carrying the second launch's
// window activation token.
const (
	dataActivationToken = "activation_token" // XDG_ACTIVATION_TOKEN (Wayland)
	dataStartupID       = "startup_id"       // DESKTOP_STARTUP_ID (X11)
)

// activationToken picks the token a second launch forwarded, preferring the
// Wayland one.
func activationToken(data map[string]string) string {
	if token := data[dataActivationToken]; token != "" {
		return token
	}
	return data[dataStartupID]
}
