package main

import (
	"encoding/xml"
	"os"
	"strings"
	"testing"
)

func TestLaunchAgentPlist(t *testing.T) {
	exe := "/Applications/Cairn & Co.app/Contents/MacOS/cairn-desktop"
	data := launchAgentPlist(exe)

	// Decode the plist's dict as alternating keys and values.
	var doc struct {
		Dict struct {
			Items []struct {
				XMLName xml.Name
				Text    string   `xml:",chardata"`
				Strings []string `xml:"string"`
			} `xml:",any"`
		} `xml:"dict"`
	}
	if err := xml.Unmarshal(data, &doc); err != nil {
		t.Fatalf("plist is not valid XML: %v\n%s", err, data)
	}
	values := map[string]any{}
	items := doc.Dict.Items
	for i := 0; i+1 < len(items); i += 2 {
		if items[i].XMLName.Local != "key" {
			t.Fatalf("item %d is <%s>, want <key>", i, items[i].XMLName.Local)
		}
		v := items[i+1]
		switch v.XMLName.Local {
		case "array":
			values[items[i].Text] = v.Strings
		case "true", "false":
			values[items[i].Text] = v.XMLName.Local
		default:
			values[items[i].Text] = v.Text
		}
	}

	if got := values["Label"]; got != launchAgentLabel {
		t.Errorf("Label = %v, want %s", got, launchAgentLabel)
	}
	args, _ := values["ProgramArguments"].([]string)
	if strings.Join(args, "|") != exe+"|--background" {
		t.Errorf("ProgramArguments = %q, want [%q --background]", args, exe)
	}
	if got := values["RunAtLoad"]; got != "true" {
		t.Errorf("RunAtLoad = %v, want true", got)
	}
	if !strings.Contains(string(data), "Cairn &amp; Co.app") {
		t.Error("executable path is not XML-escaped")
	}
}

func TestLaunchAgentWriteAndRemove(t *testing.T) {
	path := launchAgentPath(t.TempDir())

	if on, err := launchAgentExists(path); err != nil || on {
		t.Fatalf("before writing: exists = %t, %v; want false", on, err)
	}
	if err := writeLaunchAgent(path, "/usr/local/bin/cairn-desktop"); err != nil {
		t.Fatalf("writeLaunchAgent: %v", err)
	}
	if on, err := launchAgentExists(path); err != nil || !on {
		t.Fatalf("after writing: exists = %t, %v; want true", on, err)
	}
	data, err := os.ReadFile(path)
	if err != nil || !strings.Contains(string(data), "<string>/usr/local/bin/cairn-desktop</string>") {
		t.Errorf("written plist = %q, %v", data, err)
	}
	if err := removeLaunchAgent(path); err != nil {
		t.Fatalf("removeLaunchAgent: %v", err)
	}
	if on, _ := launchAgentExists(path); on {
		t.Error("LaunchAgent still exists after removal")
	}
	if err := removeLaunchAgent(path); err != nil {
		t.Errorf("removing a missing LaunchAgent: %v", err)
	}
}
