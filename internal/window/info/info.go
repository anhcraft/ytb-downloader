package info

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"os"
	"runtime"
	"time"
	"ytb-downloader/internal/constants"
	"ytb-downloader/internal/resource"
	"ytb-downloader/internal/ui/component"
	"ytb-downloader/internal/window"
)

type dynamicStats struct {
	numGC, pauseTotal, alloc, sys, heapAlloc *widget.Label
}

var win fyne.Window

func OpenInfo(app fyne.App) fyne.Window {
	if win != nil {
		win.RequestFocus()
		return win
	}

	ds := &dynamicStats{
		numGC:      widget.NewLabel(""),
		pauseTotal: widget.NewLabel(""),
		alloc:      widget.NewLabel(""),
		sys:        widget.NewLabel(""),
		heapAlloc:  widget.NewLabel(""),
	}

	done := make(chan struct{})
	go updateStats(ds, done)

	win = app.NewWindow("Info")
	win.Resize(fyne.NewSize(constants.InfoWindowWidth, constants.InfoWindowHeight))
	win.SetContent(content(ds, win))
	win.SetFixedSize(true)
	win.SetPadded(true)
	win.SetIcon(resource.ProgramIcon)
	win.Show()
	win.SetOnClosed(func() {
		close(done)
		win = nil
	})

	return win
}

func content(ds *dynamicStats, win fyne.Window) fyne.CanvasObject {
	appVersionLabel := component.NewWrappedCopyableLabel(fyne.CurrentApp().Metadata().Version, win)
	execFileLabel := component.NewWrappedCopyableLabel(ExecutableFile(), win)

	osLabel := component.NewWrappedCopyableLabel(runtime.GOOS, win)
	archLabel := component.NewWrappedCopyableLabel(runtime.GOARCH, win)
	goVersionLabel := component.NewWrappedCopyableLabel(runtime.Version(), win)

	homepageBtn := widget.NewButton("Homepage", func() {
		err := window.OpenURL("https://github.com/anhcraft/ytb-downloader")
		if err != nil {
			dialog.ShowError(err, win)
			return
		}
	})
	homepageBtn.SetIcon(theme.HomeIcon())

	return container.NewVBox(
		container.NewHBox(
			layout.NewSpacer(),
			homepageBtn,
		),
		widget.NewForm(
			widget.NewFormItem("App Version", appVersionLabel),
			widget.NewFormItem("Executable File", execFileLabel),
		),
		widget.NewSeparator(),
		widget.NewForm(
			widget.NewFormItem("Operating System", osLabel),
			widget.NewFormItem("Architecture", archLabel),
			widget.NewFormItem("Go Version", goVersionLabel),
		),
		widget.NewSeparator(),
		widget.NewForm(
			widget.NewFormItem("GC Cycles", ds.numGC),
			widget.NewFormItem("Total GC Pause", ds.pauseTotal),
			widget.NewFormItem("Memory Allocated", ds.alloc),
			widget.NewFormItem("System Memory", ds.sys),
			widget.NewFormItem("Heap Allocated", ds.heapAlloc),
		),
	)
}

func updateStats(ds *dynamicStats, done <-chan struct{}) {
	ticker := time.NewTicker(time.Millisecond * 500)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			var memStats runtime.MemStats
			runtime.ReadMemStats(&memStats)

			// TODO thread-safe?
			ds.numGC.SetText(fmt.Sprintf("%d", memStats.NumGC))
			ds.pauseTotal.SetText(fmt.Sprintf("%d ms", memStats.PauseTotalNs/1e6))
			ds.alloc.SetText(fmt.Sprintf("%.2f MB", float64(memStats.Alloc)/1024/1024))
			ds.sys.SetText(fmt.Sprintf("%.2f MB", float64(memStats.Sys)/1024/1024))
			ds.heapAlloc.SetText(fmt.Sprintf("%.2f MB", float64(memStats.HeapAlloc)/1024/1024))
		case <-done:
			return
		}
	}
}

func ExecutableFile() string {
	exec, err := os.Executable()
	if err != nil {
		return ""
	}
	return exec
}
