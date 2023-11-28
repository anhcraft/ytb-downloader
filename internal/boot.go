package internal

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"ytb-downloader/internal/settings"
	"ytb-downloader/internal/ui/theme"
	"ytb-downloader/internal/window"
)

var myApp fyne.App

func Init() {
	settings.Load()
	myApp = app.New()
	myApp.Settings().SetTheme(&theme.CustomTheme{})
	window.OpenMenu(myApp)
}
