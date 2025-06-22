package tui

import (
	"fmt"

	"github.com/rivo/tview"
)

func newLoginPage(app *App, callback func()) tview.Primitive {
	status := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter).
		SetTextColor(Styles.MoreContrastBackgroundColor).
		SetText("See your OGS \"Application\" or register one at\nhttps://online-go.com/oauth2/applications/")

	cField := tview.NewInputField().
		SetLabel("Client ID").
		SetFieldWidth(42).
		SetText(app.client.ClientID)
	sField := tview.NewInputField().
		SetLabel("Client Secret").
		SetFieldWidth(42).
		SetText(app.client.ClientSecret).
		SetPlaceholder("Required for confidential client").
		SetPlaceholderTextColor(Styles.MoreContrastBackgroundColor)
	uField := tview.NewInputField().
		SetLabel("Username").
		SetFieldWidth(42)
	pField := tview.NewInputField().
		SetLabel("Password").
		SetFieldWidth(42).
		SetMaskCharacter('*')

	form := tview.NewForm().
		AddFormItem(cField).
		AddFormItem(sField).
		AddFormItem(uField).
		AddFormItem(pField)
	form.SetButtonsAlign(tview.AlignCenter).
		AddButton("Submit", func() {
			app.client.ClientID = cField.GetText()
			app.client.ClientSecret = sField.GetText()
			if err := app.client.Login(uField.GetText(), pField.GetText()); err != nil {
				status.SetText(fmt.Sprintf("[red]%v[-]", err))
				app.tui.SetFocus(form.GetFormItemByLabel("Password"))
			} else {
				callback()
				status.SetText("[green]Success, switching to home page ...")
			}
		}).
		AddButton("Quit", func() {
			app.tui.Stop()
		}).
		SetTitle(" Login to OGS ").
		SetBorder(true)

	status.SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter).
		SetTextColor(Styles.MoreContrastBackgroundColor).
		SetText("See your OGS \"Application\" or register one at\nhttps://online-go.com/oauth2/applications/")

	// Center align the form and bottom status in a 3x3 grid
	grid := tview.NewGrid().
		SetRows(-1, 13, -1).
		SetColumns(-1, 60, -1).
		AddItem(form, 1, 1, 1, 1, 0, 0, true).
		AddItem(status, 2, 1, 1, 1, 0, 0, false)
	return grid
}
