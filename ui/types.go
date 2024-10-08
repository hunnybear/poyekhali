package ui

import (
	"errors"

	"github.com/gdamore/tcell/v2"
)

type cell struct {
	x int
	y int
}

func (c cell) draw(s tcell.Screen, style tcell.Style) {
	draw(s, c.x, c.y, ' ', style)
}

func (c cell) drawChar(s tcell.Screen, char rune, style tcell.Style) {
	draw(s, c.x, c.y, char, style)
}

type BorderWidths struct {
	top    int
	right  int
	bottom int
	left   int
}

func (bw BorderWidths) asSlice() [4]int {
	return [4]int{bw.top, bw.right, bw.bottom, bw.left}
}

func (bw BorderWidths) max() int {
	return max(bw.top, bw.right, bw.bottom, bw.left)
}

func (bw BorderWidths) min() int {
	return min(bw.top, bw.right, bw.bottom, bw.left)
}

type BorderStyle struct {
	widths    BorderWidths
	drawStyle tcell.Style
}

func NewBorderWidths(widths ...int) (BorderWidths, error) {
	if len(widths) == 1 {
		return BorderWidths{widths[0], widths[0], widths[0], widths[0]}, nil
	} else if len(widths) == 4 {
		return BorderWidths{widths[0], widths[1], widths[2], widths[3]}, nil
	} else {
		return BorderWidths{1, 1, 1, 1}, errors.New("NewBorderWidts requires either 1 or 4 ints")
	}
}
