package component

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type CopyableLabel struct {
	widget.Label

	win fyne.Window
}

func (c *CopyableLabel) MinSize() fyne.Size {
	c.ExtendBaseWidget(c)
	return c.BaseWidget.MinSize()
}

func NewCopyableLabel(text string, win fyne.Window) *CopyableLabel {
	label := &CopyableLabel{}
	label.win = win
	label.Text = text
	label.ExtendBaseWidget(label)
	return label
}

func NewWrappedCopyableLabel(text string, win fyne.Window) *CopyableLabel {
	label := NewCopyableLabel(text, win)
	label.Truncation = fyne.TextTruncateOff
	label.Wrapping = fyne.TextWrapWord
	return label
}

func (c *CopyableLabel) TappedSecondary(e *fyne.PointEvent) {
	menu := fyne.NewMenu("",
		fyne.NewMenuItem("Copy", func() {
			c.win.Clipboard().SetContent(c.Text)
		}),
	)

	widget.ShowPopUpMenuAtPosition(menu, fyne.CurrentApp().Driver().CanvasForObject(c), e.AbsolutePosition)
}
