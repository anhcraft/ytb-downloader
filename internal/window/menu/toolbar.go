package menu

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"ytb-downloader/internal/handle/request"
	"ytb-downloader/internal/resource"
	"ytb-downloader/internal/settings"
	"ytb-downloader/internal/window/info"
	"ytb-downloader/internal/window/profile"
)

var profileSelector *widget.Select
var profileListRefreshTrigger func()

func RefreshProfileSelector() {
	profileSelector.SetOptions(settings.GetProfileNames())
	profileSelector.SetSelected(settings.GetProfile().Name)
}

func toolbar(app fyne.App) fyne.CanvasObject {
	profileSelector = widget.NewSelect(
		[]string{},
		func(value string) {
			settings.SelectProfile(value)
			if profileListRefreshTrigger != nil {
				profileListRefreshTrigger()
			}
		})
	RefreshProfileSelector()

	return container.NewHBox(
		layout.NewSpacer(),
		profileSelector,
		widget.NewToolbar(
			widget.NewToolbarAction(theme.GridIcon(), func() {
				profileListRefreshTrigger = profile.OpenProfile(app, func() {
					RefreshProfileSelector()
				}, func() {
					profileListRefreshTrigger = nil
				})
			}),
			widget.NewToolbarAction(resource.EraserIcon, func() {
				table.ScrollToTop()
				request.GetTable().Clear()
				table.Refresh()
			}),
			widget.NewToolbarAction(theme.InfoIcon(), func() {
				info.OpenInfo(app)
			}),
		),
	)
}
