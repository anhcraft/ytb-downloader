package handle

import (
	"bufio"
	"errors"
	"io"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"ytb-downloader/internal/format"
	"ytb-downloader/internal/settings"
)

var isDownloading bool
var downloadLock sync.Mutex

func Download(onUpdate func(progress float64), onError func(err error), onFinish func(), format string) bool {
	downloadLock.Lock()
	success := false
	if !isDownloading {
		isDownloading = true
		success = true
		// Copies the current queue of processes allowing for parallel processing
		// and adding URLs while processing
		processing := make([]*Process, 0)
		for _, v := range processes {
			if v.Status == Queued {
				// Override formats
				v.Format = format
				processing = append(processing, v)
			}
		}
		go _download(onUpdate, func() {
			isDownloading = false
			onFinish()
		}, onError, processing)
	}
	downloadLock.Unlock()
	return success
}

func _download(onUpdate func(progress float64), onFinish func(), onError func(err error), processing []*Process) {
	totalProgress := float64(len(processing) * 100)
	progress := float64(0)

	for _, v := range processing {
		func() {
			var err error
			defer func() {
				if err != nil {
					v.Status = Error
					onError(err)
				}
			}()

			log.Printf("downloading %s\n", v.URL)
			v.Status = Downloading
			onUpdate(progress / totalProgress)

			args := []string{"--ignore-errors", "--no-warnings",
				"--progress", "--newline",
				//"--progress-template", "%(progress)j",
				"--concurrent-fragments", strconv.Itoa(settings.Get().ConcurrentDownloads),
				"--abort-on-unavailable-fragments",
				"-P", settings.Get().GetDownloadFolder(),
				"--embed-thumbnail", "--embed-metadata"}

			if fp := settings.Get().GetFFmpegPath(); fp != "" {
				args = append(args, "--ffmpeg-location", fp)
			}

			if v.Format == format.VideoOnly {
				args = append(args, "-f", "bestvideo")
			} else if v.Format == format.AudioOnly {
				args = append(args, "-f", "bestaudio")
			}

			args = append(args, v.URL)

			cmd := exec.Command("./yt-dlp.exe", args...)
			log.Printf("Executing command %s\n", cmd.String())

			stdout, err1 := cmd.StdoutPipe()
			if err1 != nil {
				log.Println("Error creating StdoutPipe:", err1)
				err = err1
				onUpdate(progress / totalProgress)
				return
			}

			stderr, err1 := cmd.StderrPipe()
			if err1 != nil {
				log.Println("Error creating StderrPipe:", err1)
				err = err1
				onUpdate(progress / totalProgress)
				return
			}

			if err1 = cmd.Start(); err1 != nil {
				log.Println("Error starting command:", err1)
				err = err1
				onUpdate(progress / totalProgress)
				return
			}

			reader := bufio.NewReader(io.MultiReader(stdout, stderr))

			for {
				line, err1 := reader.ReadString('\n')
				if err1 != nil {
					if err1 == io.EOF {
						break
					}
					log.Println("Error reading line:", err1)
					break
				}

				// analyze log
				line = strings.TrimSpace(line)
				if p, ok := extractPercentage(line); ok {
					onUpdate((progress + p) / totalProgress)
				} else if strings.Contains(line, "ffmpeg not found") {
					err = errors.New("error due to ffmpeg not installed")
				} else {
					log.Println(line)
				}
			}

			if err1 = cmd.Wait(); err1 != nil && err == nil {
				log.Println("Error on running command:", err1)
				err = err1
				onUpdate(progress / totalProgress)
				return
			}

			v.Status = Done
			progress += 100
			onUpdate(progress / totalProgress)
		}()
	}

	onFinish()
}
