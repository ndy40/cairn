//go:build linux && !gtk3

package main

/*
#cgo pkg-config: gtk4
#include <stdlib.h>
#include <gtk/gtk.h>
*/
import "C"

import (
	"unsafe"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// applyActivationToken hands GTK the activation token (Wayland
// XDG_ACTIVATION_TOKEN, or X11 DESKTOP_STARTUP_ID) forwarded by a second
// launch, so the next gtk_window_present inside Show/Focus uses it. Without a
// token GNOME on Wayland refuses to raise the window and shows "Cairn is
// ready" instead. Wails has no API for this, hence the cgo.
func applyActivationToken(w *application.WebviewWindow, token string) {
	if token == "" {
		return
	}
	native := w.NativeWindow()
	if native == nil {
		return
	}
	debugf("applying activation token")
	application.InvokeSync(func() {
		cToken := C.CString(token)
		defer C.free(unsafe.Pointer(cToken))
		C.gtk_window_set_startup_id((*C.GtkWindow)(native), cToken)
	})
}
