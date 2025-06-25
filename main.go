// Package main offers a OGS client to play on the terminal.
package main

import (
	"log"

	"github.com/ymattw/googs"

	"github.com/ymattw/tenuki/internal/tui"
)

func main() {
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
	// TODO: load from XDG_HOME
	c, err := googs.LoadClient("secret.json")
	if err != nil {
		// Use maybe-set fields from the incomplete client
		return googs.NewClient(c.ClientID, c.ClientSecret), nil
	}
	return c, nil
}
