package tui

import (
	"fmt"
	"os"

	"github.com/ymattw/tenuki/internal/config"
)

func logoutFunc(app *App) func() {
	return func() {
		app.popUp(fmt.Sprintf("Logout %s? Client ID/secret and username/password needed next time.", app.client.Username),
			[]string{"Logout", "Quit instead", "Cancel"},
			map[string]func(){
				"Logout": func() {
					os.Remove(config.SecretPath(app.client.Username))
					app.tui.Stop()
				},
				"Quit instead": func() {
					app.tui.Stop()
				},
			})
	}
}
