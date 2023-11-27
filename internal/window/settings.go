package window

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"ytb-downloader/internal/format"
	"ytb-downloader/internal/resource"
	"ytb-downloader/internal/settings"
)

func OpenSettings(app fyne.App) fyne.Window {
	win = app.NewWindow("Settings")
	ctn := container.NewVBox(settingsContainer())
	win.SetContent(ctn)
	win.Resize(fyne.NewSize(600, 400))
	win.SetFixedSize(true)
	win.SetPadded(true)
	win.SetIcon(resource.ProgramIcon)
	//win.CenterOnScreen()
	win.Show()
	return win
}

func settingsContainer() fyne.CanvasObject {
	ffmpegLabel := widget.NewLabel("FFmpeg Path")
	ffmpegPath := widget.NewLabel(settings.Get().FFmpegPath)
	ffmpegSelector := container.NewHBox(
		ffmpegPath,
		layout.NewSpacer(),
		widget.NewButton("...", func() {
			dialog.ShowFileOpen(func(uri fyne.URIReadCloser, err error) {
				if uri != nil {
					settings.Get().FFmpegPath = uri.URI().Path()
					settings.Save()
					ffmpegPath.SetText(uri.URI().Path())
				}
			}, win)
		}),
	)

	concurrentDownloads := binding.NewFloat()
	_ = concurrentDownloads.Set(float64(settings.Get().ConcurrentDownloads))
	concurrentDownloads.AddListener(binding.NewDataListener(func() {
		if v, e := concurrentDownloads.Get(); e == nil {
			settings.Get().ConcurrentDownloads = int(v)
			settings.Save()
		}
	}))
	concurrentDownloadsLabel := widget.NewLabel("Concurrent Downloads")
	concurrentDownloadsValue := widget.NewLabelWithData(binding.FloatToStringWithFormat(concurrentDownloads, "%.0f"))
	concurrentDownloadsSelector := container.NewBorder(
		nil,
		nil,
		concurrentDownloadsValue,
		nil,
		widget.NewSliderWithData(1, 5, concurrentDownloads),
	)

	concurrentFragments := binding.NewFloat()
	_ = concurrentFragments.Set(float64(settings.Get().ConcurrentFragments))
	concurrentFragments.AddListener(binding.NewDataListener(func() {
		if v, e := concurrentFragments.Get(); e == nil {
			settings.Get().ConcurrentFragments = int(v)
			settings.Save()
		}
	}))
	concurrentFragmentsLabel := widget.NewLabel("Concurrent Fragments")
	concurrentFragmentsValue := widget.NewLabelWithData(binding.FloatToStringWithFormat(concurrentFragments, "%.0f"))
	concurrentFragmentsSelector := container.NewBorder(
		nil,
		nil,
		concurrentFragmentsValue,
		nil,
		widget.NewSliderWithData(1, 10, concurrentFragments),
	)

	thumbnailLabel := widget.NewLabel("Embed Thumbnail")
	thumbnailSelector := widget.NewSelect(
		[]string{format.Default, format.VideoOnly, format.AudioOnly},
		func(value string) {
			settings.Get().EmbedThumbnail = value
			settings.Save()
		})
	thumbnailSelector.SetSelected(settings.Get().EmbedThumbnail)

	return container.New(
		layout.NewFormLayout(),
		ffmpegLabel, ffmpegSelector,
		concurrentDownloadsLabel, concurrentDownloadsSelector,
		concurrentFragmentsLabel, concurrentFragmentsSelector,
		thumbnailLabel, thumbnailSelector,
	)
}
