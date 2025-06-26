package config

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/ymattw/googs"
)

const Secret = "secret.json"

func SecretPath(username string) string {
	return filepath.Join(xdg.StateHome, "tenuki", username, Secret)
}

// TODO: Delete this migration code in a month
func MigrateSecret() {
	c, err := googs.LoadClient(Secret)
	if err != nil {
		return
	}

	newPath := SecretPath(c.Username)
	if _, err := os.Stat(newPath); errors.Is(err, os.ErrNotExist) {
		os.Rename(Secret, newPath)
	}
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
