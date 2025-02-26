package layout

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type VLayout struct {
	size []float32
}

func NewVLayout(size int, ratio ...float32) *VLayout {
	h := &VLayout{
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

func (v *VLayout) SetSize(i int, size float32) {
	if size < 0 {
		v.size[i] = -1
	} else if size > 1 {
		v.size[i] = 1
	} else {
		v.size[i] = size
	}
}

func (v *VLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	minSize := fyne.NewSize(0, 0)
	for _, child := range objects {
		if !child.Visible() {
			continue
		}
		minSize = minSize.Max(child.MinSize())
	}
	return minSize
}

func (v *VLayout) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	pos := fyne.NewPos(0, 0)
	for i, o := range objects {
		o.Move(pos)
		if v.size[i] < 0 {
			s := o.MinSize()
			o.Resize(s)
			pos = pos.AddXY(0, s.Height+theme.Padding())
		} else {
			h := containerSize.Height * v.size[i]
			o.Resize(fyne.NewSize(containerSize.Width, h))
			pos = pos.AddXY(0, h+theme.Padding())
		}
	}
}
