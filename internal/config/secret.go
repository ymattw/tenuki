package config

import (
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
)

const Secret = "secret.json"

func SecretPath(username string) string {
	return filepath.Join(xdg.StateHome, "tenuki", username, Secret)
}

// Return the matching $XDG_STATE_HOME/tenuki/<username|*>/secret.json files
func SearchSecrets(username string) []string {
	if username != "" {
		p := SecretPath(username)
		if _, err := os.Stat(SecretPath(username)); err == nil {
			return []string{p}
		}
		return nil
	}

	var res []string
	entries, _ := os.ReadDir(filepath.Join(xdg.StateHome, "tenuki"))
	for _, entry := range entries {
		if entry.IsDir() {
			p := SecretPath(entry.Name())
			if _, err := os.Stat(p); err == nil {
				res = append(res, p)
			}
		}
	}
	return res
}
