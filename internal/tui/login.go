package tui

import (
	"fmt"

	"github.com/rivo/tview"
)

type loginPage struct {
	grid     *tview.Grid
	form     *tview.Form
	status   *tview.TextView
	callback func()
}

func newLoginPage(app *App, callback func()) Page {
	p := &loginPage{
		grid:     tview.NewGrid(),
		form:     tview.NewForm(),
		status:   tview.NewTextView(),
		callback: callback,
	}

	p.status.SetDynamicColors(true)
	p.form.
		SetFieldBackgroundColor(Styles.SecondaryTextColor).
		AddInputField("Username", "", 32, nil, nil).
		AddPasswordField("Password", "", 32, '*', nil).
		AddButton("Submit", func() {
			username := p.form.GetFormItemByLabel("Username").(*tview.InputField).GetText()
			password := p.form.GetFormItemByLabel("Password").(*tview.InputField).GetText()
			if err := app.client.Login(username, password); err != nil {
				p.status.SetText(fmt.Sprintf("[red]%v[-]", err))
				p.form.GetFormItemByLabel("Password").(*tview.InputField).SetText("")
				app.tui.SetFocus(p.form.GetFormItem(0))
			} else {
				callback()
				p.status.SetText("[green]Success, switching to home page ...")
			}
		}).
		SetTitle(" Login to OGS ").
		SetBorder(true)

	// Center align the form and bottom status in a 3x3 grid
	p.grid.SetRows(0, 10, 0)
	p.grid.SetColumns(0, 50, 0)
	p.grid.AddItem(p.form, 1, 1, 1, 1, 0, 0, true)
	p.grid.AddItem(p.status, 2, 1, 1, 1, 0, 0, false)
	return p
}

func (p *loginPage) Root() tview.Primitive {
	return p.grid
}

func (p *loginPage) Focusables() []tview.Primitive {
	return []tview.Primitive{p.form}
}

func (p *loginPage) Refresh(app *App) error {
	return nil
}

func (p *loginPage) Render(app *App) {
	p.status.Clear()
}

func (p *loginPage) Leave(app *App) {
	app.tui.Stop()
}
