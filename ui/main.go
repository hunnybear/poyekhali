package ui

import (
	"errors"
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/hunnybear/poyekhali/tcellUtil"
	"github.com/hunnybear/poyekhali/types"
)

func UI() {
	screen, quit := tcellUtil.NewScreen()

	defer quit()
	tcellUtil.WriteText(screen, 5, 5, "abcdef", tcell.StyleDefault.Background(tcell.ColorPurple).Foreground(tcell.ColorOrange).Blink(true))
	screen.Show()
}

func StartMissionControl() types.QuitFunction {

	mc, quit, err := types.NewMissionControl(types.Auto)
	if err != nil {
		panic(err)
	}

	mc.StatusDebug("It begins!")

	renderer := NewCheckerRenderer(
		2,
		tcell.StyleDefault.Background(tcell.ColorRed).Foreground(tcell.ColorGreen))
	bspPane := mc.Pane.(*types.BSPPane)
	_, err = bspPane.NewViewChild(renderer, types.NewDefaultPDConstraints())
	if err != nil {
		quit()
		panic(err)
	}
	bspPane.Draw(*mc.Context, 0, 0)
	return quit
}

func UIOnTest() {

	s, quit := tcellUtil.NewScreen()
	defer quit()
	drawStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorPurple)
	tcellUtil.DrawLine(s, tcellUtil.Cell{10, 10}, tcellUtil.Cell{10, 32}, drawStyle)
	boxStyle2, err := tcellUtil.NewBorderStyle(
		tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorBlue))

	if err != nil {
		panic(err)
	}
	const (
		topTop    uint16 = 5
		topBottom        = 25
		leftLeft         = 5
		leftRight        = 25
		midTop           = 27
		midBottom        = 45
		midLeft          = 27
		midRight         = 205
	)
	// character
	tcellUtil.Draw(s, uint16(25), uint16(25), 'a', tcell.StyleDefault.Foreground(tcell.ColorYellow).Background(tcell.ColorBlueViolet))
	topLeftRect, err := tcellUtil.NewRect(topTop, leftRight, topBottom, leftLeft)
	if err != nil {
		panic(err)
	}
	midLeftBox, err := tcellUtil.NewRect(midTop, leftRight, midBottom, leftLeft)
	if err != nil {
		panic(err)
	}
	midTopRect, err := tcellUtil.NewRect(topTop, midRight, topBottom, midLeft)
	if err != nil {
		panic(err)
	}

	// Blue 1 cell border rect
	err = tcellUtil.DrawRect(s, topLeftRect, boxStyle2, "a rect eh?")
	if err != nil {
		panic(err)
	}

	// Green box [mid left]
	err = tcellUtil.DrawBox(
		s, midLeftBox,
		tcell.StyleDefault.Background(tcell.ColorGreen).Foreground(tcell.ColorRed),
	)
	if err != nil {
		panic(err)
	}
	bw, _ := tcellUtil.NewBorderWidths(1)
	boxBorder, err := tcellUtil.NewBorderStyle(
		tcellUtil.NewBGColorStyle(tcell.ColorLightGreen), bw)
	if err != nil {
		panic(err)
	}
	bw, _ = tcellUtil.NewBorderWidths(2, 1, 3, 1)
	bs, err := tcellUtil.NewBorderStyle(
		tcell.StyleDefault.Foreground(tcell.ColorRed).Background(tcell.ColorLightCyan),
		bw,
	)
	if err != nil {
		panic(err)
	}

	// coral box [mid top]
	err = tcellUtil.DrawBox(
		s, midTopRect,
		tcell.StyleDefault.Background(tcell.ColorCoral).Foreground(tcell.ColorDarkBlue),
		tcellUtil.NewBoxOptions().WithBorderStyle(bs),
	)

	if err != nil {
		panic(err)
	}
	i := 0
	debugStyle := tcell.StyleDefault
	ctrlC := 0
	shown := false
	for {

		s.Show()
		ev := s.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventResize:
			tcellUtil.WriteText(s, midLeft, midTop, "resize", tcell.StyleDefault)
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyCtrlC {
				ctrlC += 1
				if ctrlC >= 3 {

					tcellUtil.WriteText(s, midLeft, midTop, "exiting ctrlc", debugStyle)
					return
				}
			}
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {

				tcellUtil.WriteText(s, midLeft, midTop, fmt.Sprintf("got a key %v %b i: %d c: %d", ev.Key(), shown, i, ctrlC), debugStyle)
				if shown {

					tcellUtil.WriteText(s, midLeft, midTop, "ending", debugStyle)
					time.Sleep(time.Second)
					return
				}

				err = tcellUtil.DrawRect(s, midTopRect, boxStyle2, fmt.Sprintf("%d: BW: %v", i, boxBorder))
				if err != nil {
					panic(err)
				}
				shown = true
			}

		case *tcell.EventTime:
			tcellUtil.WriteText(s, midLeft, midTop+2, fmt.Sprintf("time event %d", i), tcell.StyleDefault)
		case *tcell.EventInterrupt:
			tcellUtil.WriteText(s, midLeft, midTop+2, fmt.Sprintf("interrupt event %d", i), tcell.StyleDefault)
		default:
			panic(errors.New(fmt.Sprintf("unrecognized %T event %v", ev, ev)))
		}
		tcellUtil.WriteText(s, midLeft, midTop+1, fmt.Sprintf("%b i: %d", shown, i), debugStyle)
		debugStyle = debugStyle.Reverse(true)
		i += 1
		if i >= 33 {
			panic(errors.New("PANICING PANICING"))
		}
	}
}
