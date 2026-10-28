package main

import (
	"bufio"
	"bytes"
	_ "embed"
	"strings"
)

// buildConfig is build/config.yml, whose info.version the release workflow
// sets before building, so the app reports the version it was built as.
//
//go:embed build/config.yml
var buildConfig []byte

// appVersion is the version shown in About.
var appVersion = versionFromConfig(buildConfig)

// versionFromConfig reads info.version from build/config.yml without a YAML
// parser: the first "version:" line indented under the top-level "info:" key.
func versionFromConfig(data []byte) string {
	inInfo := false
	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		if !strings.HasPrefix(line, " ") && !strings.HasPrefix(line, "\t") {
			inInfo = strings.HasPrefix(trimmed, "info:")
			continue
		}
		if !inInfo || !strings.HasPrefix(trimmed, "version:") {
			continue
		}
		value := strings.TrimSpace(strings.TrimPrefix(trimmed, "version:"))
		if i := strings.Index(value, " #"); i >= 0 {
			value = strings.TrimSpace(value[:i])
		}
		return strings.Trim(value, `"'`)
	}
	return "dev"
}
