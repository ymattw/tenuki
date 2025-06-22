package tui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/ymattw/googs"
)

const (
	// Full-width characters for stones and grid, best choices by far.
	GridChar       = '〸'
	HoshiChar      = '＊'
	BlackStone     = '⚫'
	WhiteStone     = '⚪'
	DeadBlackStone = '◾'
	DeadWhiteStone = '◽'
)

var (
	hoshiPoints = map[int][]googs.OriginCoordinate{
		9: {
			{X: 2, Y: 2}, {X: 2, Y: 6},
			{X: 4, Y: 4},
			{X: 6, Y: 2}, {X: 6, Y: 6},
		},
		13: {
			{X: 3, Y: 3}, {X: 3, Y: 9},
			{X: 6, Y: 6},
			{X: 9, Y: 3}, {X: 9, Y: 9},
		},
		19: {
			{X: 3, Y: 3}, {X: 3, Y: 9}, {X: 3, Y: 15},
			{X: 9, Y: 3}, {X: 9, Y: 9}, {X: 9, Y: 15},
			{X: 15, Y: 3}, {X: 15, Y: 9}, {X: 15, Y: 15},
		},
	}
)

type Stone int

const (
	Empty Stone = iota
	Black
	White
)

type Cell struct {
	stone      Stone
	isLastMove bool
	isHoshi    bool
	isRemoval  bool
}

func newCell(g *googs.GameState, row, col int) Cell {
	isHoshi := false
	hPoints := hoshiPoints[g.BoardSize()]
	for _, h := range hPoints {
		if h.X == col && h.Y == row {
			isHoshi = true
		}
	}
	return Cell{
		stone:      Stone(g.Board[row][col]),
		isLastMove: g.LastMove.X == col && g.LastMove.Y == row,
		isHoshi:    isHoshi,
		isRemoval:  g.Removal[row][col] == 1,
	}
}

func (c Cell) content() rune {
	if c.stone == Empty && c.isHoshi {
		return HoshiChar
	}
	if c.stone == Black && c.isRemoval {
		return DeadBlackStone
	}
	if c.stone == White && c.isRemoval {
		return DeadWhiteStone
	}
	return map[Stone]rune{
		Empty: GridChar,
		Black: BlackStone,
		White: WhiteStone,
	}[c.stone]
}

func (c Cell) foreground(theme string) tcell.Color {
	return boardThemes[theme].GridFG
}

func (c Cell) background(theme string) tcell.Color {
	bg := boardThemes[theme].BoardBG

	if c.isLastMove && c.stone == Black && !c.isRemoval {
		bg = boardThemes[theme].LastBlackBG
	} else if c.isLastMove && c.stone == White && !c.isRemoval {
		bg = boardThemes[theme].LastWhiteBG
	}
	return tcell.NewHexColor(int32(bg))
}

func colLabel(col int) rune {
	letter := 'Ａ' + rune(col) // Full-width Latin capital A
	if col >= 8 {
		letter += 1
	}
	return letter
}

// Board layout:
//
//	  ＡＢＣＤＥＦＧＨＪ
//	9 〸〸〸〸〸〸〸〸〸 9
//	8 〸〸〸〸〸〸〸〸〸 8
//	7 〸〸⚪〸〸⚫＊〸〸 7
//	6 〸〸〸〸〸〸〸〸〸 6
//	5 〸〸〸＊〸〸〸〸〸 5
//	4 〸〸〸〸〸〸⚪〸〸 4
//	3 〸〸⚫〸〸⚫⚪〸〸 3
//	2 〸〸〸〸〸〸〸〸〸 2
//	1 〸〸〸〸〸〸〸〸〸 1
//	  ＡＢＣＤＥＦＧＨＪ
func (p *gamePage) drawBoard(screen tcell.Screen, x, y int) (int, int, int, int) {
	size := p.gameState.BoardSize()
	whoseTurn := p.game.WhoseTurn(p.gameState)

	// Top coordinate labels (A, B, C, ... skipping I)
	for c := 0; c < size; c++ {
		// NOTE: 3-char offset for row numbers on the left, label runes
		// are Full-width.
		screen.SetContent(x+3+c*2, y, colLabel(c), nil, StyleDefault)
	}

	for row := 0; row < size; row++ {
		// Left side coordinate label (19, 18, .., 1) and a space
		left := fmt.Sprintf("%2d ", size-row)
		for i, r := range left {
			screen.SetContent(x+i, y+1+row, r, nil, StyleDefault)
		}

		for col := 0; col < size; col++ {
			cell := newCell(p.gameState, row, col)
			style := StyleDefault.
				Foreground(cell.foreground(p.boardTheme)).
				Background(cell.background(p.boardTheme))
			// Cursor use current shape in cell with reversed fg
			if col == p.cursor.X && row == p.cursor.Y {
				color := cond(whoseTurn == googs.PlayerBlack, tcell.ColorBlack, tcell.ColorWhite)
				style = style.Background(color)
			}
			// NOTE: cell runes are Full-width.
			screen.SetContent(x+3+col*2, y+1+row, cell.content(), nil, style)
		}

		// A space and right side coordinate label (19, 18, .., 1)
		right := fmt.Sprintf(" %-2d", size-row)
		for i, r := range right {
			screen.SetContent(x+3+size*2+i, y+1+row, r, nil, StyleDefault)
		}
	}

	// Bottom coordinate labels (A, B, C, ... skipping I)
	for c := 0; c < size; c++ {
		// NOTE: 3-char offset for row numbers on the left, label runes
		// are Full-width.
		screen.SetContent(x+3+c*2, y+1+size, colLabel(c), nil, StyleDefault)
	}
	screen.Show()
	return x, y, size*2 + 6, size + 2
}
