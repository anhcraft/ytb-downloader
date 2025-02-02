package layout

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type HLayout struct {
	size []float32
}

func NewHLayout(size int, ratio ...float32) *HLayout {
	h := &HLayout{
		size: make([]float32, size),
	}
	for i := 0; i < size; i++ {
		h.size[i] = 1
	}
	for i := 0; i < len(ratio); i++ {
		h.size[i] = ratio[i]
	}
	return h
}

func (h *HLayout) SetSize(i int, size float32) {
	if size < 0 {
		h.size[i] = -1
	} else if size > 1 {
		h.size[i] = 1
	} else {
		h.size[i] = size
	}
}

func (h *HLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	minSize := fyne.NewSize(0, 0)
	for _, child := range objects {
		if !child.Visible() {
			continue
		}
		minSize = minSize.Max(child.MinSize())
	}
	return minSize
}

func (h *HLayout) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	pos := fyne.NewPos(0, 0)
	for i, o := range objects {
		o.Move(pos)
		if h.size[i] < 0 {
			s := o.MinSize()
			o.Resize(s)
			pos = pos.AddXY(s.Width+theme.Padding(), 0)
		} else {
			w := containerSize.Width * h.size[i]
			o.Resize(fyne.NewSize(w, containerSize.Height))
			pos = pos.AddXY(w+theme.Padding(), 0)
		}
	}
}
