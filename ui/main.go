package ui

import (
	"errors"
	"fmt"
	"log"
	"time"

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

func draw(s tcell.Screen, x, y int, char rune, style tcell.Style) {
	s.SetContent(x, y, char, nil, style)
}

func newScreen() tcell.Screen {
	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := s.Init(); err != nil {
		log.Fatalf("%+v", err)
	}

	// Set default text style
	defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	s.SetStyle(defStyle)

	// Clear screen
	s.Show()
	s.Clear()
	return s
}

func drawLine(s tcell.Screen, start, end cell, style tcell.Style) error {
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
			draw(s, start.x, y_min+i, ' ', style)
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
			draw(s, x_min+i, start.y, ' ', style)
		}
	} else {
		return errors.New("Drawline only draws horizontal or vertical lines")
	}
	s.Show()
	return nil
}

func drawBox(s tcell.Screen, start, end cell, style tcell.Style)

func Ui() {
	fmt.Println("WHellllllp")
	s := newScreen()
	quit := func() {
		// You have to catch panics in a defer, clean up, and
		// re-raise them - otherwise your application can
		// die without leaving any diagnostic trace.
		maybePanic := recover()
		s.Fini()
		if maybePanic != nil {
			panic(maybePanic)
		}
	}
	defer quit()
	time.Sleep(3 * time.Second)
	drawStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorPurple)
	drawLine(s, cell{10, 10}, cell{10, 32}, drawStyle)
	time.Sleep(3 * time.Second)
}

func main() {}
