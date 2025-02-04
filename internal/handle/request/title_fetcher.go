package request

import (
	"os"
	"os/exec"
	"strings"
	"ytb-downloader/internal/handle/logger"
	"ytb-downloader/internal/settings"
)

func FetchTitles(req []*Request, callback func()) {
	go func() {
		for _, r := range req {
			fetchTitle(r)
		}
		callback()
	}()
}

func fetchTitle(req *Request) bool {
	if req.TitleFetched() || req.Custom() {
		logger.Queue.Printf("SKIPPED fetching title for %s", req.RawUrl())
		return true
	}

	logger.Queue.Printf("fetching title for %s", req.RawUrl())

	// hash the url because it might contain illegal characters (OS-dependent)
	temp, err := os.CreateTemp("", hash(req.RawUrl()))

	if err != nil {
		logger.Queue.Println("error creating temp file:", err)
		return false
	}

	defer func(temp *os.File) {
		if err := temp.Close(); err != nil {
			logger.Queue.Println("error closing temp file:", err)
		}
	}(temp)

	tempPath := temp.Name()
	req.SetTitleFetchCmdArgs(append(settings.Get().ExtraYtdlpOptionsAsArray(),
		"--skip-download",
		"--ignore-errors",
		"--no-warnings",
		"--print-to-file",
		"title",
		tempPath,
		req.url.String(),
	))
	cmd := exec.Command(settings.Get().GetYtdlpPath(), req.TitleFetchCmdArgs()...)
	DecorateCmd(cmd)
	logger.Queue.Println("executing command", cmd.String())

	if err := cmd.Run(); err != nil {
		logger.Queue.Println("error running command:", err)
		return false
	}

	bytes, err := os.ReadFile(tempPath)
	if err != nil {
		logger.Queue.Println("error reading temp file:", err)
		return false
	}

	title := strings.ReplaceAll(string(bytes), "\r", "")
	title = strings.ReplaceAll(title, "\n", "")
	req.SetTitle(title)
	req.SetTitleFetched(true)
	return true
}
