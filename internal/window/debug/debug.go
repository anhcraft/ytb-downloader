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
	win.Resize(fyne.NewSize(constants.RequestDebugWindowWidth, constants.RequestDebugWindowHeight))
	win.SetContent(container.NewVBox(content(req, win)))
	win.SetFixedSize(true)
	win.SetPadded(true)
	win.SetIcon(resource.ProgramIcon)
	win.Show()
	win.SetOnClosed(func() {
		delete(active, req)
	})
	active[req] = win
	return win
}

func content(req *request.Request, win fyne.Window) fyne.CanvasObject {
	inputEntry := component.NewWrappedCopyableLabel(req.Input(), win)
	titleEntry := component.NewWrappedCopyableLabel(req.Title(), win)
	urlEntry := component.NewWrappedCopyableLabel(req.RawUrl(), win)
	formatEntry := component.NewWrappedCopyableLabel(req.Format(), win)
	errorLogEntry := component.NewWrappedCopyableLabel(fmt.Sprint(req.DownloadError()), win)

	if req.Custom() {
		filePathEntry := component.NewWrappedCopyableLabel(req.FilePath(), win)
		return container.New(
			layout.NewFormLayout(),
			widget.NewLabel("Input:"), inputEntry,
			widget.NewLabel("File Path:"), filePathEntry,
			widget.NewLabel("Video Title:"), titleEntry,
			widget.NewLabel("Video URL:"), urlEntry,
			widget.NewLabel("Format:"), formatEntry,
			widget.NewLabel("Error Log:"), errorLogEntry,
		)
	}

	titleFetchCmdEntry := component.NewWrappedCopyableLabel(req.GetTitleFetchCommand(), win)
	downloadCmdEntry := component.NewWrappedCopyableLabel(req.GetDownloadCommand(), win)

	return container.New(
		layout.NewFormLayout(),
		widget.NewLabel("Input:"), inputEntry,
		widget.NewLabel("Video Title:"), titleEntry,
		widget.NewLabel("Video URL:"), urlEntry,
		widget.NewLabel("Format:"), formatEntry,
		widget.NewLabel("Title-Fetch Command:"), titleFetchCmdEntry,
		widget.NewLabel("Download Command:"), downloadCmdEntry,
		widget.NewLabel("Error Log:"), errorLogEntry,
	)
}
