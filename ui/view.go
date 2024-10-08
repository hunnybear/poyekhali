package ui

import "github.com/gdamore/tcell/v2"

// view.go : views, which is to say, tiles

type View struct {
	screen    tcell.Screen
	minWidth  int
	maxWidth  int
	minHeight int
	maxHeight int
}
