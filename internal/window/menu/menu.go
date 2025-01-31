package menu

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/rs/zerolog"
	"image/color"
	"net/url"
	"strings"
	"ytb-downloader/internal/constants"
	"ytb-downloader/internal/format"
	"ytb-downloader/internal/handle"
	"ytb-downloader/internal/resource"
	"ytb-downloader/internal/settings"
	"ytb-downloader/internal/ui/component"
	layout2 "ytb-downloader/internal/ui/layout"
	settingsWindow "ytb-downloader/internal/window/settings"
)

var win fyne.Window
var table *widget.Table
var input *widget.Entry
var progressBar binding.Float
var logger zerolog.Logger

func OpenMenu(app fyne.App) fyne.Window {
	logger = settings.Get().GetLogger().With().Str("scope", "gui/menu").Logger()
	win = app.NewWindow("Yt-dlp GUI")
	ctn := container.NewVBox(
		header(app),
		container.New(layout2.NewHLayout(2, 0.35, 0.65), leftSide(), rightSide()),
		footer(),
	)
	win.SetContent(ctn)
	win.Resize(fyne.NewSize(constants.MainWindowWidth, constants.MainWindowHeight))
	win.SetFixedSize(true)
	win.SetPadded(true)
	win.SetIcon(resource.ProgramIcon)
	win.SetMaster()
	win.CenterOnScreen()
	win.ShowAndRun()
	return win
}

func appendUrl(s string) {
	input.SetText(input.Text + "\n" + s)
}

func header(app fyne.App) fyne.CanvasObject {
	toolbar := widget.NewToolbar(
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(resource.EraserIcon, func() {
			_ = progressBar.Set(0)
			handle.ClearProcesses()
		}),
		widget.NewToolbarAction(theme.SettingsIcon(), func() {
			settingsWindow.OpenSettings(app)
		}),
	)
	return toolbar
}

func footer() fyne.CanvasObject {
	bg := canvas.NewRectangle(color.RGBA{R: 220, G: 220, B: 220, A: 255})
	bg.SetMinSize(fyne.NewSize(1, 5))
	fg := canvas.NewRectangle(color.RGBA{R: 86, G: 186, B: 245, A: 255})
	fg.SetMinSize(fyne.NewSize(1, 5))
	cont := container.New(layout2.NewZLayout(2), bg, fg)
	progressBar = binding.NewFloat()
	progressBar.AddListener(binding.NewDataListener(func() {
		v, _ := progressBar.Get()
		cont.Layout.(*layout2.ZLayout).SetSize(1, float32(v))
		cont.Refresh()
	}))
	_ = progressBar.Set(0)
	return cont
}

func leftSide() fyne.CanvasObject {
	space := canvas.NewRectangle(color.Transparent)
	space.SetMinSize(fyne.NewSize(0, 30))
	sep := canvas.NewLine(color.Gray16{Y: 0xcccc})
	cont := container.NewVBox(
		topLeft(),
		space,
		sep,
		space,
		bottomLeft(),
		space,
	)
	return cont
}

func topLeft() fyne.CanvasObject {
	input = widget.NewMultiLineEntry()
	input.SetPlaceHolder("Enter URL(s) of videos, playlists, etc")
	input.SetMinRowsVisible(16)

	btn := widget.NewButton("Fetch", func() {
		for _, v := range strings.Split(input.Text, "\n") {
			v = strings.TrimSpace(v)
			if _, err := url.ParseRequestURI(v); err == nil {
				handle.SubmitUrl(v, settings.Get().Format, func() {
					table.Refresh()
				})
			} else {
				logger.Printf("invalid URL: %s\n", v)
			}
		}
		input.SetText("")
		table.Refresh()
	})
	btn.SetIcon(theme.SearchIcon())

	return container.NewVBox(input, btn)
}

func truncateString(s string, max int) string {
	if len(s) > max {
		return s[:max]
	}
	return s
}

func bottomLeft() fyne.CanvasObject {
	fmtLabel := widget.NewLabel("Format")
	fmtSelector := widget.NewSelect(
		[]string{format.Default, format.VideoOnly, format.AudioOnly},
		func(value string) {
			settings.Get().Format = value
			settings.Save()
		})
	fmtSelector.SetSelected(settings.Get().Format)

	downloadToLabel := widget.NewLabel("Download To")
	downloadTo := widget.NewLabel(truncateString(settings.Get().DownloadFolder, 30))
	downloadFolder := container.NewHBox(
		downloadTo,
		layout.NewSpacer(),
		widget.NewButton("...", func() {
			component.OpenFolderSelector(settings.Get().DownloadFolder, func(uri fyne.ListableURI, err error) {
				if uri != nil {
					settings.Get().DownloadFolder = uri.Path()
					settings.Save()
					downloadTo.SetText(truncateString(uri.Path(), 30))
				}
			}, win)
		}),
	)

	downloadBtn := widget.NewButton("Download", func() {
		if !handle.Download(func(progress float64) {
			_ = progressBar.Set(progress)
			table.Refresh()
		}, func(err error, url string) {
			dialog.ShowError(err, win)
			appendUrl(url) // append error url
			table.Refresh()
		}, func() {
			table.Refresh()
		}, settings.Get().Format) {
			dialog.ShowInformation("Warning", "Downloading... Please wait!", win)
		}
	})
	downloadBtn.SetIcon(theme.DownloadIcon())

	grid := container.New(
		layout.NewFormLayout(),
		fmtLabel, fmtSelector,
		downloadToLabel, downloadFolder,
		layout.NewSpacer(), downloadBtn,
	)
	return grid
}

func rightSide() fyne.CanvasObject {
	table = widget.NewTable(
		func() (int, int) {
			return handle.CountProcess(), 2
		},
		func() fyne.CanvasObject {
			label := NewTableEntry()
			label.Truncation = fyne.TextTruncateEllipsis
			return label
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			p := handle.GetProcess(i.Row)
			entry := o.(*TableEntry)
			entry.url = p.URL
			if i.Col == 0 {
				entry.Alignment = fyne.TextAlignLeading
				entry.SetText(p.Name)
			} else {
				entry.Alignment = fyne.TextAlignCenter
				entry.SetText(p.Status.String())
			}
		})
	table.ShowHeaderRow = true
	table.UpdateHeader = func(id widget.TableCellID, o fyne.CanvasObject) {
		l := o.(*widget.Label)
		l.Alignment = fyne.TextAlignCenter
		if id.Row < 0 {
			if id.Col == 0 {
				l.SetText("Name")
			}
			if id.Col == 1 {
				l.SetText("Status")
			}
		}
	}
	table.SetColumnWidth(0, 520)
	table.SetColumnWidth(1, 120)
	return table
}

type TableEntry struct {
	widget.Label
	menu *fyne.Menu
	url  string
}

func (b *TableEntry) TappedSecondary(e *fyne.PointEvent) {
	widget.ShowPopUpMenuAtPosition(b.menu, fyne.CurrentApp().Driver().CanvasForObject(b), e.AbsolutePosition)
}

func NewTableEntry() *TableEntry {
	b := &TableEntry{}
	b.menu = fyne.NewMenu("",
		fyne.NewMenuItem("Copy URL", func() {
			win.Clipboard().SetContent(b.url)
		}),
		fyne.NewMenuItem("Copy Title", func() {
			win.Clipboard().SetContent(b.Text)
		}),
	)
	b.ExtendBaseWidget(b)
	return b
}
