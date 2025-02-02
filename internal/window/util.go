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
