package component

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type CopyableLabel struct {
	widget.Label

	win fyne.Window
}

func NewCopyableLabel(text string, win fyne.Window) *CopyableLabel {
	label := &CopyableLabel{}
	label.win = win
	label.Text = text
	label.Truncation = fyne.TextTruncateEllipsis
	label.ExtendBaseWidget(label)
	return label
}

func NewWrappedCopyableLabel(text string, win fyne.Window, height float32) *container.Scroll {
	label := NewCopyableLabel(text, win)
	label.Truncation = fyne.TextTruncateOff
	label.Wrapping = fyne.TextWrapBreak
	scroll := container.NewVScroll(label)
	scroll.SetMinSize(fyne.NewSize(0, height))
	return scroll
}

func (l *CopyableLabel) TappedSecondary(e *fyne.PointEvent) {
	menu := fyne.NewMenu("",
		fyne.NewMenuItem("Copy", func() {
			l.win.Clipboard().SetContent(l.Text)
		}),
	)

	widget.ShowPopUpMenuAtPosition(menu, fyne.CurrentApp().Driver().CanvasForObject(l), e.AbsolutePosition)
}
