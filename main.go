// Package main offers a OGS client to play on the terminal.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/ymattw/googs"

	"github.com/ymattw/gote/internal/tui"
)

var (
	clientID     = flag.String("c", "", "client ID")
	clientSecret = flag.String("s", "", "client secret")
)

func main() {
	flag.Parse()
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
	if err == nil {
		return c, nil
	}
	if !os.IsNotExist(err) {
		return nil, err
	}
	if *clientID == "" {
		return nil, fmt.Errorf("new login needed but clientID (-c) and clientSecret (-s) are not specfied")
	}
	return googs.NewClient(*clientID, *clientSecret), nil
}
