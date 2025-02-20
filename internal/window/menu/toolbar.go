package menu

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"ytb-downloader/internal/handle/request"
	"ytb-downloader/internal/resource"
	"ytb-downloader/internal/window/info"
	settingsWindow "ytb-downloader/internal/window/settings"
)

func toolbar(app fyne.App) fyne.CanvasObject {
	toolbar := widget.NewToolbar(
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(resource.EraserIcon, func() {
			table.ScrollToTop()
			request.GetTable().Clear()
			table.Refresh()
		}),
		widget.NewToolbarAction(theme.SettingsIcon(), func() {
			settingsWindow.OpenSettings(app)
		}),
		widget.NewToolbarAction(theme.InfoIcon(), func() {
			info.OpenInfo(app)
		}),
	)
	return toolbar
}
