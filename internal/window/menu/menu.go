package menu

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"ytb-downloader/internal/constants"
	"ytb-downloader/internal/handle/request"
	"ytb-downloader/internal/resource"
	layout2 "ytb-downloader/internal/ui/layout"
)

// NOTE on concurrency model:
// The whole UI is rendered by the main thread
// While the download scheduler and worker are independent goroutines
// We accept data race here... as UI updates are purely for visual

var win fyne.Window
var table *widget.Table
var input *widget.Entry

func OpenMenu(app fyne.App) fyne.Window {
	CheckUpdate(func(latest bool, currVer string, latestVer string, err error) {
		if !latest {
			dialog.ShowInformation(
				"Update",
				fmt.Sprintf("Latest version: %s\nCurrent version: %s\nPlease update the app!", latestVer, currVer),
				win,
			)
		}
	})

	request.GetQueue().SetUpdateCallback(func(req *request.Request) {
		// TODO thread-safe?
		if req.Status() == request.StatusFailed {
			input.SetText(input.Text + "\n" + req.RawUrl())
		}
		table.Refresh()
	})

	win = app.NewWindow("Yt-dlp GUI")
	win.SetContent(container.New(
		layout2.NewVLayout(2, 0.3, 0.7),
		container.NewVBox(
			toolbar(app),
			container.New(layout2.NewHLayout(2, 0.5, 0.5), requestInput(), requestSettings()),
			layout2.VSpace(10),
		),
		requestTable(),
	))
	win.Resize(fyne.NewSize(constants.MainWindowWidth, constants.MainWindowHeight))
	win.SetFixedSize(true)
	win.SetPadded(true)
	win.SetIcon(resource.ProgramIcon)
	win.SetMaster()
	win.CenterOnScreen()
	win.ShowAndRun()

	return win
}
