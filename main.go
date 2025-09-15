// Package main offers a OGS client to play on the terminal.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/ymattw/googs"

	"github.com/ymattw/tenuki/internal/config"
	"github.com/ymattw/tenuki/internal/tui"
)

var (
	showVersion = flag.Bool("V", false, "Print version and exit")
	username    = flag.String("u", "", "OGS username, only needed for switching accounts")

	// To be set by compiler via -ldflags
	buildVersion string
	buildDate    string
	buildCommit  string
)

func main() {
	flag.Parse()
	if *showVersion {
		fmt.Printf("Tenuki version: %s\nbuilt on: %s\ncommit: %s\n",
			buildVersion, buildDate, buildCommit)
		os.Exit(0)
	}

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
