package ui

import (
	"encoding/json"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/hunnybear/poyekhali/tcellUtil"
)

type ConfRGBColor [3]int32

func (c ConfRGBColor) ToRGBColor() tcell.Color {
	return tcell.NewRGBColor(c[0], c[1], c[2])
}

type TestLine struct {
	Start tcellUtil.Cell
	End   tcellUtil.Cell
	Text  string `json:"text,omitempty"`
}

type TestBox struct {
	Rect    tcellUtil.Rectangle
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
	PauseDiv  int          `json:"pauseDiv"`
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
	content, err := os.ReadFile(*testFilePtr)
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
	s, quit := tcellUtil.NewScreen()
	defer quit()
	for _, drawing := range config.Drawings {

		drawingStyle := tcell.StyleDefault.Foreground(drawing.TextColor.ToRGBColor()).Background(drawing.DrawColor.ToRGBColor())

		for _, line := range drawing.TestLines {
			tcellUtil.DrawLine(s, line.Start, line.End, drawingStyle, line.Text)
			time.Sleep(drawing.getPause())
		}

		for _, box := range drawing.TestBoxes {
			boxWidths, err := tcellUtil.NewBorderWidths(box.Widths...)
			if err != nil {
				panic(err)
			}
			boxStyle, err := tcellUtil.NewBorderStyle(drawingStyle, boxWidths)
			if err != nil {
				panic(err)
			}

			tcellUtil.DrawRect(s, box.Rect, boxStyle)
			time.Sleep(drawing.getPause())
		}
	}
	s.Show()
	time.Sleep(config.getPause())
}
