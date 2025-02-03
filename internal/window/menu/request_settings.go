package menu

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"ytb-downloader/internal/format"
	"ytb-downloader/internal/handle/request"
	"ytb-downloader/internal/settings"
)

func requestSettings() fyne.CanvasObject {
	fmtLabel := widget.NewLabel("Format")
	fmtSelector := widget.NewSelect(
		[]string{format.Default, format.VideoOnly, format.AudioOnly},
		func(value string) {
			settings.Get().SetFormat(value)
			settings.Save()
		})
	fmtSelector.SetSelected(settings.Get().GetFormat())

	fetchBtn := widget.NewButton("Fetch", func() {
		// The work is done async so we do not clear the input if it is in progress
		if FetchInput(input.Text, func(req []*request.Request) {
			request.FetchTitles(req, func() {
				// TODO thread-safe?
				table.Refresh()
			})
		}) {
			// we must immediately clear the input instead of clearing after the work is done
			//  because one could type input while work is in progress
			input.SetText("")
		}
	})
	fetchBtn.SetIcon(theme.SearchIcon())

	downloadBtn := widget.NewButton("Download", func() {
		request.SupplyQueue(request.GetTable().GetAllByStatus(request.StatusInQueue))
	})
	downloadBtn.SetIcon(theme.DownloadIcon())

	return container.NewVBox(
		container.New(
			layout.NewFormLayout(),
			fmtLabel, fmtSelector,
		),
		container.NewHBox(
			layout.NewSpacer(),
			fetchBtn,
			downloadBtn,
		),
	)
}
