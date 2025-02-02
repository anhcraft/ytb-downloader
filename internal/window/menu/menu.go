package menu

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/rs/zerolog"
	"strings"
	"ytb-downloader/internal/constants"
	"ytb-downloader/internal/format"
	"ytb-downloader/internal/handle/request"
	"ytb-downloader/internal/resource"
	"ytb-downloader/internal/settings"
	"ytb-downloader/internal/ui/component"
	layout2 "ytb-downloader/internal/ui/layout"
	settingsWindow "ytb-downloader/internal/window/settings"
)

// NOTE on concurrency model:
// The whole UI is rendered by the main thread
// While the download scheduler and worker are independent goroutines
// We accept data race here... as UI updates are purely for visual

var win fyne.Window
var table *widget.Table
var input *widget.Entry
var logger zerolog.Logger

func OpenMenu(app fyne.App) fyne.Window {
	request.GetQueue().SetUpdateCallback(func(req *request.Request) {
		// TODO thread-safe?
		if req.Status() == request.StatusFailed {
			input.SetText(input.Text + "\n" + req.RawUrl())
		}
		table.Refresh()
	})

	logger = settings.Get().GetLogger().With().Str("scope", "gui/menu").Logger()
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

func toolbar(app fyne.App) fyne.CanvasObject {
	toolbar := widget.NewToolbar(
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(resource.EraserIcon, func() {
			request.GetTable().Clear()
			table.Refresh()
		}),
		widget.NewToolbarAction(theme.SettingsIcon(), func() {
			settingsWindow.OpenSettings(app)
		}),
	)
	return toolbar
}

func requestInput() fyne.CanvasObject {
	input = widget.NewMultiLineEntry()
	input.SetPlaceHolder("Enter URL(s) of videos, playlists, etc")
	input.SetMinRowsVisible(5)
	return input
}

func truncateString(s string, max int) string {
	if len(s) > max {
		return s[:max]
	}
	return s
}

func requestSettings() fyne.CanvasObject {
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
			downloadToLabel, downloadFolder,
		),
		container.NewHBox(
			layout.NewSpacer(),
			fetchBtn,
			downloadBtn,
		),
	)
}

func requestTable() fyne.CanvasObject {
	table = widget.NewTable(
		func() (int, int) {
			return request.GetTable().Count(), 2
		},
		func() fyne.CanvasObject {
			label := NewTableEntry()
			label.Truncation = fyne.TextTruncateEllipsis
			return label
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			entry := o.(*TableEntry)
			entry.req = request.GetTable().Get(i.Row)
			if i.Col == 0 {
				entry.Alignment = fyne.TextAlignLeading
				entry.SetText(entry.req.Title())
			} else {
				entry.Alignment = fyne.TextAlignCenter
				if entry.req.Status() == request.StatusDownloading && len(entry.req.DownloadProgress()) > 0 {
					entry.SetText(
						fmt.Sprintf(
							"Downloading %s (%s/%s) at %s ETA %s",
							entry.req.DownloadProgress(),
							entry.req.DownloadedSize(),
							entry.req.DownloadTotalSize(),
							entry.req.DownloadSpeed(),
							entry.req.DownloadEta(),
						),
					)
				} else {
					entry.SetText(entry.req.DescribeStatus())
				}
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

	table.SetColumnWidth(0, constants.MainWindowWidth*0.6)
	table.SetColumnWidth(1, constants.MainWindowWidth*0.38)

	return table
}

type TableEntry struct {
	widget.Label

	menu *fyne.Menu
	req  *request.Request
}

func (b *TableEntry) TappedSecondary(e *fyne.PointEvent) {
	widget.ShowPopUpMenuAtPosition(b.menu, fyne.CurrentApp().Driver().CanvasForObject(b), e.AbsolutePosition)
}

func NewTableEntry() *TableEntry {
	entry := &TableEntry{}
	entry.menu = fyne.NewMenu("",
		fyne.NewMenuItem("Copy URL", func() {
			win.Clipboard().SetContent(entry.req.RawUrl())
		}),
		fyne.NewMenuItem("Copy Title", func() {
			win.Clipboard().SetContent(entry.req.Title())
		}),
		fyne.NewMenuItem("Copy Title-Fetch command", func() {
			win.Clipboard().SetContent(settings.Get().GetYTdlpPath() + " " + strings.Join(entry.req.TitleFetchCmdArgs(), " "))
		}),
		fyne.NewMenuItem("Copy Download command", func() {
			win.Clipboard().SetContent(settings.Get().GetYTdlpPath() + " " + strings.Join(entry.req.DownloadCmdArgs(), " "))
		}),
		fyne.NewMenuItem("Terminate", func() {
			entry.req.SetStatus(request.StatusTerminated)
			table.Refresh()
		}),
		fyne.NewMenuItem("Remove", func() {
			if entry.req.Status() == request.StatusDownloading {
				return
			}
			request.GetTable().Remove(entry.req)
			table.Refresh()
		}),
	)
	entry.ExtendBaseWidget(entry)
	return entry
}
