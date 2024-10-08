package ui

import (
	"errors"
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
)

var debugStyle tcell.Style

func init() {
	debugStyle = tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorRed)
}

func draw(s tcell.Screen, x, y int, char rune, style tcell.Style) {
	s.SetContent(x, y, char, nil, style)
}

func writeText(s tcell.Screen, x, y int, text string, style tcell.Style) {
	for i, char := range text {
		draw(s, x+i, y, char, style)
	}
}

func debug(s tcell.Screen, text string, printMessage ...bool) {
	return
	// Bottom row, right col
	if len(printMessage) > 0 && printMessage[0] == true {
		//fmt.Println("[[DEBUG]]: ", text)
	}
	maxCol, maxRow := s.Size()
	maxCol -= 1
	maxRow -= 1
	//tood handle too long lines
	//if len(text)+5 < maxCol {
	// add leader
	//	drawHLine(s, maxRow, 0, maxCol-len(text)-5, debugStyle)
	//}
	//fmt.Println("writing text", text)
	//writeText(s, maxCol-5-len(text), maxRow, text, debugStyle)
	//fmt.Println("drawing line", text)
	//drawHLine(s, maxRow, maxCol-5, maxCol, debugStyle)
}

func drawHLine(s tcell.Screen, y, start, end int, style tcell.Style, debugMsg ...string) error {
	if len(debugMsg) > 0 {
		return drawLine(s, cell{x: start, y: y}, cell{x: end, y: y}, style, debugMsg[0])
	} else {
		return drawLine(s, cell{x: start, y: y}, cell{x: end, y: y}, style)
	}
}

func drawVLine(s tcell.Screen, x, start, end int, style tcell.Style, debugMsg ...string) error {
	if len(debugMsg) > 0 {
		return drawLine(s, cell{x: x, y: start}, cell{x: x, y: end}, style, debugMsg[0])
	} else {
		return drawLine(s, cell{x: x, y: start}, cell{x: x, y: end}, style)
	}
}

func drawLine(s tcell.Screen, start, end cell, style tcell.Style, debugMsg ...string) error {
	var debugText string
	getDebugText := func(i int) rune {
		if i >= len(debugText) {
			return ' '
		}
		return rune(debugText[i])
	}
	if len(debugMsg) > 1 {
		return errors.New("Cannot pass more than one `debug` arg to drawLine")
	} else if len(debugMsg) == 1 {
		debugText = fmt.Sprintf("%s (%v, %v)", debugMsg, start, end)
	} else {
		debugText = ""
	}
	var x_min, x_max, y_min, y_max int
	if start.x == end.x {
		if start.y > end.y {
			y_min = end.y
			y_max = start.y
		} else {
			y_min = start.y
			y_max = end.y
		}
		for i := range y_max - y_min {
			draw(s, start.x, y_min+i, getDebugText(i), style)
		}
	} else if start.y == end.y {
		if start.x > end.x {
			x_min = end.x
			x_max = start.x
		} else {
			x_min = start.x
			x_max = end.x
		}
		for i := range x_max - x_min {
			draw(s, x_min+i, start.y, getDebugText(i), style)
		}
	} else {
		return errors.New("Drawline only draws horizontal or vertical lines")
	}
	s.Show()
	return nil
}

func drawBox(s tcell.Screen, topLeft, bottomRight cell, style BorderStyle, debug_pause ...time.Duration) error {
	var pause func(tcell.Screen, string)
	if len(debug_pause) == 0 {
		pause = func(_ tcell.Screen, _ string) {}
	} else if len(debug_pause) > 1 {
		return errors.New("May not pass more than one int debug pause to drawBox")
	} else {
		pause = func(s tcell.Screen, text string) {
			debug(s, text, true)
			time.Sleep(debug_pause[0])
		}
	}
	for pass := range style.widths.max() {
		pause(s, fmt.Sprint("doing Hline for pass ", pass))
		drawHLine(s, topLeft.y-pass, topLeft.x+pass, bottomRight.x+pass, style.drawStyle, fmt.Sprintf("pass %d", pass))
		pause(s, fmt.Sprint("doing Vline for pass ", pass))
		drawVLine(s, bottomRight.x-pass, topLeft.y-1-pass, bottomRight.y+pass, style.drawStyle, fmt.Sprintf("pass %d", pass))

	}

	return nil
}
