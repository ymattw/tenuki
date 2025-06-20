package tui

import (
	"fmt"
	"strings"

	"github.com/rivo/uniseg"
)

// Equivalent to Python `return x if b else y`
func cond[T any](b bool, x, y T) T {
	if b {
		return x
	}
	return y
}

// Trim a unicode string up to given max display width
func trimString(s string, maxWidth int) string {
	var b strings.Builder
	width := 0
	graphemes := uniseg.NewGraphemes(s)
	for graphemes.Next() {
		g := graphemes.Str()
		width += uniseg.StringWidth(g)
		if width > maxWidth {
			break
		}
		b.WriteString(g)
	}

	return b.String()
}

func keyHint(desc string) string {
	parts := strings.SplitN(desc, " ", 2)
	if len(parts) > 1 {
		// "CR play" => "[::r]CR[::-] play"
		return fmt.Sprintf("[::r]%s[::-] %s", parts[0], parts[1])
	}
	// "Pass" => "[::r]P[::-]ass"
	return fmt.Sprintf("[::r]%c[::-]%s", desc[0], desc[1:])
}

func keyHints(descs []string) string {
	var hints []string
	descs = append(descs, commonKeyDescriptions...)
	for _, desc := range descs {
		hints = append(hints, keyHint(desc))
	}
	return strings.Join(hints, " ⋅ ")
}
