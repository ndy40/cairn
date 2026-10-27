//go:build !linux || gtk3

package main

import "github.com/wailsapp/wails/v3/pkg/application"

// applyActivationToken is only needed on Linux (GTK4); macOS raises the
// window without a token.
func applyActivationToken(*application.WebviewWindow, string) {}
