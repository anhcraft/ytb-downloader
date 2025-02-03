package request

import (
	"net/url"
	"os"
	"os/exec"
	"strings"
	"ytb-downloader/internal/handle/logger"
	"ytb-downloader/internal/scripting"
	"ytb-downloader/internal/settings"
)

func ParseRequest(input string) []*Request {
	scriptCode := settings.Get().LoadScriptFile()
	res := make([]*Request, 0)

	for _, v := range strings.Split(input, "\n") {
		v = strings.TrimSpace(v)

		if scriptCode != nil {
			logger.Queue.Println("running script at input:", v)
			result, err := scripting.HandleDownload(scriptCode, v)

			if err != nil {
				logger.Queue.Println("error running script:", err)
				continue
			}

			switch result.Action {
			case "skip":
				continue
			case "override":
				v = result.Value
			default:
				// do nothing
			}
		}

		u, err := url.ParseRequestURI(v)

		if err != nil {
			continue
		}

		RewriteYoutubeShortLink(u)

		// flatten YouTube playlist
		if strings.HasSuffix(u.Host, "youtube.com") && u.Path == "/playlist" {
			flattenYoutubePlaylist(&res, v)
			continue
		}

		res = append(res, NewRequest(u))

	}

	return res
}

func flattenYoutubePlaylist(queue *[]*Request, link string) {
	logger.Queue.Println("new playlist link", link)

	// ./yt-dlp.exe --flat-playlist --ignore-errors --no-warnings --print-to-file "title,url" "temp.txt" ""
	// somehow printing into the console does not support UTF8
	// so the workaround is using a temporary file

	// hash the url because it might contain illegal characters (OS-dependent)
	temp, err := os.CreateTemp("", hash(link))
	if err != nil {
		logger.Queue.Println("error creating temp file:", err)
		return
	}

	defer func(temp *os.File) {
		if err := temp.Close(); err != nil {
			logger.Queue.Println("error closing temp file:", err)
		}
	}(temp)

	tempPath := temp.Name()
	cmd := exec.Command(
		settings.Get().GetYtdlpPath(),
		append(settings.Get().ExtraYtdlpOptionsAsArray(),
			"--skip-download",
			"--flat-playlist",
			"--ignore-errors",
			"--no-warnings",
			"--print-to-file",
			"url,title",
			tempPath,
			link,
		)...,
	)
	DecorateCmd(cmd)
	logger.Queue.Printf("executing command %s", cmd.String())

	if err := cmd.Run(); err != nil {
		logger.Queue.Println("error running command:", err)
		return
	}

	bytes, err := os.ReadFile(tempPath)

	if err != nil {
		logger.Queue.Println("error reading temp file:", err)
		return
	}

	lines := strings.Split(strings.ReplaceAll(string(bytes), "\r", ""), "\n")
	logger.Queue.Printf("found %d videos in the playlist", len(lines)>>1)

	for i := 0; i+1 < len(lines); i += 2 {
		logger.Queue.Println(lines[i+1] + ": " + lines[i])

		// TODO better way to check private videos
		if strings.Contains(lines[i+1], "[Private video]") {
			continue
		}

		if u, err := url.ParseRequestURI(lines[i]); err == nil {
			req := NewRequest(u)
			req.SetTitle(lines[i+1])
			req.SetTitleFetched(true)
			*queue = append(*queue, req)
		}
	}
}
