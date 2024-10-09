package ui

import (
	"errors"

	"github.com/gdamore/tcell/v2"
)

type cell struct {
	x uint8
	y uint8
}

type quitFunction func()

var defaultStyle, debugStyle, warningStyle, errorStyle tcell.Style

func init() {
	defaultStyle = tcell.StyleDefault.Foreground(tcell.ColorBlack).Background(tcell.ColorWhite)
	debugStyle = tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorBlue)
	warningStyle = tcell.StyleDefault.Foreground(tcell.ColorBlack).Background(tcell.ColorOrange)
	errorStyle = tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorRed)

}

type BSPDivision int

const (
	Horizontal BSPDivision = iota
	Vertical
)

func (c cell) draw(s tcell.Screen, style tcell.Style) {
	draw(s, c.x, c.y, ' ', style)
}

func (c cell) drawChar(s tcell.Screen, char rune, style tcell.Style) {
	draw(s, c.x, c.y, char, style)
}

type PaneDimensionConstraints struct {
	MinWidth  int
	MaxWidth  int
	MinHeight int
	MaxHeight int
}

func NewPaneDimensionConstraints(maxW, maxH, minW, minH int) PaneDimensionConstraints {
	return PaneDimensionConstraints{minW, maxW, minH, maxH}
}

type PaneDimensions struct {
	Width  uint8
	Height uint8
	PaneDimensionConstraints
}

func NewDefaultPDConstraints() PaneDimensionConstraints {
	return NewPaneDimensionConstraints(-1, -1, 1, 1)
}

func NewPaneDimensions(constraints PaneDimensionConstraints) *PaneDimensions {
	return &PaneDimensions{
		PaneDimensionConstraints: constraints,
		Height:                   0,
		Width:                    0,
	}
}

func (pd PaneDimensions) GetDims() (uint8, uint8) {
	return pd.Width, pd.Height
}
func (pd PaneDimensions) GetMinDims() (int, int) {
	return pd.MinWidth, pd.MinHeight
}

func (pd PaneDimensions) GetMaxDims() (int, int) {
	return pd.MaxWidth, pd.MaxHeight
}

func (pd *PaneDimensions) SetWidth(width uint8) error {
	if width < 1 {
		return errors.New("cannot Set Width to negative number")
	}
	pd.Width = width
	return nil
}

func (pd *PaneDimensions) SetHeight(height uint8) error {
	if height < 1 {
		return errors.New("cannot set Height to negative number")
	}
	pd.Height = height
	return nil
}

func (p *PaneDimensions) SetSize(width, height uint8) error {
	maxWidth, maxHeight := p.GetMaxDims()
	minWidth, minHeight := p.GetMinDims()
	checkWidth := int(width)
	checkHeight := int(height)
	if maxWidth > 0 && checkWidth > maxWidth {
		return errors.New("cannot set size of pane larger than maxwidth")
	} else if minWidth > 0 && checkWidth < minWidth {
		return errors.New("cannot set size of pane smallwer than minwidth")
	} else if maxHeight > 0 && checkHeight > maxHeight {
		return errors.New("cannot set size of pane larger than MaxHeight")
	} else if minHeight > 0 && checkHeight < minHeight {
		return errors.New("cannot set size of pane smaller than minheight")
	}
	return nil
}

type Pane interface {
	Draw(uint8, uint8) error // row, col
	SetWidth(uint8) error
	SetHeight(uint8) error
	GetDims() (uint8, uint8) // width, height
	GetMinDims() (int, int)  // width, height
	GetMaxDims() (int, int)  // width, height
}
type ViewPane struct {
	View
	PaneDimensions
}

func (vp *ViewPane) Draw(row, col uint8) error {
	return vp.View.Draw(row, col, vp.Width, vp.Height)
}

type Children map[Pane]int

// BSPPane is also a Pane
type BSPPane struct {
	PaneDimensions
	Children
	Division BSPDivision
}

func (bsp *BSPPane) Draw(row, col uint8) error {
	for child, offsetInt := range bsp.Children {
		offset := uint8(offsetInt)
		if bsp.Division == Horizontal {
			col = col + offset
		} else if bsp.Division == Vertical {
			row = row + offset
		}
		if err := child.Draw(row, col); err != nil {
			return err
		}
	}
	return nil
}

func (bsp *BSPPane) ReCalc() error {
	if bsp.Division == Horizontal {
		width := bsp.Width / uint8(len(bsp.Children))
		remainder := bsp.Width / uint8(len(bsp.Children))
		remainderWidth := width + uint8(1)
		for child, _ := range bsp.Children {
			if remainder > 0 {
				child.SetWidth(remainderWidth)
				remainder -= 1
			} else {
				child.SetWidth(width)
			}
		}
	} else {
		// Vertical divsiion
		height := bsp.Height / uint8(len(bsp.Children))
		remainder := bsp.Height % uint8(len(bsp.Children))
		remainderHeight := height + uint8(1)
		for child, _ := range bsp.Children {
			if remainder > 0 {
				child.SetHeight(remainderHeight)
				remainder -= 1
			} else {
				child.SetHeight(height)
			}
		}
	}
	return nil
}

func newBSP(division BSPDivision, constraints ...PaneDimensionConstraints) (*BSPPane, error) {
	var useConstraints PaneDimensionConstraints
	if err := util.expectOneOptional(constraints); err != nil {
		return nil, err
	}
	return &BSPPane{
		*NewPaneDimensions(useConstraints),
		Children{},
		division,
	}, nil
}

func (bsp *BSPPane) AddChild(child Pane) (int, error) {
	bsp.Children[child] = 0
	if err := bsp.ReCalc(); err != nil {
		return -1, err
	}
	return len(bsp.Children), nil
}

func (bsp *BSPPane) NewBSPChild(division BSPDivision, constraints ...PaneDimensionConstraints) (Pane, error) {

	panePtr, err := newBSP(division, constraints...)
	if err != nil {
		return nil, err
	}
	bsp.AddChild(panePtr)
	return panePtr, nil

}

func (bsp *BSPPane) newViewChild(view *View, constraints ...PaneDimensionConstraints) (Pane, error) {

	panePtr := &ViewPane{}

	return panePtr, nil
}

type Window struct {
	screen  tcell.Screen
	pane    *Pane
	context *MissionControlContext
	StatusBar
	BorderStyle
}

func NewWindow(ctx poyekhaliContext, cols, rows uint8, pane *Pane) *Window {
	return &Window{screen: screen, pane, cols, Height: rows}
}

type MissionControl struct {
	Window
}

func NewMissionControl(context ...*MissionControlContext) (*MissionControl, quitFunction) {

	screen, quit := newScreen()
	rows, cols := screen.Size()
	return &MissionControl{Window: *NewWindow(screen, uint8(cols), uint8(rows))}, quit
}

func (mc *MissionControl) StatusDebug(message string) {
	mc.StatusBar.Write(message, debugStyle)
}

type BorderWidths struct {
	top    uint8
	right  uint8
	bottom uint8
	left   uint8
}
