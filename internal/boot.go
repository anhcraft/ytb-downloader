package internal

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"log"
	"runtime"
	"time"
	"ytb-downloader/internal/handle/downloader"
	"ytb-downloader/internal/handle/logger"
	"ytb-downloader/internal/settings"
	"ytb-downloader/internal/ui/theme"
	"ytb-downloader/internal/window/menu"
)

var myApp fyne.App

func Init() {
	settings.Load()
	logger.InitLogger()
	downloader.InitDownloadScheduler()

	myApp = app.New()
	myApp.Settings().SetTheme(&theme.CustomTheme{})
	menu.OpenMenu(myApp)
	myApp.Lifecycle().SetOnStopped(func() {
		downloader.TerminateDownloadScheduler()
	})
}

func InitGcLog() {
	go func() {
		for {
			time.Sleep(3 * time.Second)

			var memStats runtime.MemStats
			runtime.ReadMemStats(&memStats)

			log.Printf(
				"NumGC: %d, PauseTotal: %dms, Alloc: %d KB, Sys: %d KB, HeapAlloc: %d KB, HeapSys: %d KB",
				memStats.NumGC,
				memStats.PauseTotalNs/1e6,
				memStats.Alloc/1024,
				memStats.Sys/1024,
				memStats.HeapAlloc/1024,
				memStats.HeapSys/1024,
			)
		}
	}()
}
