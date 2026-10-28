package main

import "testing"

func TestVersionFromConfig(t *testing.T) {
	config := `version: '3'

info:
  productName: "Cairn" # The name of the application
  version: "0.8.1" # The application version
  # version: "9.9.9"

other:
  version: "1.0"
`
	if got := versionFromConfig([]byte(config)); got != "0.8.1" {
		t.Errorf("versionFromConfig = %q, want 0.8.1", got)
	}
	if got := versionFromConfig([]byte("version: '3'\n")); got != "dev" {
		t.Errorf("versionFromConfig without info = %q, want dev", got)
	}
	if appVersion == "" || appVersion == "dev" {
		t.Errorf("appVersion = %q, want build/config.yml's info.version", appVersion)
	}
}
