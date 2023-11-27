package resource

import "fyne.io/fyne/v2"

var ProgramIcon fyne.Resource
var EraserIcon fyne.Resource

func init() {
	ProgramIcon, _ = fyne.LoadResourceFromPath("assets/ytb.svg")
	EraserIcon, _ = fyne.LoadResourceFromPath("assets/eraser.svg")
}
