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
	"sync/atomic"
	"syscall"
	"ytb-downloader/internal/format"
	"ytb-downloader/internal/settings"
	"ytb-downloader/internal/thumbnail"
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
	var progress atomic.Int32
	var semaphore sync.WaitGroup
	semaphore.Add(len(processing))
	jobs := make(chan *Process, len(processing))

	for i := 0; i < settings.Get().ConcurrentDownloads; i++ {
		go func() {
			for job := range jobs {
				func() {
					var err error
					defer func() {
						if err != nil {
							job.Status = Error
							onError(err)
						}
						semaphore.Done()
					}()

					log.Printf("downloading %s\n", job.URL)
					job.Status = Downloading
					onUpdate(float64(progress.Load()) / totalProgress)

					args := []string{"--ignore-errors", "--no-warnings",
						"--progress", "--newline",
						//"--progress-template", "%(progress)j",
						"--concurrent-fragments", strconv.Itoa(settings.Get().ConcurrentFragments),
						"--abort-on-unavailable-fragments",
						"-P", settings.Get().GetDownloadFolder()}

					if fp := settings.Get().GetFFmpegPath(); fp != "" {
						args = append(args, "--ffmpeg-location", fp)
					}

					embedThumbnail := settings.Get().EmbedThumbnail
					shouldEmbedThumbnail := embedThumbnail != thumbnail.Never &&
						(embedThumbnail == thumbnail.Always || job.Format == embedThumbnail)

					if shouldEmbedThumbnail {
						args = append(args, "--embed-thumbnail")
					}

					// Choose the best quality format
					// Remux the video to mp4 or audio to m4a to support thumbnail embedding
					if job.Format == format.VideoOnly {
						args = append(args, "-f", "bestvideo")
						if shouldEmbedThumbnail {
							args = append(args, "--remux-video", "mp4")
						}
					} else if job.Format == format.AudioOnly {
						args = append(args, "-f", "bestaudio")
						if shouldEmbedThumbnail {
							args = append(args, "-x", "--audio-quality", "0", "--audio-format", "m4a")
						}
					} else if shouldEmbedThumbnail {
						args = append(args, "--merge-output-format", "mp4")
					}

					args = append(args, job.URL)

					cmd := exec.Command(settings.Get().GetYTdlpPath(), args...)
					cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
					log.Printf("Executing command %s\n", cmd.String())

					stdout, err1 := cmd.StdoutPipe()
					if err1 != nil {
						log.Println("Error creating StdoutPipe:", err1)
						err = err1
						onUpdate(float64(progress.Load()) / totalProgress)
						return
					}

					stderr, err1 := cmd.StderrPipe()
					if err1 != nil {
						log.Println("Error creating StderrPipe:", err1)
						err = err1
						onUpdate(float64(progress.Load()) / totalProgress)
						return
					}

					if err1 = cmd.Start(); err1 != nil {
						log.Println("Error starting command:", err1)
						err = err1
						onUpdate(float64(progress.Load()) / totalProgress)
						return
					}

					reader := bufio.NewReader(io.MultiReader(stdout, stderr))
					lastPercentage := int32(0)

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
							deltaPercentage := int32(p) - lastPercentage
							lastPercentage = int32(p)
							onUpdate(float64(progress.Add(deltaPercentage)) / totalProgress)
						} else if strings.Contains(line, "ffmpeg not found") {
							err = errors.New("error due to ffmpeg not installed")
						} else if strings.Contains(line, "Requested format is not available") {
							err = errors.New("requested format is not available")
						} else if strings.HasPrefix(line, "ERROR:") {
							err = errors.New(line)
						} else {
							log.Println("[yt-dlp]", line)
						}
					}

					defer func() {
						onUpdate(float64(progress.Add(100-lastPercentage)) / totalProgress)
					}()

					if err1 = cmd.Wait(); err1 != nil && err == nil {
						log.Println("Error on running command:", err1)
						err = err1
						onUpdate(float64(progress.Load()) / totalProgress)
						return
					}

					job.Status = Done
					onUpdate(float64(progress.Load()) / totalProgress)
				}()
			}
		}()
	}

	for _, v := range processing {
		jobs <- v
	}
	close(jobs)

	semaphore.Wait()
	onFinish()
}
