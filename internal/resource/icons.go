package resource

import "fyne.io/fyne/v2"

var ProgramIcon fyne.Resource

func init() {
	ProgramIcon, _ = fyne.LoadResourceFromPath("assets/ytb.svg")
}
