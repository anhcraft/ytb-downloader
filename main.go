package main

import (
	"os"
	"ytb-downloader/internal"
)

func main() {
	for _, arg := range os.Args[1:] {
		if arg == "--gc" {
			internal.InitGcLog()
			break
		}
	}

	internal.Init()
}
