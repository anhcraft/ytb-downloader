package debug

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"ytb-downloader/internal/constants"
	"ytb-downloader/internal/handle/request"
	"ytb-downloader/internal/resource"
	"ytb-downloader/internal/ui/component"
)

// Go GC is non-moving so it is safe to use pointer as map key
var active = make(map[*request.Request]fyne.Window)

func OpenRequestDebugViewer(app fyne.App, req *request.Request) fyne.Window {
	if win, ok := active[req]; ok {
		win.RequestFocus()
		return win
	}

	win := app.NewWindow("Debug")
	win.SetContent(container.NewVBox(content(req, win)))
	win.SetFixedSize(true)
	win.SetPadded(true)
	win.SetIcon(resource.ProgramIcon)
	win.Resize(fyne.NewSize(constants.RequestDebugWindowWidth, constants.RequestDebugWindowHeight))
	win.Show()
	win.SetOnClosed(func() {
		delete(active, req)
	})
	active[req] = win
	return win
}

func content(req *request.Request, win fyne.Window) fyne.CanvasObject {
	titleEntry := component.NewCopyableLabel(req.Title(), win)
	urlEntry := component.NewCopyableLabel(req.RawUrl(), win)
	formatEntry := component.NewCopyableLabel(req.Format(), win)
	titleFetchCmdEntry := component.NewWrappedCopyableLabel(req.GetTitleFetchCommand(), win, 120)
	downloadCmdEntry := component.NewWrappedCopyableLabel(req.GetDownloadCommand(), win, 200)
	errorLogEntry := component.NewWrappedCopyableLabel(fmt.Sprint(req.DownloadError()), win, 120)

	return container.New(
		layout.NewFormLayout(),
		widget.NewLabel("Video Title:"), titleEntry,
		widget.NewLabel("Video URL:"), urlEntry,
		widget.NewLabel("Format:"), formatEntry,
		widget.NewLabel("Title-Fetch Command:"), titleFetchCmdEntry,
		widget.NewLabel("Download Command:"), downloadCmdEntry,
		widget.NewLabel("Error Log:"), errorLogEntry,
	)
}
