package tui

import (
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/ymattw/googs"
)

type App struct {
	client *googs.Client
	tui    *tview.Application
	root   *tview.Pages
	logger *tview.TextView
	pages  map[string]Page

	// Connection measurement (milliseconds)
	drift      int64
	latency    int64
	pingTicker *time.Ticker

	// Next actionable board to move on, key is gameID
	nextBoard     map[int64]*googs.GameListEntry
	currentGameID int64
}

type Page interface {
	Root() tview.Primitive         // Root view of the page
	Focusables() []tview.Primitive // Focusable primitives, [0] = default
	Refresh(*App) error            // Refresh data for redraw, must NOT touch UI
	Render(*App)                   // Update content of the widgets
	Leave(*App)                    // Clean up and switch page (when Esc pressed)
}

func NewApp(client *googs.Client) *App {
	app := &App{
		client:     client,
		tui:        tview.NewApplication(),
		root:       tview.NewPages(),
		pages:      make(map[string]Page),
		pingTicker: time.NewTicker(10 * time.Second),
		nextBoard:  make(map[int64]*googs.GameListEntry),
	}

	// Too small screen leads to tab switching focus to invisble
	// primitives and cause app to hang.
	app.tui.SetBeforeDrawFunc(func(screen tcell.Screen) bool {
		w, h := screen.Size()
		if w < 70 || h < 30 {
			msg := "Screen too small, make it at least 70x30."
			tview.Print(screen, msg, 0, 0, len(msg), tview.AlignLeft, solarizedRed)
			return true
		}
		return false
	})

	app.initLogger()
	app.info("App initialized")

	// Precreated pages which are not required to follow Page interface
	app.root.AddPage("login", newLoginPage(app, func() {
		// TODO: save to .local/state/
		if err := app.client.Save("secret.json"); err != nil {
			panic(err)
		}
		app.onLoggedIn()
	}), true, false)

	return app
}

func (app *App) addPage(name string, page Page) {
	app.pages[name] = page
	app.root.AddPage(name, page.Root(), true, false)
	app.setupCommonKeys(page)
}

// This is expected to be called only once upon logged in
func (app *App) onLoggedIn() {
	app.client.NetPing(0, 0) // Initial ping
	app.client.OnNetPong(func(drift, latency int64) {
		app.info("Server pong drift=%d latency=%d", drift, latency)
		app.drift, app.latency = drift, latency
	})
	app.client.OnActiveGame(func(g *googs.GameListEntry) {
		app.info("Active game update of #%d", g.ID)
		delete(app.nextBoard, g.ID)
		switch g.Phase {
		case googs.FinishedPhase:
		case googs.StoneRemovalPhase:
			if (g.Black.ID == app.client.UserID && g.Black.AcceptedStones == nil) ||
				(g.White.ID == app.client.UserID && g.White.AcceptedStones == nil) {
				app.nextBoard[g.ID] = g
			}
		case googs.PlayPhase:
			if g.PlayerToMove == app.client.UserID {
				app.nextBoard[g.ID] = g
			}
		}
	})

	go func() {
		for range app.pingTicker.C {
			app.client.NetPing(app.drift, app.latency)
		}
	}()

	app.addPage("home", newHomePage(app))
	app.addPage("watch", newWatchPage(app))
	app.switchToPage("home")
}

func (app *App) Run() error {
	if app.client.LoggedIn() {
		app.onLoggedIn()
	} else {
		app.root.SwitchToPage("login")
	}
	app.tui.SetRoot(app.root, true)
	app.info("App started running")
	return app.tui.Run()
}

// Always safe to call from no matter where
func (app *App) redraw(fn func()) {
	fn = cond(fn != nil, fn, func() {})
	go func() {
		app.tui.QueueUpdateDraw(fn)
	}()
}

func (app *App) switchToPage(name string) {
	app.loading(
		func() error {
			return app.pages[name].Refresh(app)
		},
		func() {
			app.pages[name].Render(app)
			app.root.SwitchToPage(name)
			if len(app.pages[name].Focusables()) > 0 {
				app.tui.SetFocus(app.pages[name].Focusables()[0])
			}
		},
	)
}

func (app *App) removePage(name string) {
	app.root.RemovePage(name)
}

func (app *App) switchToNewGamePage(gameID int64, returnPage string) {
	pageName := fmt.Sprintf("%d", gameID)
	if returnPage == "" {
		returnPage, _ = app.root.GetFrontPage()
	}
	app.addPage(pageName, newGamePage(app, gameID, returnPage))
	app.currentGameID = gameID
	app.switchToPage(pageName)
}

func (app *App) nextGameEntry() *googs.GameListEntry {
	if len(app.nextBoard) == 0 {
		return nil
	}

	gameIDs := make([]int64, 0, len(app.nextBoard))
	for id := range app.nextBoard {
		gameIDs = append(gameIDs, id)
	}
	sort.Slice(gameIDs, func(i, j int) bool {
		return gameIDs[i] < gameIDs[j]
	})

	found := false
	for _, id := range gameIDs {
		if found {
			return app.nextBoard[id]
		}
		if id == app.currentGameID {
			found = true
		}
	}
	return app.nextBoard[gameIDs[0]]
}

var commonKeyDescriptions = []string{"Home", "Next", "Watch", "quit"}

// Set up common shortcuts. Page Root() must be a Grid.
func (app *App) setupCommonKeys(p Page) {
	tabIndex := 0

	p.Root().(*tview.Grid).SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if len(p.Focusables()) == 0 {
			return event
		}

		_, onInputField := app.tui.GetFocus().(*tview.InputField)

		switch event.Key() {
		case tcell.KeyESC:
			if onInputField {
				// Refocus to the first (main) focusable widget
				app.tui.SetFocus(p.Focusables()[0])
				tabIndex = 0
			} else {
				p.Leave(app)
			}
			return nil
		case tcell.KeyTAB:
			tabIndex = (tabIndex + 1) % len(p.Focusables())
			app.tui.SetFocus(p.Focusables()[tabIndex])
			return nil
		case tcell.KeyBacktab: // Shift+Tab
			tabIndex = (tabIndex - 1 + len(p.Focusables())) % len(p.Focusables())
			app.tui.SetFocus(p.Focusables()[tabIndex])
			return nil
		}

		if onInputField {
			return event
		}

		switch event.Rune() {
		case 'D':
			app.showLogger()
			return nil
		case 'H':
			app.switchToPage("home")
			return nil
		case 'N':
			curPage, _ := app.root.GetFrontPage()
			if _, err := strconv.ParseInt(curPage, 10, 64); err == nil {
				// Dynamic game page, leave (to "home") first
				p.Leave(app)
			}
			// If no next, this actually stays at "home" (desired)
			if g := app.nextGameEntry(); g != nil {
				app.switchToNewGamePage(g.ID, "")
			}
			return nil
		case 'W':
			app.switchToPage("watch")
			return nil
		case 'q':
			p.Leave(app)
			return nil
		}
		return event
	})
}
