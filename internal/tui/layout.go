package tui

import (
	"time"

	"github.com/rivo/tview"
)

func (app *App) popUp(message string, buttons []string, callbacks map[string]func()) *tview.Modal {
	returnPage, _ := app.root.GetFrontPage()
	returnFocus := app.tui.GetFocus()
	popupPage := returnPage + "-popup"
	modal := tview.NewModal().
		SetText(message).
		AddButtons(buttons).
		SetDoneFunc(func(_ int, label string) {
			app.root.RemovePage(popupPage)
			app.root.SwitchToPage(returnPage)
			app.tui.SetFocus(returnFocus)
			if cb := callbacks[label]; cb != nil {
				cb()
			}
		})
	app.root.AddPage(popupPage, modal, false, true)
	app.tui.SetFocus(modal)
	return modal
}

func (app *App) confirm(message string, callback func()) {
	app.popUp(message, []string{"Yes", "No"}, map[string]func(){"Yes": callback})
}

func (app *App) loading(refresh func() error, render func()) {
	pageName := "loading-page"
	done := make(chan struct{})

	// Do background work
	go func() {
		err := refresh()
		close(done)
		app.root.RemovePage(pageName) // Safe to do even when not shown

		app.redraw(func() {
			if err != nil {
				app.popUp("[red]"+err.Error(), []string{"OK"}, nil)
			} else {
				render()
			}
		})
	}()

	// Show loading modal after 500ms if not done yet
	timer := time.NewTimer(500 * time.Millisecond)
	select {
	case <-timer.C:
		app.redraw(func() {
			modal := tview.NewModal().SetText("Still refreshing data ...")
			app.root.AddPage(pageName, modal, false, true)
		})
	case <-done:
		timer.Stop()
	}
}
