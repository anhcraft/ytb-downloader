package menu

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func requestInput() fyne.CanvasObject {
	input = widget.NewMultiLineEntry()
	input.SetPlaceHolder("Enter URL(s) of videos, playlists, etc")
	input.SetMinRowsVisible(5)
	return input
}
