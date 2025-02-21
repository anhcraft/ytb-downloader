package profile

import (
	"errors"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"strings"
	"ytb-downloader/internal/constants"
	"ytb-downloader/internal/settings"
	"ytb-downloader/internal/ui/component"
	"ytb-downloader/internal/window"
	settingsWindow "ytb-downloader/internal/window/settings"
)

var win fyne.Window
var app fyne.App
var contentContainer *fyne.Container
var profileUpdateCallback func()
var profileListRefreshTrigger func()

func OpenProfile(_app fyne.App, _profileUpdateCallback func(), onClose func()) func() {
	if win != nil {
		win.RequestFocus()
		return profileListRefreshTrigger
	}

	profileUpdateCallback = _profileUpdateCallback
	app = _app
	win = _app.NewWindow("Profiles")
	contentContainer = container.NewPadded()
	switchToProfileList()
	win.SetContent(contentContainer)
	win.Resize(fyne.NewSize(constants.ProfileWindowWidth, constants.ProfileWindowHeight))
	win.SetFixedSize(true)
	win.SetPadded(true)
	win.SetIcon(theme.GridIcon())
	win.Show()
	win.CenterOnScreen()
	win.SetOnClosed(func() {
		app = nil
		win = nil
		profileUpdateCallback = nil
		profileListRefreshTrigger = nil
		onClose()
	})
	profileListRefreshTrigger = func() {
		if win != nil && contentContainer != nil {
			switchToProfileList()
		}
	}

	return profileListRefreshTrigger
}

func switchToProfileList() {
	contentContainer.Objects = []fyne.CanvasObject{profileListContent()}
	contentContainer.Refresh()
}

func switchToProfileAdd() {
	contentContainer.Objects = []fyne.CanvasObject{profileAddContent()}
	contentContainer.Refresh()
}

func profileListContent() fyne.CanvasObject {
	profileList := container.NewVBox()
	borderedList := container.NewPadded(container.NewVBox(profileList))
	border := container.NewVBox(widget.NewSeparator(), borderedList, widget.NewSeparator())

	selectedProfile := settings.GetProfile()

	for _, p := range settings.GetProfiles() {
		check := canvas.NewText("✔", color.RGBA{G: 255, A: 255})
		profileName := widget.NewLabel(p.Name)
		nameContainer := container.NewHBox(profileName, check)

		if p.Name != selectedProfile.Name {
			check.Hide()
		}

		locateBtn := widget.NewButtonWithIcon("", theme.SearchIcon(), func() {
			window.OpenExplorer(p.Path)
		})

		deleteBtn := widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
			if p.Name == selectedProfile.Name {
				dialog.ShowError(errors.New("cannot delete selected profile"), win)
				return
			}
			dialog.ShowConfirm("Delete Profile", "Are you sure you want to delete profile "+p.Name+"?", func(ok bool) {
				if ok {
					settings.DeleteProfile(p.Name)
					switchToProfileList()
					profileUpdateCallback()
				}
			}, win)
		})

		settingsBtn := widget.NewButtonWithIcon("", theme.SettingsIcon(), func() {
			settingsWindow.OpenSettings(app, p)
		})

		if p.Name == selectedProfile.Name {
			deleteBtn.Disable()
		}

		buttonGroup := container.NewHBox(layout.NewSpacer(), deleteBtn, locateBtn, settingsBtn, layout.NewSpacer())
		row := container.NewHBox(nameContainer, layout.NewSpacer(), buttonGroup)
		profileList.Add(row)
	}

	addProfileBtn := widget.NewButton("Add Profile", func() {
		switchToProfileAdd()
	})
	addProfileRow := container.NewHBox(layout.NewSpacer(), addProfileBtn)

	return container.NewVBox(addProfileRow, border)
}

func profileAddContent() fyne.CanvasObject {
	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("Enter Profile Name")
	pathEntry := widget.NewEntry()
	pathEntry.SetPlaceHolder("Enter Profile Path")
	pathEntrySelector := container.NewBorder(
		nil,
		nil,
		nil,
		widget.NewButton("...", func() {
			defaultPath := pathEntry.Text
			if len(defaultPath) == 0 {
				defaultPath = settings.GetProfile().Path
			}
			component.OpenFileSelector(defaultPath, func(uri fyne.URIReadCloser, err error) {
				if uri != nil {
					pathEntry.SetText(uri.URI().Path())
				}
			}, win)
		}),
		pathEntry,
	)

	confirmBtn := widget.NewButton("Confirm", func() {
		name := strings.TrimSpace(nameEntry.Text)
		path := strings.TrimSpace(pathEntry.Text)
		if name == "" {
			dialog.ShowError(errors.New("name cannot be empty"), win)
			return
		}
		if path == "" {
			dialog.ShowError(errors.New("path cannot be empty"), win)
			return
		}
		err := settings.AddProfile(settings.Profile{
			Name: name,
			Path: path,
		})
		if err == nil {
			switchToProfileList()
			profileUpdateCallback()
		} else {
			dialog.ShowError(err, win)
		}
	})

	backBtn := widget.NewButton("Back", func() {
		switchToProfileList()
	})

	title := canvas.NewText("Add New Profile", color.Black)
	title.TextSize = 24

	form := container.NewVBox(
		title,
		container.New(
			layout.NewFormLayout(),
			widget.NewLabel("Name:"), nameEntry,
			widget.NewLabel("Path:"), pathEntrySelector,
		),
		container.NewHBox(layout.NewSpacer(), confirmBtn, backBtn),
	)

	return container.NewVBox(form)
}
