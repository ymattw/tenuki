package tui

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/ymattw/googs"
)

type gamePage struct {
	grid    *tview.Grid
	next    *tview.Button
	home    *tview.Button
	watch   *tview.Button
	logout  *tview.Button
	title   *tview.TextView
	bPlayer *tview.TextView
	wPlayer *tview.TextView
	board   *tview.Box
	status  *tview.TextView
	hint    *tview.TextView
	chat    *tview.Table
	message *tview.InputField

	returnPage string
	gameID     int64            // Orignal input
	game       *googs.Game      // Loaded game
	gameState  *googs.GameState // Loaded game state
	cursor     *googs.OriginCoordinate
	clock      *googs.Clock
	ticker     *time.Ticker
	chats      []*googs.GameChatLine
	chatsLock  sync.Mutex
}

func newGamePage(app *App, gameID int64, returnPage string) Page {
	p := &gamePage{
		grid:    tview.NewGrid(),
		next:    tview.NewButton("Next (0)"),
		home:    tview.NewButton("Home"),
		watch:   tview.NewButton("Watch"),
		logout:  tview.NewButton("Logout"),
		title:   tview.NewTextView(),
		bPlayer: tview.NewTextView(),
		board:   tview.NewBox(),
		wPlayer: tview.NewTextView(),
		status:  tview.NewTextView(),
		hint:    tview.NewTextView(),
		chat:    tview.NewTable(),
		message: tview.NewInputField(),

		returnPage: returnPage,
		gameID:     gameID,
		game:       &googs.Game{},      // Avoid nil deference
		gameState:  &googs.GameState{}, // Avoid nil deference
		cursor:     &googs.OriginCoordinate{},
		ticker:     time.NewTicker(time.Second),
	}

	// Update Next label and clock displays every second, keep it simple
	// instead of dynamically reset
	go func() {
		for range p.ticker.C {
			updated := p.updatePlayers()
			newLabel := fmt.Sprintf("Next (%d)", len(app.nextBoard))
			if updated || newLabel != p.next.GetLabel() {
				app.redraw(func() {
					p.next.SetLabel(newLabel)
				})
			}
		}
	}()

	p.next.SetSelectedFunc(func() {
		if g := app.nextGameEntry(); g != nil {
			p.Leave(app)
			app.switchToNewGamePage(g.ID, "home") // Current page is removed
		}
	})
	p.home.SetSelectedFunc(func() {
		app.switchToPage("home")
	})
	p.watch.SetSelectedFunc(func() {
		app.switchToPage("watch")
	})
	p.logout.SetSelectedFunc(logoutFunc(app))

	p.title.SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter)
	// bPlayer and wPlayer use secondary color to not conflict with
	// highlighted widget in focus.
	p.bPlayer.SetDynamicColors(true).
		SetTextColor(Styles.SecondaryTextColor).
		SetTextAlign(tview.AlignCenter).
		SetTitleColor(Styles.SecondaryTextColor).
		SetTitleAlign(tview.AlignCenter).
		SetBorder(true).
		SetBorderColor(Styles.SecondaryTextColor)
	p.wPlayer.SetDynamicColors(true).
		SetTextColor(Styles.SecondaryTextColor).
		SetTextAlign(tview.AlignCenter).
		SetTitleColor(Styles.SecondaryTextColor).
		SetTitleAlign(tview.AlignCenter).
		SetBorder(true).
		SetBorderColor(Styles.SecondaryTextColor)
	p.board.SetBorder(true).
		SetFocusFunc(func() { p.board.SetBorderColor(Styles.PrimaryTextColor) }).
		SetBlurFunc(func() { p.board.SetBorderColor(Styles.SecondaryTextColor) })
	p.status.SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter)
	p.hint.SetDynamicColors(true).
		SetTextColor(Styles.SecondaryTextColor).
		SetTextAlign(tview.AlignCenter)

	p.chat.SetSelectable(true, true).
		SetBorder(true).
		SetTitle(" Chat ").
		SetTitleAlign(tview.AlignCenter).
		SetTitleColor(Styles.SecondaryTextColor).
		SetBorderColor(Styles.SecondaryTextColor).
		SetFocusFunc(func() { p.chat.SetBorderColor(Styles.PrimaryTextColor) }).
		SetBlurFunc(func() { p.chat.SetBorderColor(Styles.SecondaryTextColor) })

	p.message.
		SetFieldBackgroundColor(Styles.PrimitiveBackgroundColor).
		SetPlaceholder("> Enter message ...").
		SetPlaceholderStyle(tcell.StyleDefault.Background(Styles.PrimitiveBackgroundColor)).
		SetDoneFunc(func(key tcell.Key) {
			if key == tcell.KeyEnter {
				app.client.GameChat(p.gameID, p.gameState.MoveNumber, p.message.GetText())
				p.message.SetText("")
			}
		}).
		SetFocusFunc(func() {
			app.redraw(func() {
				p.message.SetFieldTextColor(Styles.PrimaryTextColor).
					SetPlaceholder("")
			})
		},
		).
		SetBlurFunc(func() {
			app.redraw(func() {
				p.message.SetFieldTextColor(Styles.SecondaryTextColor).
					SetPlaceholder("> Enter message ...")
			})
		})

	return p
}

func (p *gamePage) Root() tview.Primitive {
	return p.grid
}

func (p *gamePage) Focusables() []tview.Primitive {
	return []tview.Primitive{p.board, p.chat, p.message, p.next, p.home, p.watch, p.logout}
}

func (p *gamePage) Refresh(app *App) error {
	if err := p.refreshGame(app); err != nil {
		return err
	}
	if err := p.refreshGameState(app); err != nil {
		return err
	}
	if err := app.client.ChatJoin(p.game.GameID); err != nil {
		return err
	}
	if err := app.client.GameConnect(p.game.GameID); err != nil {
		return err
	}
	return nil
}

func (p *gamePage) resetLayout() {
	navbar := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(nil, 0, 1, false). // left spacer
		AddItem(p.next, 10, 0, false).
		AddItem(nil, 1, 0, false). // gap
		AddItem(p.home, 10, 0, false).
		AddItem(nil, 1, 0, false). // gap
		AddItem(p.watch, 10, 0, false).
		AddItem(nil, 1, 0, false). // gap
		AddItem(p.logout, 10, 0, false)

	bPlayerFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(p.bPlayer, 7, 1, false)
	wPlayerFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(p.wPlayer, 7, 1, false)

	// Align the elements in a 11x7 grid
	p.grid.SetRows(
		1,                    // navbar
		-1,                   // spacer
		1,                    // title
		-1,                   // spacer
		p.game.BoardSize()+2, // board with labels
		-1,                   // spacer
		1,                    // status
		1,                    // hint
		-3,                   // chat
		1,                    // message
		-1,                   // spacer
	)
	p.grid.SetColumns(
		-1,                       // spacer
		18,                       // black player
		1,                        // gap
		3+p.game.BoardSize()*2+3, // board with labels
		1,                        // gap
		18,                       // white player
		-1,                       // spacer
	)
	// Row 0: navbar, span 7 columns
	p.grid.AddItem(navbar, 0, 0, 1, 7, 1, 0, false)
	// Row 1: spacer
	// Row 2: spacer, title (span 5 columns), spacer
	p.grid.AddItem(p.title, 2, 1, 1, 5, 1, 0, false)
	// Row 3: spacer
	// Row 4: spacer, bPlayer, gap, board, gap, wPlayer, spacer
	p.grid.AddItem(bPlayerFlex, 4, 1, 1, 1, 5 /*minWidth*/, 0, false)
	p.grid.AddItem(p.board, 4, 3, 1, 1, 0, 0, true)
	p.grid.AddItem(wPlayerFlex, 4, 5, 1, 1, 5 /*minWidth*/, 0, false)
	// Row 5: spacer
	// Row 6: spacer, status (5 columns), spacer
	p.grid.AddItem(p.status, 6, 1, 1, 5, 1, 0, false)
	// Row 7: spacer, hint (5 columns), spacer
	p.grid.AddItem(p.hint, 7, 1, 1, 5, 1, 0, false)
	// Row 8: chat (7 columns)
	p.grid.AddItem(p.chat, 8, 0, 1, 7, 1, 50, false)
	// Row 9: message (7 columns)
	p.grid.AddItem(p.message, 9, 0, 1, 7, 1, 50, false)
}

func (p *gamePage) gameTitle() string {
	speed := map[string]string{
		"blitz":          "⚡",
		"rapid":          "⏩",
		"live":           "⏱️",
		"correspondence": "🐢",
	}[p.game.TimeControl.Speed]

	shortRule := map[string]string{
		"aga":      "AGA",
		"chinese":  "CN",
		"ing":      "ING",
		"japanese": "JP",
		"korean":   "KR",
		"nz":       "NZ",
	}[strings.ToLower(p.game.Rules)]

	ranked := cond(p.game.Ranked, "ranked", "unranked")
	private := cond(p.game.Private, "🔒", "")
	// Note '❶' != '⓿' + 1
	handicap := rune(cond(p.game.Handicap > 0, '❶'+p.game.Handicap-1, '⓿'))
	return fmt.Sprintf("#%d %s %s %s %s, %c +%.1f, %s %s",
		p.game.GameID, trimString(p.game.GameName, 30), speed, shortRule, p.game.TimeControl, handicap, p.game.Komi, ranked, private)
}

func (p *gamePage) Render(app *App) {
	p.resetLayout()
	p.setupKeys(app) // p.board is dynamical

	p.title.SetText(p.gameTitle())
	p.clock = &p.game.Clock // Initial game clock
	p.updatePlayers()
	p.updateStatusAndHint(app)

	p.board.SetDrawFunc(func(screen tcell.Screen, x, y, width, height int) (int, int, int, int) {
		return p.drawBoard(screen, x, y)
	})

	app.client.OnGameData(p.game.GameID, func(g *googs.Game) {
		p.refreshGame(app)
		app.redraw(func() { p.updateStatusAndHint(app) })
	})

	app.client.OnGamePhase(p.game.GameID, func(phase googs.GamePhase) {
		p.refreshGame(app)
		p.refreshGameState(app) // gameState has removal and outcome
		app.redraw(func() { p.updateStatusAndHint(app) })
	})

	app.client.OnGameRemovedStones(p.game.GameID, func(r *googs.RemovedStones) {
		p.refreshGame(app)
		p.refreshGameState(app) // for dead stones drawing
		app.redraw(func() { p.updateStatusAndHint(app) })
	})

	app.client.OnGameRemovedStonesAccepted(p.game.GameID, func(r *googs.RemovedStonesAccepted) {
		var who string
		if p.game.IsMyGame(app.client.UserID) {
			who = cond(r.PlayerID == app.client.UserID, "You have", "Opponent has")
		} else {
			who = cond(r.PlayerID == p.game.BlackPlayerID, "Black has", "White has")
		}
		app.redraw(func() {
			p.status.SetText("[red]" + who + " accepted stone removal[-]")
			if r.Phase == googs.FinishedPhase {
				p.status.SetText("[green]" + r.Result() + "[-]")
				p.hint.SetText(keyHints(nil))
			}
		})
	})

	app.client.OnMove(p.game.GameID, func(m *googs.GameMove) {
		p.refreshGameState(app)
		app.redraw(func() { p.updateStatusAndHint(app) })
	})

	app.client.OnClock(p.game.GameID, func(c *googs.Clock) {
		// DP("onClock: %s", formatObject(c))
		p.clock = c
	})

	app.client.OnGameChat(p.game.GameID, func(chat *googs.GameChat) {
		p.chatsLock.Lock()
		p.chats = insertSortedChats(p.chats, &chat.Line)
		p.updateChatTable()
		p.chatsLock.Unlock()
		app.redraw(nil)
	})
}

func insertSortedChats(lines []*googs.GameChatLine, newLine *googs.GameChatLine) []*googs.GameChatLine {
	for _, line := range lines {
		if line.ChatID == newLine.ChatID {
			return lines // duplicate found, skip insert
		}
	}
	index := sort.Search(len(lines), func(i int) bool {
		return lines[i].Date.After(newLine.Date.Time)
	})
	lines = append(lines, nil)           // extend slice by 1
	copy(lines[index+1:], lines[index:]) // shift elements right
	lines[index] = newLine               // insert at index
	return lines
}

func (p *gamePage) refreshGame(app *App) error {
	g, err := app.client.Game(p.gameID)
	if err != nil {
		return err
	}
	p.game = g
	return nil
}

func (p *gamePage) updatePlayer(t *tview.TextView, c googs.PlayerColor) bool {
	// Use a blinking dot to indicate who is on turn
	// TODO: online status
	title := cond(p.game.WhoseTurn(p.gameState) == c,
		fmt.Sprintf(" %s [::l]•[-] ", c),
		fmt.Sprintf(" %s ", c))
	clock := p.clock.ComputeClock(&p.game.TimeControl, c)
	// FIXME: support more clock systems
	style := cond(clock != nil && clock.SuddenDeath, "[red]", "")
	player := cond(c == googs.PlayerBlack, p.game.BlackPlayer(), p.game.WhitePlayer())
	text := fmt.Sprintf("\n%s\n\n%s%s[-]", player, style, clock)

	if title == t.GetTitle() && text == t.GetText(false) {
		return false
	}
	t.SetTitle(title)
	t.SetText(text)
	return true
}

func (p *gamePage) updatePlayers() bool {
	if p.game.Phase != googs.PlayPhase {
		return false
	}
	b := p.updatePlayer(p.bPlayer, googs.PlayerBlack)
	w := p.updatePlayer(p.wPlayer, googs.PlayerWhite)
	return b && w
}

func (p *gamePage) updateStatusAndHint(app *App) {
	isMyGame := p.game.IsMyGame(app.client.UserID)
	switch p.game.Phase {
	case googs.PlayPhase:
		p.status.SetText(p.game.Status(p.gameState, app.client.UserID))
		p.hint.SetText(cond(p.gameState.IsMyTurn(app.client.UserID),
			keyHints([]string{"←↓↑→hjkl move cursor", "CR play", "Pass", "Resign"}),
			cond(isMyGame,
				keyHints([]string{"Resign"}),
				keyHints(nil))))
	case googs.StoneRemovalPhase:
		p.status.SetText(fmt.Sprintf("%s phase", p.game.Phase))
		p.hint.SetText(cond(isMyGame,
			keyHints([]string{"Accept"}),
			keyHints(nil)))
	case googs.FinishedPhase:
		p.status.SetText("[green]" + p.game.Result() + "[-]")
		p.hint.SetText(keyHints(nil))
	}
}

func (p *gamePage) updateChatTable() {
	// No mutex lock needed
	chatCount := len(p.chats)
	p.chat.Clear()

	for row := 0; row < chatCount; row++ {
		line := p.chats[row]
		player := &googs.Player{
			Professional: line.Professional != 0,
			Rank:         line.Ranking,
			Username:     line.Username,
		}
		p.chat.SetCell(row, 0, tview.NewTableCell(line.Date.Format(time.DateTime)).SetTextColor(Styles.SecondaryTextColor))
		p.chat.SetCell(row, 1, tview.NewTableCell(fmt.Sprintf("%d", line.MoveNumber)).SetTextColor(Styles.SecondaryTextColor).SetAlign(tview.AlignRight))
		p.chat.SetCell(row, 2, tview.NewTableCell(player.String()).SetTextColor(Styles.TertiaryTextColor))
		p.chat.SetCell(row, 3, tview.NewTableCell(strings.TrimSpace(line.Body)).SetTextColor(Styles.SecondaryTextColor))
	}
	if chatCount > 0 {
		p.chat.Select(chatCount-1, 3)
		p.chat.SetSelectedStyle(p.chat.GetCell(chatCount-1, 3).Style.Foreground(Styles.PrimaryTextColor))
	}
}

func (p *gamePage) refreshGameState(app *App) error {
	g, err := app.client.GameState(p.gameID)
	if err != nil {
		return err
	}
	p.gameState = g

	p.cursor.X, p.cursor.Y = -1, -1 // Hide cursor
	if p.gameState.IsMyTurn(app.client.UserID) {
		if p.gameState.LastMove.IsPass() {
			p.cursor.X = p.gameState.BoardSize() / 2
			p.cursor.Y = p.gameState.BoardSize() / 2
		} else {
			p.cursor.X = p.gameState.LastMove.X
			p.cursor.Y = p.gameState.LastMove.Y
		}
	}

	return nil
}

func (p *gamePage) Leave(app *App) {
	// Disconnect game, stop refresh
	p.ticker.Stop()
	app.client.GameDisconnect(p.game.GameID)
	app.removePage(fmt.Sprintf("%d", p.game.GameID))
	app.switchToPage(p.returnPage)
}

func (p *gamePage) setupKeys(app *App) {
	size := p.gameState.BoardSize()
	p.board.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		myTurn := p.gameState.IsMyTurn(app.client.UserID)

		if event.Key() == tcell.KeyLeft || event.Rune() == 'h' {
			if myTurn && p.cursor.X > 0 {
				p.cursor.X--
			}
			return nil
		} else if event.Key() == tcell.KeyDown || event.Rune() == 'j' {
			if myTurn && p.cursor.Y < size-1 {
				p.cursor.Y++
			}
			return nil
		} else if event.Key() == tcell.KeyUp || event.Rune() == 'k' {
			if myTurn && p.cursor.Y > 0 {
				p.cursor.Y--
			}
			return nil
		} else if event.Key() == tcell.KeyRight || event.Rune() == 'l' {
			if myTurn && p.cursor.X < size-1 {
				p.cursor.X++
			}
			return nil
		} else if event.Key() == tcell.KeyEnter {
			if myTurn && p.cursor.X != -1 && p.cursor.Y != -1 && p.gameState.Board[p.cursor.Y][p.cursor.X] == 0 {
				app.client.GameMove(p.game.GameID, p.cursor.X, p.cursor.Y)
				return nil
			}
		} else if event.Rune() == 'P' {
			if myTurn {
				app.confirm("Pass?", func() {
					app.client.PassTurn(p.game.GameID)
				})
				return nil
			}
		} else if event.Rune() == 'R' {
			if p.game.IsMyGame(app.client.UserID) {
				app.confirm("Resign?", func() {
					app.client.GameResign(p.game.GameID)
				})
				return nil
			}
		} else if event.Rune() == 'A' {
			if p.game.IsMyGame(app.client.UserID) && p.game.Phase == googs.StoneRemovalPhase {
				app.confirm("Accept stone removal?", func() {
					app.client.GameRemovedStonesAccept(p.game.GameID, p.gameState)
				})
				return nil
			}
		}

		return event
	})
}
