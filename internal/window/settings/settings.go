package settings

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"ytb-downloader/internal/constants"
	"ytb-downloader/internal/constants/downloadmode"
	"ytb-downloader/internal/constants/thumbnail"
	"ytb-downloader/internal/resource"
	"ytb-downloader/internal/settings"
	"ytb-downloader/internal/ui/component"
	"ytb-downloader/internal/window"
)

var win fyne.Window

func OpenSettings(app fyne.App) fyne.Window {
	if win != nil {
		win.RequestFocus()
		return win
	}

	win = app.NewWindow("Settings")
	win.SetContent(settingsContainer())
	win.Resize(fyne.NewSize(constants.SettingWindowWidth, constants.SettingWindowHeight))
	win.SetFixedSize(true)
	win.SetPadded(true)
	win.SetIcon(resource.ProgramIcon)
	//win.CenterOnScreen()
	win.Show()
	win.SetOnClosed(func() {
		settings.Save()
		win = nil
	})

	return win
}

func settingsContainer() fyne.CanvasObject {
	downloadFolderLabel := widget.NewLabel("Download Folder")
	downloadFolderInput := component.NewAutoSaveInput(settings.Get().GetDownloadFolder, func(val string) {
		settings.Get().SetDownloadFolder(val)
		settings.Save()
	}, requirePathIsFolder)
	downloadFolderSelector := container.NewBorder(
		nil,
		nil,
		nil,
		widget.NewButton("...", func() {
			component.OpenFolderSelector(settings.Get().GetDownloadFolder(), func(uri fyne.ListableURI, err error) {
				if uri != nil {
					downloadFolderInput.SetText(uri.Path())
				}
			}, win)
		}),
		downloadFolderInput,
	)

	ytdlpLabel := widget.NewLabel("Yt-dlp Path")
	ytdlpPathInput := component.NewAutoSaveInput(settings.Get().GetYtdlpPath, func(val string) {
		settings.Get().SetYtdlpPath(val)
		settings.Save()
	}, requirePathIsFile)
	ytdlpSelector := container.NewBorder(
		nil,
		nil,
		nil,
		widget.NewButton("...", func() {
			component.OpenFileSelector(settings.Get().GetYtdlpPath(), func(uri fyne.URIReadCloser, err error) {
				if uri != nil {
					ytdlpPathInput.SetText(uri.URI().Path())
				}
			}, win)
		}),
		ytdlpPathInput,
	)

	ffmpegLabel := widget.NewLabel("FFmpeg Path")
	ffmpegPathInput := component.NewAutoSaveInput(settings.Get().GetFfmpegPath, func(val string) {
		settings.Get().SetFfmpegPath(val)
		settings.Save()
	}, requirePathIsFile)
	ffmpegSelector := container.NewBorder(
		nil,
		nil,
		nil,
		widget.NewButton("...", func() {
			component.OpenFileSelector(settings.Get().GetFfmpegPath(), func(uri fyne.URIReadCloser, err error) {
				if uri != nil {
					ffmpegPathInput.SetText(uri.URI().Path())
				}
			}, win)
		}),
		ffmpegPathInput,
	)

	concurrentDownloads := binding.NewFloat()
	_ = concurrentDownloads.Set(float64(settings.Get().GetConcurrentDownloads()))
	concurrentDownloads.AddListener(binding.NewDataListener(func() {
		if v, e := concurrentDownloads.Get(); e == nil {
			settings.Get().SetConcurrentDownloads(uint32(v))
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
	_ = concurrentFragments.Set(float64(settings.Get().GetConcurrentFragments()))
	concurrentFragments.AddListener(binding.NewDataListener(func() {
		if v, e := concurrentFragments.Get(); e == nil {
			settings.Get().SetConcurrentFragments(uint32(v))
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

	disallowOverwriteLabel := widget.NewLabel("Disallow Overwrite")
	disallowOverwriteSelector := widget.NewSelect(
		[]string{downloadmode.Default, downloadmode.CustomDownloadOnly, downloadmode.YtdlpDownloadOnly},
		func(value string) {
			settings.Get().SetDisallowOverwrite(value)
			settings.Save()
		})
	disallowOverwriteSelector.SetSelected(settings.Get().GetDisallowOverwrite())

	thumbnailLabel := widget.NewLabel("Embed Thumbnail")
	thumbnailSelector := widget.NewSelect(
		[]string{thumbnail.Always, thumbnail.VideoOnly, thumbnail.AudioOnly, thumbnail.Never},
		func(value string) {
			settings.Get().SetEmbedThumbnail(value)
			settings.Save()
		})
	thumbnailSelector.SetSelected(settings.Get().GetEmbedThumbnail())

	logPathLabel := widget.NewLabel("Path to log file")
	logPathInput := component.NewAutoSaveInput(settings.Get().GetLogPath, func(val string) {
		settings.Get().SetLogPath(val)
		settings.Save()
	}, requirePathIsFile)
	logPathSelector := container.NewBorder(
		nil,
		nil,
		nil,
		widget.NewButton("...", func() {
			component.OpenFileSelector(settings.Get().GetLogPath(), func(uri fyne.URIReadCloser, err error) {
				if uri != nil {
					logPathInput.SetText(uri.URI().Path())
				}
			}, win)
		}),
		logPathInput,
	)

	extraYtpOptLabel := widget.NewLabel("Extra Yt-dlp options (space separated)")
	extraYtpOptInputBinding := binding.NewString()
	_ = extraYtpOptInputBinding.Set(settings.Get().GetExtraYtdlpOptions())
	extraYtpOptInputBinding.AddListener(binding.NewDataListener(func() {
		v, _ := extraYtpOptInputBinding.Get()
		settings.Get().SetExtraYtdlpOptions(v)
	}))
	extraYtpOptInput := widget.NewEntryWithData(extraYtpOptInputBinding)

	locateSettingFile := widget.NewButton("Locate settings file", func() {
		window.OpenExplorer(settings.SETTINGS_FILE)
	})
	locateSettingFile.SetIcon(theme.SearchIcon())

	clearSettings := widget.NewButton("Reset settings", func() {
		dialog.ShowConfirm("Reset Settings", "Are you sure you want to reset all settings?", func(b bool) {
			if b {
				settings.Reset()
				settings.Save()
				win.SetContent(settingsContainer()) // reload
			}
		}, win)
	})
	clearSettings.SetIcon(theme.DeleteIcon())

	scriptFileLabel := widget.NewLabel("Script File")
	scriptFileInput := component.NewAutoSaveInput(settings.Get().GetScriptFile, func(val string) {
		settings.Get().SetScriptFile(val)
		settings.Save()
	}, requirePathIsFileOrAbsent)
	scriptFileSelector := container.NewBorder(
		nil,
		nil,
		nil,
		widget.NewButton("...", func() {
			component.OpenFileSelector(settings.Get().GetScriptFile(), func(uri fyne.URIReadCloser, err error) {
				if uri != nil {
					scriptFileInput.SetText(uri.URI().Path())
				}
			}, win)
		}),
		scriptFileInput,
	)

	return container.NewVBox(
		container.NewHBox(
			layout.NewSpacer(),
			locateSettingFile,
			clearSettings,
		),
		container.New(
			layout.NewFormLayout(),
			downloadFolderLabel, downloadFolderSelector,
			ytdlpLabel, ytdlpSelector,
			ffmpegLabel, ffmpegSelector,
			concurrentDownloadsLabel, concurrentDownloadsSelector,
			concurrentFragmentsLabel, concurrentFragmentsSelector,
			disallowOverwriteLabel, disallowOverwriteSelector,
			thumbnailLabel, thumbnailSelector,
			logPathLabel, logPathSelector,
			extraYtpOptLabel, extraYtpOptInput,
			scriptFileLabel, scriptFileSelector,
		),
	)
}
