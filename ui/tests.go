package ui

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/gdamore/tcell/v2"
)

type ConfRGBColor [3]int32

func (c ConfRGBColor) ToRGBColor() tcell.Color {
	return tcell.NewRGBColor(c[0], c[1], c[2])
}

type TestLine struct {
	Start cell
	End   cell
	Text  string `json:"text,omitempty"`
}

type TestBox struct {
	Start   cell
	End     cell
	Widths  []int  `json:"widths,omitempty"`
	Caption string `json:"caption,omitempty"`
}

type TestDrawing struct {
	TextColor ConfRGBColor `json:"textColor"`
	DrawColor ConfRGBColor `json:"drawColor"`
	TestWords []string     `json:"testWords"`
	TestLines []TestLine   `json:"testLines"`
	TestBoxes []TestBox    `json:"testBoxes"`
	Pause     int          `json:"pause"`
	PauseDiv  int          `json:"pauseDiv`
}

func (td TestDrawing) getPause() time.Duration {
	pause_div := td.PauseDiv
	if pause_div < 1 {
		pause_div = 1
	}
	return time.Second * time.Duration(td.Pause/pause_div)
}

type TestConfig struct {
	Drawings  []TestDrawing `json:"drawings"`
	Pause     uint16        `json:"pause"`
	Pause_div uint16        `json:"pauseDiv"`
}

func (tc TestConfig) getPause() time.Duration {
	pause_div := tc.Pause_div
	if pause_div < 1 {
		pause_div = 1
	}
	return time.Second * time.Duration(tc.Pause/pause_div)
}

func TestUIFromFile(testFilePtr *string) {
	content, err := ioutil.ReadFile(*testFilePtr)
	if err != nil {
		panic(err)
	}
	TestUI(content)
}

func TestUI(content []byte) {

	var config TestConfig

	err := json.Unmarshal(content, &config)
	if err != nil {
		panic(err)
	}
	if err = json.Unmarshal(content, &config); err != nil {
		panic(err)
	}
	s := newScreen()
	//w, h := s.Size()
	//fmt.Println("Screen size w/h", w, h)
	//fmt.Println("sleeping", config.getPause())
	//time.Sleep(config.getPause())
	for i, drawing := range config.Drawings {

		debug(s, fmt.Sprintf("doing drawing %d", i), true)
		drawingStyle := tcell.StyleDefault.Foreground(drawing.TextColor.ToRGBColor()).Background(drawing.DrawColor.ToRGBColor())

		for j, line := range drawing.TestLines {
			debug(s, fmt.Sprintf("doing line %d for drawing %d", j, i), true)
			drawLine(s, line.Start, line.End, drawingStyle, line.Text)
			time.Sleep(drawing.getPause())
		}

		for j, box := range drawing.TestBoxes {
			debug(s, fmt.Sprintf("doing box %d for drawing %d", j, i), true)
			boxWidths, err := NewBorderWidths(box.Widths...)
			if err != nil {
				panic(err)
			}
			boxStyle := BorderStyle{
				widths:    boxWidths,
				drawStyle: drawingStyle,
			}
			drawBox(s, box.Start, box.End, boxStyle, drawing.getPause())
			time.Sleep(drawing.getPause())
		}
	}
	s.Show()
	time.Sleep(config.getPause())
}
