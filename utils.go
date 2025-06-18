package main

import (
	"encoding/json"
	"os/exec"
	"strings"
)

func cloneActivity(activity DiscordActivity) (DiscordActivity, error) {
	var clone DiscordActivity

	data, err := json.Marshal(activity)
	if err != nil {
		return clone, err
	}

	err = json.Unmarshal(data, &clone)
	return clone, err
}

func runExternalScript(scriptPath string) (string, error) {
    out, err := exec.Command("/bin/sh", "-c", scriptPath).Output()
    if err != nil {
        return "", err
    }
    return strings.TrimSpace(string(out)), nil
}


