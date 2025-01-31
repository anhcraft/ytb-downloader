package component

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"os"
	"path/filepath"
)

func getDir(path string) string {
	info, err := os.Stat(path)
	if err != nil {
		return ""
	}

	// path could miss the ending slash, leading to wrong filepath.Dir result
	// so we check if such path is actually a dir first
	if info.IsDir() {
		return path
	}

	parentDir := filepath.Dir(path)
	return parentDir
}

func OpenFileSelector(defaultLocation string, callback func(fyne.URIReadCloser, error), parent fyne.Window) {
	if len(defaultLocation) == 0 {
		defaultLocation, _ = os.UserHomeDir()
	}

	selector := dialog.NewFileOpen(callback, parent)
	defaultLocation = getDir(defaultLocation)
	uri := storage.NewFileURI(defaultLocation)

	if listableUri, err := storage.ListerForURI(uri); err == nil {
		selector.SetLocation(listableUri)
	}

	selector.Show()
}

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
