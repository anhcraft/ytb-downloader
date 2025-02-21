package menu

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"ytb-downloader/internal/constants/format"
	"ytb-downloader/internal/handle/request"
	"ytb-downloader/internal/resource"
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

	fetchBtn := widget.NewButtonWithIcon("Fetch", theme.SearchIcon(), func() {
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

	downloadBtn := widget.NewButtonWithIcon("Download", theme.DownloadIcon(), func() {
		request.SupplyQueue(request.GetTable().GetAllByStatus(request.StatusInQueue))
	})

	clearBtn := widget.NewButtonWithIcon("Clear", resource.EraserIcon, func() {
		table.ScrollToTop()
		request.GetTable().Clear()
		table.Refresh()
	})

	return container.NewVBox(
		layout.NewSpacer(),
		container.New(
			layout.NewFormLayout(),
			fmtLabel, fmtSelector,
		),
		container.NewHBox(
			layout.NewSpacer(),
			fetchBtn,
			downloadBtn,
			clearBtn,
		),
	)
}
