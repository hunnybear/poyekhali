package ui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/hunnybear/poyekhali/tcellUtil"
	"github.com/hunnybear/poyekhali/types"
	"github.com/hunnybear/poyekhali/util"
)

// view.go : views, which is to say, tiles

func NewCheckerRenderer(checkSize uint16, style tcell.Style) types.Renderer {

	renderer := func(screen tcell.Screen, row, col, height, width uint16) error {
		styleA := style
		styleB := style.Reverse(true)
		widths, err := tcellUtil.NewBorderWidths(-1)
		if err != nil {
			return err
		}
		cursorRow := row

		for checkRow := range height / checkSize {
			cursorCol := col
			for checkCol := range width / checkSize {

				var boxStyle tcellUtil.BorderStyle
				var err error
				if checkRow%2 == checkCol%2 {
					boxStyle, err = tcellUtil.NewBorderStyle(
						styleA, widths,
					)
					if err != nil {
						return err
					}
				} else {
					boxStyle, err = tcellUtil.NewBorderStyle(
						styleB, widths,
					)
					if err != nil {
						return err
					}
				}
				rect, err := tcellUtil.NewRect(
					cursorRow, cursorCol+checkSize,
					cursorRow+checkSize, col,
				)
				if err != nil {
					return err
				}
				if err := tcellUtil.DrawRect(
					screen, rect, boxStyle); err != nil {
					return err
				}
			}
			// do remainder of row (TODO)
		}

		for thisRow := range height {
			for thisCheckCol := range width / checkSize {
				thisCol := col + (thisCheckCol * checkSize)
				if thisRow%2 == 0 && thisCheckCol%2 == 0 {
					style = style.Reverse(true)
					srev := style.Reverse(false)
					fmt.Println("styles", style, srev)
				}
				if err = tcellUtil.DrawHLine(
					screen, row+thisRow, thisCol, thisCol+checkSize, style); err != nil {
					return err
				}
			}
		}

		return nil
	}

	return renderer
}

type StatusBar struct {
	tcell.Screen
	row          uint16
	width        uint16
	defaultStyle tcell.Style
}

func (bar *StatusBar) Write(message string, style ...tcell.Style) error {
	drawStyle, err := util.ExpectOneOptional(style, "style", "Write")
	//tood handle too long lines
	if err != nil {
		return err
	}
	if len(message)+5 < int(bar.width) {
		// add leader
		tcellUtil.DrawHLine(bar, bar.row, 0, bar.row-uint16(len(message)-5), *drawStyle)
	}
	tcellUtil.WriteText(bar, bar.width-uint16(5-len(message)), bar.row, message, *drawStyle)
	tcellUtil.DrawHLine(bar, bar.row, bar.width-5, bar.width, *drawStyle)
	return nil
}
