package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/ndy40/cairn/internal/updater"
)

func runUpdate(args []string) {
	fs := flag.NewFlagSet("update", flag.ContinueOnError)
	fs.Usage = func() {
		printCommandHelp("update")
		os.Exit(0)
	}
	checkOnly := fs.Bool("check", false, "check for updates without applying them")
	ext := fs.Bool("extension", false, "update the Vicinae extension instead of the CLI binary")
	if err := fs.Parse(args); err != nil {
		os.Exit(1)
	}
	if *ext {
		runUpdateExtension(*checkOnly)
	} else {
		runUpdateBinary(*checkOnly)
	}
}

func runUpdateBinary(checkOnly bool) {
	latest, available, err := updater.CheckLatestVersion(version)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "cairn update: %v\n", err)
		os.Exit(1)
	}
	if !available {
		fmt.Printf("cairn: already up to date (%s)\n", latest)
		return
	}
	if checkOnly {
		fmt.Printf("cairn: current version %s, latest %s (update available)\n", version, latest)
		return
	}
	fmt.Printf("cairn: current version %s, latest %s\n", version, latest)
	if err := updater.UpdateBinary(version, latest); err != nil {
		switch {
		case errors.Is(err, updater.ErrChecksumMismatch):
			_, _ = fmt.Fprintln(os.Stderr, "cairn update: checksum mismatch for downloaded binary")
			os.Exit(3)
		case errors.Is(err, updater.ErrPermission):
			_, _ = fmt.Fprintln(os.Stderr, "cairn update: permission denied: cannot write to install directory")
			os.Exit(4)
		default:
			_, _ = fmt.Fprintf(os.Stderr, "cairn update: %v\n", err)
			os.Exit(1)
		}
	}
}

func runUpdateExtension(checkOnly bool) {
	dir, installed := updater.DetectExtension()
	if !installed {
		fmt.Println("cairn: extension not installed; run the install script with --with-extension to install it")
		return
	}
	current, latest, available, err := updater.CheckExtensionVersion(dir)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "cairn update: %v\n", err)
		os.Exit(1)
	}
	if !available {
		fmt.Printf("cairn: extension already up to date (%s)\n", latest)
		return
	}
	if checkOnly {
		fmt.Printf("cairn: extension current %s, latest %s (update available)\n", current, latest)
		return
	}
	fmt.Printf("cairn: extension current %s, latest %s\n", current, latest)
	if err := updater.UpdateExtension(dir, latest); err != nil {
		switch {
		case errors.Is(err, updater.ErrChecksumMismatch):
			_, _ = fmt.Fprintln(os.Stderr, "cairn update: checksum mismatch for extension archive")
			os.Exit(3)
		default:
			_, _ = fmt.Fprintf(os.Stderr, "cairn update: %v\n", err)
			os.Exit(1)
		}
	}
}
