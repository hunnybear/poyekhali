package tcellUtil

import (
	"errors"
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/hunnybear/poyekhali/util"
)

type Side int

const (
	Top    Side = iota
	Bottom Side = iota
	Right  Side = iota
	Left   Side = iota
)

func NewBGColorStyle(color tcell.Color) tcell.Style {
	return tcell.StyleDefault.Background(color)
}

type BorderWidths struct {
	Top    uint16
	Right  uint16
	Bottom uint16
	Left   uint16
}

func (bw BorderWidths) asSlice() [4]uint16 {
	return [4]uint16{bw.Top, bw.Right, bw.Bottom, bw.Left}
}

func (bw BorderWidths) Max() uint16 {
	return max(bw.Top, bw.Right, bw.Bottom, bw.Left)
}

func (bw BorderWidths) Min() uint16 {
	return min(bw.Top, bw.Right, bw.Bottom, bw.Left)
}

func NewBorderWidths(widths ...int) (BorderWidths, error) {
	if len(widths) == 1 {
		return BorderWidths{
			uint16(widths[0]),
			uint16(widths[0]),
			uint16(widths[0]),
			uint16(widths[0]),
		}, nil
	} else if len(widths) == 4 {
		return BorderWidths{
			uint16(widths[0]),
			uint16(widths[1]),
			uint16(widths[2]),
			uint16(widths[3])}, nil
	} else {
		return BorderWidths{1, 1, 1, 1}, errors.New("NewBorderWidts requires either 1 or 4 ints")
	}
}

type BorderStyle struct {
	Widths    BorderWidths
	DrawStyle tcell.Style
}

func NewBorderStyle(drawStyle tcell.Style, borderWidths ...BorderWidths) (BorderStyle, error) {
	widthsPtr, err := util.ExpectOneOptional(
		borderWidths, "borderWidths", "NewBorderStyle")

	if err != nil {
		return BorderStyle{}, err
	}
	var widths BorderWidths
	if widthsPtr == nil {
		widths, err = NewBorderWidths(1)
		if err != nil {
			return BorderStyle{}, err
		}
	} else {
		widths = *widthsPtr
	}
	return BorderStyle{Widths: widths, DrawStyle: drawStyle}, nil
}

func Draw(s tcell.Screen, x, y uint16, char rune, style tcell.Style) {
	s.SetContent(int(x), int(y), char, nil, style)
}

func WriteText(s tcell.Screen, x, y uint16, text string, style tcell.Style) {
	for i, char := range text {
		ui8 := uint16(i)
		Draw(s, x+ui8, y, char, style)
	}
}

func DrawHLine(s tcell.Screen, y, start, end uint16, style tcell.Style, debugMsg ...string) error {
	if len(debugMsg) > 0 {
		return DrawLine(s, Cell{X: start, Y: y}, Cell{X: end, Y: y}, style, debugMsg[0])
	} else {
		return DrawLine(s, Cell{X: start, Y: y}, Cell{X: end, Y: y}, style)
	}
}

func DrawVLine(s tcell.Screen, x, start, end uint16, style tcell.Style, debugMsg ...string) error {
	if len(debugMsg) > 0 {
		return DrawLine(s, Cell{X: x, Y: start}, Cell{X: x, Y: end}, style, debugMsg[0])
	} else {
		return DrawLine(s, Cell{X: x, Y: start}, Cell{X: x, Y: end}, style)
	}
}

func DrawLine(s tcell.Screen, start, end Cell, style tcell.Style, debugMsg ...string) error {
	var debugText string
	if len(debugMsg) > 1 {
		return errors.New("cannot pass more than one `debug` arg to drawLine")
	} else if len(debugMsg) == 1 {
		debugText = debugMsg[0]
	} else {
		debugText = ""
	}
	getDebugText := func(i uint16) rune {
		ii := int(i)
		if ii >= len(debugText) {
			return ' '
		}
		return rune(debugText[i])
	}
	var x_min, x_max, y_min, y_max uint16
	if start.X == end.X {
		if start.Y > end.Y {
			y_min = end.Y
			y_max = start.Y
		} else {
			y_min = start.Y
			y_max = end.Y
		}
		for i := range y_max - y_min + 1 {
			Draw(s, start.X, y_min+i, getDebugText(i), style)
		}
	} else if start.Y == end.Y {
		if start.X > end.X {
			x_min = end.X
			x_max = start.X
		} else {
			x_min = start.X
			x_max = end.X
		}
		for i := range x_max - x_min + 1 {
			Draw(s, x_min+i, start.Y, getDebugText(i), style)
		}
	} else {
		return errors.New("Drawline only draws horizontal or vertical lines")
	}
	//s.Show()
	return nil
}

// TODO: title position, truncate/wrap text
type BoxOptions struct {
	Title string
	TitleSides
	*BorderStyle
}

func NewBoxOptions() BoxOptions {
	return BoxOptions{}
}

func (opts BoxOptions) WithTitle(title string) BoxOptions {
	opts.Title = title
	return opts
}

func (opts BoxOptions) WithStyle(style tcell.Style) BoxOptions {
	newBorderStylePtr := &BorderStyle{}
	*newBorderStylePtr = *opts.BorderStyle
	opts.BorderStyle = newBorderStylePtr
	opts.BorderStyle.DrawStyle = style
	return opts
}

func (opts BoxOptions) WithBorderStyle(style BorderStyle) BoxOptions {
	opts.BorderStyle = &style
	return opts
}

func DrawBox(s tcell.Screen, rect Rectangle, style tcell.Style, options ...BoxOptions) error {
	boxOptions, err := util.ExpectOneOptional(options, "options", "DrawBox")

	if err != nil {
		return err
	}
	// offsets
	lOffset := uint16(0)
	rOffset := uint16(0)
	tOffset := uint16(0)
	bOffset := uint16(0)

	if boxOptions != nil {
		lOffset = boxOptions.Widths.Left
		rOffset = boxOptions.Widths.Right
		bOffset = boxOptions.Widths.Bottom
		tOffset = boxOptions.Widths.Top
	}

	if err := rect.Validate(); err != nil {
		return err
	}
	// Fill

	for row := range rect.BottomRight.Y + 1 - rect.TopLeft.Y - bOffset - tOffset {
		err := DrawHLine(s,
			rect.TopLeft.Y+tOffset+row,
			rect.TopLeft.X+lOffset,
			rect.BottomRight.X-rOffset, style, fmt.Sprintf("row %v, %d, %d, %d ,%d", row, tOffset, rOffset, bOffset, lOffset))
		if err != nil {
			return err
		}
	}

	// Border and title
	if boxOptions != nil {

		err := DrawRect(s, rect, *boxOptions.BorderStyle, boxOptions.Title)
		if err != nil {
			return err
		}
		return nil
	}
	return nil
}

func truncateString(s string, l int) string {
	if l < 0 {
		return s
	}

	truncated := []rune{}

	for i, c := range s {
		if i+2 >= l {
			return string(truncated) + ".."
		}
		truncated = append(truncated, c)
	}
	return string(truncated)
}

func DrawRect(s tcell.Screen, rect Rectangle, style BorderStyle, title ...string) error {
	titlePtr, err := util.ExpectOneOptional(title, "title", "DrawRect")
	if err != nil {
		return err
	}
	for passInt := range style.Widths.Max() {
		pass := uint16(passInt)
		// top
		if style.Widths.Top > pass {
			if err := DrawHLine(
				s, rect.TopLeft.Y, rect.TopLeft.X, rect.BottomRight.X, style.DrawStyle, *titlePtr); err != nil {
				return err
			}
			rect.TopLeft.Y += 1
		}
		// bottom
		if style.Widths.Bottom > pass {
			if err := DrawHLine(
				s, rect.BottomRight.Y,
				rect.TopLeft.X, rect.BottomRight.X,
				style.DrawStyle, *titlePtr); err != nil {
				return err
			}
			rect.BottomRight.Y -= 1
		}
		// right
		if style.Widths.Right > pass {
			if err := DrawVLine(
				s, rect.BottomRight.X,
				rect.TopLeft.Y, rect.BottomRight.Y,
				style.DrawStyle, *titlePtr); err != nil {
				return err
			}
			rect.BottomRight.X -= 1
		}

		// left
		if style.Widths.Left > pass {
			if err := DrawVLine(
				s, rect.TopLeft.X,
				rect.TopLeft.Y, rect.BottomRight.Y,
				style.DrawStyle, *titlePtr); err != nil {
				return err
			}
			rect.TopLeft.X += 1
		}
	}

	*titlePtr = fmt.Sprintf("t: %d r %d b %d l %d", style.Widths.Top, style.Widths.Right, style.Widths.Bottom, style.Widths.Top)
	if titlePtr != nil && *titlePtr != "" {
		max_len := rect.BottomRight.X - rect.TopLeft.X - style.Widths.Right - style.Widths.Left - 2
		if len(*titlePtr) > int(max_len) {

			// todo: truncate
			// probably struct of title options (position, truncate/wrap/error, text)
			//return errors.New("message is too long for box")
			*titlePtr = truncateString(*titlePtr, int(max_len)-1)
		}
		xStart := rect.TopLeft.X + style.Widths.Left
		DrawHLine(
			s, rect.TopLeft.Y+style.Widths.Top, xStart, xStart+uint16(len(*titlePtr)),
			style.DrawStyle, *titlePtr)
	}

	return nil
}
