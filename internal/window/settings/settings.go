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
	"ytb-downloader/internal/settings"
	"ytb-downloader/internal/ui/component"
)

type activeWindow struct {
	window fyne.Window
	cfg    *settings.Settings
}

var active = make(map[string]*activeWindow)

func OpenSettings(app fyne.App, profile settings.Profile) fyne.Window {
	if win, ok := active[profile.Name]; ok {
		win.window.RequestFocus()
		return win.window
	}

	win := app.NewWindow("Settings | " + profile.Name)
	active[profile.Name] = &activeWindow{
		window: win,
		cfg:    settings.LoadSettings(profile),
	}

	win.SetContent(settingsContainer(app, win, profile))
	win.Resize(fyne.NewSize(constants.SettingWindowWidth, constants.SettingWindowHeight))
	win.SetFixedSize(true)
	win.SetPadded(true)
	win.SetIcon(theme.SettingsIcon())
	win.Show()
	win.SetOnClosed(func() {
		settings.SaveSettings(profile, active[profile.Name].cfg)
		delete(active, profile.Name)
	})

	return win
}

func settingsContainer(app fyne.App, win fyne.Window, profile settings.Profile) fyne.CanvasObject {
	cfg := active[profile.Name].cfg

	clearSettings := widget.NewButtonWithIcon("Reset", theme.DeleteIcon(), func() {
		dialog.ShowConfirm("Reset Settings", "Are you sure you want to reset all settings?", func(b bool) {
			if b {
				active[profile.Name].cfg = settings.ResetSettings(profile)
				win.SetContent(settingsContainer(app, win, profile)) // reload
			}
		}, win)
	})

	////////////////////////////////////////////////////////

	downloadFolderLabel := widget.NewLabel("Download Folder")
	downloadFolderInput := component.NewAutoSaveInput(cfg.GetDownloadFolder, func(val string) {
		cfg.SetDownloadFolder(val)
		settings.SaveSettings(profile, cfg)
	}, requirePathIsFolder)
	downloadFolderSelector := container.NewBorder(
		nil,
		nil,
		nil,
		widget.NewButton("...", func() {
			component.OpenFolderSelector(cfg.GetDownloadFolder(), func(uri fyne.ListableURI, err error) {
				if uri != nil {
					downloadFolderInput.SetText(uri.Path())
				}
			}, win)
		}),
		downloadFolderInput,
	)

	ytdlpLabel := widget.NewLabel("Yt-dlp Path")
	ytdlpPathInput := component.NewAutoSaveInput(cfg.GetYtdlpPath, func(val string) {
		cfg.SetYtdlpPath(val)
		settings.SaveSettings(profile, cfg)
	}, requirePathIsFile)
	ytdlpSelector := container.NewBorder(
		nil,
		nil,
		nil,
		widget.NewButton("...", func() {
			component.OpenFileSelector(cfg.GetYtdlpPath(), func(uri fyne.URIReadCloser, err error) {
				if uri != nil {
					ytdlpPathInput.SetText(uri.URI().Path())
				}
			}, win)
		}),
		ytdlpPathInput,
	)

	ffmpegLabel := widget.NewLabel("FFmpeg Path")
	ffmpegPathInput := component.NewAutoSaveInput(cfg.GetFfmpegPath, func(val string) {
		cfg.SetFfmpegPath(val)
		settings.SaveSettings(profile, cfg)
	}, requirePathIsFile)
	ffmpegSelector := container.NewBorder(
		nil,
		nil,
		nil,
		widget.NewButton("...", func() {
			component.OpenFileSelector(cfg.GetFfmpegPath(), func(uri fyne.URIReadCloser, err error) {
				if uri != nil {
					ffmpegPathInput.SetText(uri.URI().Path())
				}
			}, win)
		}),
		ffmpegPathInput,
	)

	concurrentDownloads := binding.NewFloat()
	_ = concurrentDownloads.Set(float64(cfg.GetConcurrentDownloads()))
	concurrentDownloads.AddListener(binding.NewDataListener(func() {
		if v, e := concurrentDownloads.Get(); e == nil {
			cfg.SetConcurrentDownloads(uint32(v))
			settings.SaveSettings(profile, cfg)
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
	_ = concurrentFragments.Set(float64(cfg.GetConcurrentFragments()))
	concurrentFragments.AddListener(binding.NewDataListener(func() {
		if v, e := concurrentFragments.Get(); e == nil {
			cfg.SetConcurrentFragments(uint32(v))
			settings.SaveSettings(profile, cfg)
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
			cfg.SetDisallowOverwrite(value)
			settings.SaveSettings(profile, cfg)
		})
	disallowOverwriteSelector.SetSelected(cfg.GetDisallowOverwrite())

	thumbnailLabel := widget.NewLabel("Embed Thumbnail")
	thumbnailSelector := widget.NewSelect(
		[]string{thumbnail.Always, thumbnail.VideoOnly, thumbnail.AudioOnly, thumbnail.Never},
		func(value string) {
			cfg.SetEmbedThumbnail(value)
			settings.SaveSettings(profile, cfg)
		})
	thumbnailSelector.SetSelected(cfg.GetEmbedThumbnail())

	logPathLabel := widget.NewLabel("Path to log file")
	logPathInput := component.NewAutoSaveInput(cfg.GetLogPath, func(val string) {
		cfg.SetLogPath(val)
		settings.SaveSettings(profile, cfg)
	}, requirePathIsFileOrAbsent)
	logPathSelector := container.NewBorder(
		nil,
		nil,
		nil,
		widget.NewButton("...", func() {
			component.OpenFileSelector(cfg.GetLogPath(), func(uri fyne.URIReadCloser, err error) {
				if uri != nil {
					logPathInput.SetText(uri.URI().Path())
				}
			}, win)
		}),
		logPathInput,
	)

	extraYtpOptLabel := widget.NewLabel("Extra Yt-dlp options (space separated)")
	extraYtpOptInputBinding := binding.NewString()
	_ = extraYtpOptInputBinding.Set(cfg.GetExtraYtdlpOptions())
	extraYtpOptInputBinding.AddListener(binding.NewDataListener(func() {
		v, _ := extraYtpOptInputBinding.Get()
		cfg.SetExtraYtdlpOptions(v)
		settings.SaveSettings(profile, cfg)
	}))
	extraYtpOptInput := widget.NewEntryWithData(extraYtpOptInputBinding)

	scriptFileLabel := widget.NewLabel("Script File")
	scriptFileInput := component.NewAutoSaveInput(cfg.GetScriptFile, func(val string) {
		cfg.SetScriptFile(val)
		settings.SaveSettings(profile, cfg)
	}, requirePathIsFileOrAbsent)
	scriptFileSelector := container.NewBorder(
		nil,
		nil,
		nil,
		widget.NewButton("...", func() {
			component.OpenFileSelector(cfg.GetScriptFile(), func(uri fyne.URIReadCloser, err error) {
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
