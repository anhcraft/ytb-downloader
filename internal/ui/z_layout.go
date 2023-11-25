package ui

import (
	"fyne.io/fyne/v2"
)

type ZLayout struct {
	size []float32
}

func NewZLayout(size int) *ZLayout {
	z := &ZLayout{
		size: make([]float32, size),
	}
	for i := 0; i < size; i++ {
		z.size[i] = 1
	}
	return z
}

func (z *ZLayout) SetSize(i int, size float32) {
	if size < 0 {
		z.size[i] = -1
	} else if size > 1 {
		z.size[i] = 1
	} else {
		z.size[i] = size
	}
}

func (z *ZLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	minSize := fyne.NewSize(0, 0)
	for _, child := range objects {
		if !child.Visible() {
			continue
		}
		minSize = minSize.Max(child.MinSize())
	}
	return minSize
}

func (z *ZLayout) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	pos := fyne.NewPos(0, 0)
	for i, o := range objects {
		if z.size[i] < 0 {
			o.Resize(o.MinSize())
		} else {
			w := containerSize.Width * z.size[i]
			h := containerSize.Height
			o.Resize(fyne.NewSize(w, h))
		}
		o.Move(pos)
	}
}
