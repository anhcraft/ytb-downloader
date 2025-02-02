package menu

import (
	"fyne.io/fyne/v2"
	"ytb-downloader/internal/handle/update"
)

func CheckUpdate(callback func(latest bool, currVer string, latestVer string, err error)) {
	go func() {
		currVer := fyne.CurrentApp().Metadata().Version
		latest, s, err := update.IsLatest(currVer)
		callback(latest, currVer, s, err)
	}()
}
