package tui

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/ymattw/googs"
)

type homePage struct {
	grid   *tview.Grid
	next   *tview.Button
	watch  *tview.Button
	logout *tview.Button
	games  *tview.Table
	status *tview.TextView
	hint   *tview.TextView

	ticker   *time.Ticker
	overview *googs.Overview
}

func newHomePage(app *App) Page {
	p := &homePage{
		grid:   tview.NewGrid(),
		next:   tview.NewButton("Next (0)"),
		watch:  tview.NewButton("Watch"),
		logout: tview.NewButton("Logout"),
		games:  tview.NewTable(),
		status: tview.NewTextView(),
		hint:   tview.NewTextView(),
		ticker: time.NewTicker(time.Second),
	}

	go func() {
		for range p.ticker.C {
			newLabel := fmt.Sprintf("Next (%d)", len(app.nextBoard))
			if newLabel != p.next.GetLabel() {
				p.Refresh(app)
				app.redraw(func() {
					p.next.SetLabel(newLabel)
					p.Render(app)
				})
			}
		}
	}()

	p.next.SetSelectedFunc(func() {
		if g := app.nextGameEntry(); g != nil {
			app.switchToNewGamePage(g.ID, "")
		}
	})
	p.watch.SetSelectedFunc(func() {
		app.switchToPage("watch")
	})
	p.logout.SetSelectedFunc(logoutFunc(app))

	navbar := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(nil, 0, 1, false). // left spacer
		AddItem(p.next, 10, 0, false).
		AddItem(nil, 1, 0, false). // gap
		AddItem(p.watch, 10, 0, false).
		AddItem(nil, 1, 0, false). // gap
		AddItem(p.logout, 10, 0, false)

	p.games.SetSelectable(true, false).
		SetBorder(true).
		SetTitleAlign(tview.AlignCenter)
	p.status.SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter).
		SetTextColor(Styles.TertiaryTextColor)
	p.hint.SetDynamicColors(true).
		SetTextColor(Styles.SecondaryTextColor).
		SetTextAlign(tview.AlignCenter).
		SetText(keyHints([]string{"↓↑jk select", "CR connect"}))

	// Center align the game table and bottom hint in a 4x1 grid
	p.grid.SetRows(1, 0, 1, 1)
	p.grid.SetColumns(0)
	// Row 0: navbar
	p.grid.AddItem(navbar, 0, 0, 1, 1, 1, 0, false)
	// Row 1: games table
	p.grid.AddItem(p.games, 1, 0, 1, 1, 10, 60, true)
	// Row 2: status
	p.grid.AddItem(p.status, 2, 0, 1, 1, 1, 0, false)
	// Row 3: hint
	p.grid.AddItem(p.hint, 3, 0, 1, 1, 1, 0, false)

	p.setupKeys(app)
	return p
}

func (p *homePage) Root() tview.Primitive {
	return p.grid
}

func (p *homePage) Focusables() []tview.Primitive {
	return []tview.Primitive{p.games, p.next, p.watch, p.logout}
}

func (p *homePage) Refresh(app *App) error {
	if !app.client.LoggedIn() {
		return nil
	}

	ov, err := app.client.Overview()
	if err != nil {
		app.error("Refresh home page %v", err)
		return err
	}
	p.overview = ov
	return nil
}

func (p *homePage) Render(app *App) {
	if p.overview == nil {
		return
	}

	p.status.SetText(fmt.Sprintf("You have %d active games", len(p.overview.ActiveGames)))
	p.games.Clear()
	p.games.Select(-1, -1)
	p.games.SetTitle(fmt.Sprintf(" Active Games (%d) ", len(p.overview.ActiveGames)))

	// Headers
	headers := []string{"#", "Move", "Game", "Flags", "Opponent", "Clock", "Size"}
	for col, h := range headers {
		p.games.SetCell(0, col, tview.NewTableCell(h).SetSelectable(false))
	}

	// Rows
	for i, g := range p.overview.ActiveGames {
		p.games.SetCell(i+1, 0, tview.NewTableCell(fmt.Sprintf("%d", i+1)))
		p.games.SetCell(i+1, 1, tview.NewTableCell(fmt.Sprintf("%3d", len(g.Moves))))
		p.games.SetCell(i+1, 2, tview.NewTableCell(trimString(g.GameName, 30)))
		handicap := cond(g.Handicap > 0, "🤏", "")
		private := cond(g.Private, "🔒", "")
		p.games.SetCell(i+1, 3, tview.NewTableCell(handicap+private))
		p.games.SetCell(i+1, 4, tview.NewTableCell(g.Opponent(app.client.UserID).String()))
		turn := cond(g.Clock.CurrentPlayerID == g.Players.Black.ID, googs.PlayerBlack, googs.PlayerWhite)
		p.games.SetCell(i+1, 5, tview.NewTableCell(g.Clock.ComputeClock(&g.TimeControl, turn).String()))
		p.games.SetCell(i+1, 6, tview.NewTableCell(fmt.Sprintf("%dx%d ", g.BoardSize(), g.BoardSize())))

		if g.IsMyTurn(app.client.UserID) {
			for col := range headers {
				p.games.GetCell(i+1, col).SetTextColor(Styles.TertiaryTextColor)
			}
		}
	}

	p.games.SetSelectedFunc(func(row, _ int) {
		if len(p.overview.ActiveGames) < 1 {
			return
		}
		selected := p.overview.ActiveGames[row-1]
		p.status.SetText("Connecting to " + selected.URL() + " ...")
		app.switchToNewGamePage(selected.GameID, "")
	})
}

func (p *homePage) Leave(app *App) {
	app.tui.Stop()
}

func (p *homePage) setupKeys(app *App) {
	p.games.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'r':
			app.loading(
				func() error { return p.Refresh(app) },
				func() { p.Render(app) },
			)
			return nil
		}
		return event
	})
}
