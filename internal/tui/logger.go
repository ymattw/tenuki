package tui

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (app *App) initLogger() {
	app.logger = tview.NewTextView()
	app.logger.SetDynamicColors(true).
		SetScrollable(true).
		SetTextStyle(StyleDefault.Background(Styles.ContrastBackgroundColor)).
		SetBackgroundColor(Styles.ContrastBackgroundColor).
		SetBorder(true).
		SetTitle(" Logs (Esc/q to close) ").
		SetTitleAlign(tview.AlignCenter)
	app.logger.SetChangedFunc(func() { app.tui.Draw(); app.logger.ScrollToEnd() })
	app.logger.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyESC || event.Rune() == 'q' {
			app.hideLogger()
			return nil
		}
		return event
	})
}

func (app *App) showLogger() {
	grid := tview.NewGrid().
		SetRows(-3, -1).
		SetColumns(-1, -60, -1).
		AddItem(app.logger, 1, 1, 1, 1, 0, 0, true)
	app.root.AddPage("logger", grid, true, true)
	app.tui.SetFocus(app.logger)
}

func (app *App) hideLogger() {
	app.root.RemovePage("logger")
}

func (app *App) info(format string, args ...any) {
	fmt.Fprintln(app.logger, fmt.Sprintf(time.Now().Format("0102 15:04:05")+" "+format, args...))
}

func (app *App) warn(format string, args ...any) {
	fmt.Fprintln(app.logger, fmt.Sprintf("[orange]"+time.Now().Format("0102 15:04:05")+" "+format+"[-]", args...))
}

func (app *App) error(format string, args ...any) {
	fmt.Fprintln(app.logger, fmt.Sprintf("[red]"+time.Now().Format("0102 15:04:05")+" "+format+"[-]", args...))
}

func (app *App) debug(format string, args ...any) {
	fmt.Fprintln(app.logger, fmt.Sprintf("[::d]"+time.Now().Format("0102 15:04:05")+" "+format+"[::-]", args...))
}
