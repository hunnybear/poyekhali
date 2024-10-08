package ui

import (
	"log"
	"time"

	"github.com/gdamore/tcell/v2"
)

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

	// Clear screen
	s.Show()
	s.Clear()
	return s
}

func UI() {
	s := newScreen()
	s.Clear()
}

func UIOnTest() {
	s := newScreen()
	time.Sleep(3 * time.Second)
	drawStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorPurple)
	drawLine(s, cell{10, 10}, cell{10, 32}, drawStyle)
	boxStyle1 := BorderStyle{
		widths:    BorderWidths{1, 2, 1, 2},
		drawStyle: tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorRed),
	}
	drawBox(s, cell{12, 12}, cell{35, 35}, boxStyle1, 1)
	time.Sleep(1 * time.Second)
	drawBox(s, cell{15, 15}, cell{41, 55}, BorderStyle{
		widths: BorderWidths{1, 1, 1, 1}, drawStyle: tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorBlue),
	}, 3, 2)
	time.Sleep(3 * time.Second)
}
