package layout

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"image/color"
)

func VSpace(sz float32) *canvas.Rectangle {
	space := canvas.NewRectangle(color.Transparent)
	space.SetMinSize(fyne.NewSize(0, sz))
	return space
}

func HSpace(sz float32) *canvas.Rectangle {
	space := canvas.NewRectangle(color.Transparent)
	space.SetMinSize(fyne.NewSize(sz, 0))
	return space
}
