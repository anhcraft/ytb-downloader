package menu

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"ytb-downloader/internal/constants"
	"ytb-downloader/internal/handle/request"
	"ytb-downloader/internal/window/debug"
)

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

	table.SetColumnWidth(0, constants.MainWindowWidth*0.5)
	table.SetColumnWidth(1, constants.MainWindowWidth*0.48)

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
		fyne.NewMenuItem("Copy Download command", func() {
			win.Clipboard().SetContent(entry.req.GetDownloadCommand())
		}),
		fyne.NewMenuItem("Terminate", func() {
			entry.req.SetStatus(request.StatusTerminated)
			table.Refresh()
		}),
		fyne.NewMenuItem("Debug", func() {
			debug.OpenRequestDebugViewer(fyne.CurrentApp(), entry.req)
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
