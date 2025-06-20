package tui

func logoutFunc(app *App) func() {
	return func() {
		app.popUp("Logout? Username/password needed next time.",
			[]string{"Logout", "Quit instead", "Cancel"},
			map[string]func(){
				// TODO: implement logout
				"Quit instead": func() { app.tui.Stop() },
			})
	}
}
