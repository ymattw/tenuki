package tui

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/ymattw/googs"
)

type watchPage struct {
	grid   *tview.Grid
	next   *tview.Button
	home   *tview.Button
	logout *tview.Button
	games  *tview.Table
	status *tview.TextView
	hint   *tview.TextView

	ticker   *time.Ticker
	gameList *googs.GameListResponse
}

func newWatchPage(app *App) Page {
	p := &watchPage{
		grid:   tview.NewGrid(),
		next:   tview.NewButton("Next (0)"),
		home:   tview.NewButton("Home"),
		logout: tview.NewButton("Logout"),
		games:  tview.NewTable(),
		status: tview.NewTextView(),
		hint:   tview.NewTextView(),
		ticker: time.NewTicker(time.Second),
	}

	go func() {
		for range p.ticker.C {
			// Do not Refresh() here otherwise too many requests
			newLabel := fmt.Sprintf("Next (%d)", len(app.nextBoard))
			if newLabel != p.next.GetLabel() {
				app.redraw(func() {
					p.next.SetLabel(newLabel)
				})
			}
		}
	}()

	p.next.SetSelectedFunc(func() {
		if g := app.nextGameEntry(); g != nil {
			app.switchToNewGamePage(g.ID, "")
		}
	})
	p.home.SetSelectedFunc(func() {
		app.switchToPage("home")
	})
	p.logout.SetSelectedFunc(logoutFunc(app))

	navbar := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(nil, 0, 1, false). // left spacer
		AddItem(p.next, 10, 0, false).
		AddItem(nil, 1, 0, false). // gap
		AddItem(p.home, 10, 0, false).
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

func (p *watchPage) Root() tview.Primitive {
	return p.grid
}

func (p *watchPage) Focusables() []tview.Primitive {
	return []tview.Primitive{p.games, p.next, p.home, p.logout}
}

func (p *watchPage) Refresh(app *App) error {
	if !app.client.LoggedIn() {
		return nil
	}

	resp, err := app.client.GameListQuery(googs.LiveGameList, 0, 20, nil, time.Second*10)
	if err != nil {
		app.error("Refresh watch page %v", err)
		return err
	}
	p.gameList = resp
	return nil
}

func (p *watchPage) Render(app *App) {
	if p.gameList == nil {
		return
	}

	p.games.Clear()
	p.games.Select(-1, -1)
	if p.gameList == nil {
		p.status.SetText("[red]Query game list got null response[-]")
		return
	}
	p.status.SetText(fmt.Sprintf("Showing top %d of total %d %s games", len(p.gameList.Results), p.gameList.Size, p.gameList.List))

	p.games.SetTitle(fmt.Sprintf(" Live Games (%d) ", len(p.gameList.Results)))

	// Headers
	headers := []string{"Move", "Game", "Flags", "Black", "White", "Handicap", "Komi", "Phase", "Size"}
	for col, h := range headers {
		p.games.SetCell(0, col, tview.NewTableCell(h).SetSelectable(false))
	}

	// Rows
	for i, g := range p.gameList.Results {
		p.games.SetCell(i+1, 0, tview.NewTableCell(fmt.Sprintf("%3d", g.MoveNumber)))
		p.games.SetCell(i+1, 1, tview.NewTableCell(trimString(g.Name, 30)))
		bot := cond(g.BotGame, "🤖", "")
		private := cond(g.Private, "🔒", "")
		handicap := cond(g.Handicap > 0, "🤏", "")
		p.games.SetCell(i+1, 2, tview.NewTableCell(bot+private+handicap))
		p.games.SetCell(i+1, 3, tview.NewTableCell(g.Black.String()))
		p.games.SetCell(i+1, 4, tview.NewTableCell(g.White.String()))
		p.games.SetCell(i+1, 5, tview.NewTableCell(fmt.Sprintf("%d", g.Handicap)))
		p.games.SetCell(i+1, 6, tview.NewTableCell(fmt.Sprintf("%.1f", g.Komi)))
		p.games.SetCell(i+1, 7, tview.NewTableCell(string(g.Phase)))
		p.games.SetCell(i+1, 8, tview.NewTableCell(fmt.Sprintf("%dx%d ", g.Width, g.Height)))

		// Is my game
		if g.Black.ID == app.client.UserID || g.White.ID == app.client.UserID {
			for col := range headers {
				p.games.GetCell(i+1, col).SetTextColor(Styles.TertiaryTextColor)
			}
		}
	}

	p.games.SetSelectedFunc(func(row, _ int) {
		if len(p.gameList.Results) < 1 {
			return
		}
		selected := p.gameList.Results[row-1]
		p.status.SetText(fmt.Sprintf("Connecting to game %d ...", selected.ID))
		app.switchToNewGamePage(selected.ID, "")
	})
}

func (p *watchPage) Leave(app *App) {
	app.tui.Stop()
}

func (p *watchPage) setupKeys(app *App) {
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
