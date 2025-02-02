package downloader

import (
	"bufio"
	"errors"
	"io"
	"os/exec"
	"strings"
	"ytb-downloader/internal/handle/logger"
	"ytb-downloader/internal/handle/request"
	"ytb-downloader/internal/settings"
)

func download(req *request.Request, callback func()) {
	var err error

	defer func() {
		if err != nil {
			logger.Downloader.Println("error running command:", err)
			req.SetDownloadError(err)

			// SetStatus already check status transition
			// so if "Terminated"/"Completed" before, it would stay anyway
			req.SetStatus(request.StatusFailed)
			request.GetQueue().OnUpdate(req)
		}

		callback()
	}()

	logger.Downloader.Printf("downloading %s", req.RawUrl())

	cmd := exec.Command(settings.Get().GetYTdlpPath(), req.DownloadCmdArgs()...)
	request.DecorateCmd(cmd)
	logger.Downloader.Printf("executing command %s", cmd.String())

	stdout, err1 := cmd.StdoutPipe()
	if err1 != nil {
		err = err1
		return
	}

	stderr, err1 := cmd.StderrPipe()
	if err1 != nil {
		err = err1
		return
	}

	if err1 = cmd.Start(); err1 != nil {
		err = err1
		return
	}

	reader := bufio.NewReader(io.MultiReader(stdout, stderr))

	for {
		if req.Status() == request.StatusTerminated {
			err = errors.New("download terminated by user")
			break
		}

		line, err1 := reader.ReadString('\n')
		if err1 != nil {
			if err1 == io.EOF {
				break
			}
			logger.Downloader.Println("error reading line:", err1)
			break
		}

		// analyze log
		line = strings.TrimSpace(line)

		if p, ok := request.ExtractProgress(line); ok {
			req.SetDownloadProgress(p.DownloadProgress)
			req.SetDownloadedSize(p.DownloadedSize)
			req.SetDownloadTotalSize(p.DownloadTotalSize)
			req.SetDownloadSpeed(p.DownloadSpeed)
			req.SetDownloadEta(p.DownloadEta)
			request.GetQueue().OnUpdate(req)
		} else if strings.HasPrefix(line, "ERROR:") {
			err = errors.New(line)
		}

		logger.Downloader.Println("[yt-dlp]", line)
	}

	// if error occurred (on behalf of yt-dlp), the process returns non-zero exit code
	// so we prioritize the error already caught from the output
	if err == nil {
		if err1 = cmd.Wait(); err1 != nil {
			err = err1
			return
		}
	}

	if err == nil {
		req.SetStatus(request.StatusCompleted)
		request.GetQueue().OnUpdate(req)
	}
}
