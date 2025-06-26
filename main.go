// Package main offers a OGS client to play on the terminal.
package main

import (
	"flag"
	"log"

	"github.com/ymattw/googs"

	"github.com/ymattw/tenuki/internal/config"
	"github.com/ymattw/tenuki/internal/tui"
)

var (
	username = flag.String("u", "", "OGS username, only needed when you have multiple users logged in before")
)

func main() {
	flag.Parse()
	config.MigrateSecret()

	client, err := loadClient()
	if err != nil {
		log.Fatal(err)
	}

	app := tui.NewApp(client)
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

func loadClient() (*googs.Client, error) {
	secretFiles := config.SearchSecrets(*username)
	if len(secretFiles) > 1 {
		log.Fatalf("Username (-u) is needed to pick one from multiple secret files found: %q\n", secretFiles)
	}
	if len(secretFiles) == 0 {
		return googs.NewClient("", ""), nil
	}

	c, err := googs.LoadClient(secretFiles[0])
	if err != nil {
		// Use maybe-set fields from the incomplete client
		return googs.NewClient(c.ClientID, c.ClientSecret), nil
	}
	return c, nil
}
