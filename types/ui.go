package types

import (
	"context"
	"errors"

	"github.com/gdamore/tcell/v2"
	"github.com/hunnybear/poyekhali/tcellUtil"
	"github.com/hunnybear/poyekhali/util"
)

type key int

const (
	ScreenKey key = iota
)

type Cell tcellUtil.Cell

type QuitFunction func()

type PoyekhaliContext struct {
	context.Context
}

func (ctx *PoyekhaliContext) SetScreen(screen tcell.Screen) error {
	if val := ctx.Value(ScreenKey); val != nil {
		//already there
		return errors.New("context already has a screen assigned")
	}
	ctx.Context = context.WithValue(ctx.Context, ScreenKey, screen)
	return nil
}

// so lazy.
func (ctx *PoyekhaliContext) GetScreen() tcell.Screen {
	return ctx.Value(ScreenKey).(tcell.Screen)
}

type MissionControlContext struct {
	PoyekhaliContext
}

func NewMissionControlContext() (*MissionControlContext, QuitFunction) {
	ctx := MissionControlContext{
		PoyekhaliContext{context.Background()},
	}
	screen, quit := tcellUtil.NewScreen()
	ctx.SetScreen(screen)
	return &ctx, quit
}

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
	Auto
)

func (c Cell) draw(ctx PoyekhaliContext, style tcell.Style) {
	tcellUtil.Draw(ctx.Value(ScreenKey).(tcell.Screen), c.X, c.Y, ' ', style)
}

func (c Cell) drawChar(s tcell.Screen, char rune, style tcell.Style) {
	tcellUtil.Draw(s, c.X, c.Y, char, style)
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
	Width  uint16
	Height uint16
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

func (pd PaneDimensions) GetDims() (uint16, uint16) {
	return pd.Width, pd.Height
}
func (pd PaneDimensions) GetMinDims() (int, int) {
	return pd.MinWidth, pd.MinHeight
}

func (pd PaneDimensions) GetMaxDims() (int, int) {
	return pd.MaxWidth, pd.MaxHeight
}

func (pd *PaneDimensions) SetWidth(width uint16) error {
	if width < 1 {
		return errors.New("cannot Set Width to negative number")
	}
	pd.Width = width
	return nil
}

func (pd *PaneDimensions) SetHeight(height uint16) error {
	if height < 1 {
		return errors.New("cannot set Height to negative number")
	}
	pd.Height = height
	return nil
}

func (p PaneDimensions) SetSize(width, height uint16) error {
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
	Draw(PoyekhaliContext, uint16, uint16) error // row, col
	SetWidth(uint16) error
	SetHeight(uint16) error
	GetDims() (uint16, uint16) // width, height
	GetMinDims() (int, int)    // width, height
	GetMaxDims() (int, int)    // width, height
}
type ViewPane struct {
	Renderer
	PaneDimensions
}

func (vp ViewPane) Draw(ctx PoyekhaliContext, row, col uint16) error {
	return vp.Renderer(ctx.GetScreen(), row, col, vp.Height, vp.Width)
}

type Children map[Pane]int

// BSPPane is also a Pane
type BSPPane struct {
	PaneDimensions
	Children
	Division BSPDivision
}

func (bsp BSPPane) Draw(ctx PoyekhaliContext, row, col uint16) error {
	for child, offsetInt := range bsp.Children {
		offset := uint16(offsetInt)
		if bsp.Division == Horizontal {
			col = col + offset
		} else if bsp.Division == Vertical {
			row = row + offset
		}
		if err := child.Draw(ctx, row, col); err != nil {
			return err
		}
	}
	return nil
}

func (bsp BSPPane) ReCalc() error {
	if bsp.Division == Auto {
		// Set division based off of dimensions]
		if bsp.Height > bsp.Width {
			bsp.Division = Vertical
		} else {
			bsp.Division = Horizontal
		}
	}
	if bsp.Division == Horizontal {
		width := bsp.Width / uint16(len(bsp.Children))
		remainder := bsp.Width / uint16(len(bsp.Children))
		remainderWidth := width + uint16(1)
		for child := range bsp.Children {
			if remainder > 0 {
				child.SetWidth(remainderWidth)
				remainder -= 1
			} else {
				child.SetWidth(width)
			}
		}
	} else {
		// Vertical divsiion
		height := bsp.Height / uint16(len(bsp.Children))
		remainder := bsp.Height % uint16(len(bsp.Children))
		remainderHeight := height + uint16(1)
		for child := range bsp.Children {
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

func NewBSP(division BSPDivision, constraintsSlice ...PaneDimensionConstraints) (*BSPPane, error) {
	useConstraints, err := util.ExpectOneOptional(constraintsSlice, "contraints", "NewBSP")
	if err != nil {
		return nil, err
	} else if useConstraints == nil {
		useConstraintsVal := NewPaneDimensionConstraints(-1, -1, 0, 0)
		useConstraints = &useConstraintsVal
	}
	pdPtr := NewPaneDimensions(*useConstraints)
	return &BSPPane{
		*pdPtr,
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

	panePtr, err := NewBSP(division, constraints...)
	if err != nil {
		return nil, err
	}
	bsp.AddChild(panePtr)
	return panePtr, nil

}

func (bsp *BSPPane) NewViewChild(renderer Renderer, dimConstraints ...PaneDimensionConstraints) (Pane, error) {

	constraintsPtr, err := util.ExpectOneOptional(dimConstraints, "dimConstraints", "newViewChild")
	if err != nil {
		return nil, err
	}
	if constraintsPtr == nil {
		constraints := NewDefaultPDConstraints()
		constraintsPtr = &constraints
	}
	panePtr := &ViewPane{
		renderer, *NewPaneDimensions(*constraintsPtr),
	}

	return panePtr, nil
}

type Window struct {
	Pane
	Context *PoyekhaliContext
	StatusBar
	tcellUtil.BorderStyle
}

func NewWindow(ctx *PoyekhaliContext, cols, rows uint16, pane Pane, borderStyle ...tcellUtil.BorderStyle) (*Window, error) {
	if ctx == nil || ctx.Context == nil {
		panic(errors.New("context cannot be nil"))
	}
	style, err := util.ExpectOneOptional(borderStyle, "style", "NewWindow")
	if err != nil {
		return nil, err
	}
	if style == nil {
		styleVal, err := tcellUtil.NewBorderStyle(defaultStyle)
		if err != nil {
			return nil, err
		}
		style = &styleVal
	}

	return &Window{pane, ctx, *NewStatusBar(*ctx, rows-2, cols, style.DrawStyle), *style}, nil
}

type MissionControl struct {
	Window
}

func NewMissionControl(division BSPDivision) (*MissionControl, QuitFunction, error) {

	ctx, quit := NewMissionControlContext()
	rows, cols := ctx.GetScreen().Size()
	pane, err := NewBSP(division)
	if err != nil {
		return nil, func() {}, err
	}
	window, err := NewWindow(
		&ctx.PoyekhaliContext,
		uint16(cols),
		uint16(rows),
		pane,
	)
	if err != nil {
		return nil, func() {}, err
	}
	mc := &MissionControl{
		Window: *window,
	}

	return mc, quit, nil
}

func (mc *MissionControl) StatusDebug(message string) {
	mc.StatusBar.Write(message, debugStyle)
}

// view.go : views, which is to say, tiles

// Panes are the layout element. That is, they have position and size

type Renderer func(screen tcell.Screen, row, col, height, width uint16) error

type StatusBar struct {
	context      *PoyekhaliContext
	row          uint16
	width        uint16
	defaultStyle tcell.Style
}

func NewStatusBar(ctx PoyekhaliContext, row, width uint16, style tcell.Style) *StatusBar {
	return &StatusBar{
		&ctx, row, width, style,
	}
}

func (bar *StatusBar) Write(message string, style ...tcell.Style) error {
	if len(style) > 1 {
		return errors.New("may only pass one style to write to status bar")
	}
	//fmt.Println("writing", message)
	//tood handle too long lines
	if len(message)+5 < int(bar.width) {
		// add leader
		tcellUtil.DrawHLine(bar.context.GetScreen(),
			bar.row-5, 0, bar.row-uint16(len(message)-5), debugStyle)
	}
	tcellUtil.WriteText(
		bar.context.GetScreen(), bar.width-uint16(5-len(message)), bar.row-6, message, debugStyle)
	tcellUtil.DrawHLine(
		bar.context.GetScreen(), bar.row-5, bar.width-5, bar.width, debugStyle)
	return nil
}
