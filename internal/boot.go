package internal

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"ytb-downloader/internal/handle"
	"ytb-downloader/internal/settings"
	"ytb-downloader/internal/ui/theme"
	"ytb-downloader/internal/window/menu"
)

var myApp fyne.App

func Init() {
	settings.Load()
	handle.InitLogger()
	myApp = app.New()
	myApp.Settings().SetTheme(&theme.CustomTheme{})
	menu.OpenMenu(myApp)
}
