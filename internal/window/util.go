package window

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func OpenExplorer(path string) {
	if runtime.GOOS == "windows" {
		absPath, err := filepath.Abs(path)
		if err != nil {
			fmt.Println("Error resolving absolute path:", err)
			return
		}

		if _, err := os.Stat(absPath); os.IsNotExist(err) {
			fmt.Println("File or folder does not exist:", absPath)
			return
		}

		cmd := exec.Command("explorer", "/select,", absPath)
		_ = cmd.Run()
	}
}

// OpenURL https://stackoverflow.com/a/39324149
func OpenURL(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}
