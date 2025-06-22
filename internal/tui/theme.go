package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	Styles       = &tview.Styles
	StyleDefault = tcell.StyleDefault.Background(solarizedBase03).Foreground(solarizedBase0)
)

var (
	solarizedBase03  = tcell.NewHexColor(0x002b36)
	solarizedBase02  = tcell.NewHexColor(0x073642)
	solarizedBase01  = tcell.NewHexColor(0x586e75)
	solarizedBase00  = tcell.NewHexColor(0x657b83)
	solarizedBase0   = tcell.NewHexColor(0x839496)
	solarizedBase1   = tcell.NewHexColor(0x93a1a1)
	solarizedBase2   = tcell.NewHexColor(0xeee8d5)
	solarizedBase3   = tcell.NewHexColor(0xfdf6e3)
	solarizedYellow  = tcell.NewHexColor(0xb58900)
	solarizedOrange  = tcell.NewHexColor(0xcb4b16)
	solarizedRed     = tcell.NewHexColor(0xdc322f)
	solarizedMagenta = tcell.NewHexColor(0xd33682)
	solarizedViolet  = tcell.NewHexColor(0x6c71c4)
	solarizedBlue    = tcell.NewHexColor(0x268bd2)
	solarizedCyan    = tcell.NewHexColor(0x2aa198)
	solarizedGreen   = tcell.NewHexColor(0x859900)
)

func init() {
	Styles.PrimitiveBackgroundColor = solarizedBase03
	Styles.ContrastBackgroundColor = solarizedBase02     // Background color for contrasting elements, eg. inactive buttons
	Styles.MoreContrastBackgroundColor = solarizedBase01 // Background color for even more contrasting elements
	Styles.BorderColor = solarizedBase00
	Styles.TitleColor = solarizedBase00
	Styles.GraphicsColor = solarizedBase01
	Styles.PrimaryTextColor = solarizedBase0
	Styles.SecondaryTextColor = solarizedBase1
	Styles.TertiaryTextColor = solarizedBase00 // Tertiary text (e.g. subtitles, notes)
	Styles.InverseTextColor = solarizedBase3
	Styles.ContrastSecondaryTextColor = solarizedCyan // Secondary text on ContrastBackgroundColor-colored backgrounds
}
