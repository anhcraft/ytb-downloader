package resource

import "fyne.io/fyne/v2"

var RegularFont fyne.Resource
var BoldFont fyne.Resource

func init() {
	RegularFont, _ = fyne.LoadResourceFromPath("./assets/NotoSansSC-Regular.ttf")
	BoldFont, _ = fyne.LoadResourceFromPath("./assets/NotoSansSC-Bold.ttf")
}
