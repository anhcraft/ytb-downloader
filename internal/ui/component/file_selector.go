package component

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"os"
)

func OpenFolderSelector(defaultLocation string, callback func(fyne.ListableURI, error), parent fyne.Window) {
	if len(defaultLocation) == 0 {
		defaultLocation, _ = os.UserHomeDir()
	}

	selector := dialog.NewFolderOpen(callback, parent)
	uri := storage.NewFileURI(defaultLocation)

	if listableUri, err := storage.ListerForURI(uri); err == nil {
		selector.SetLocation(listableUri)
	}

	selector.Show()
}
