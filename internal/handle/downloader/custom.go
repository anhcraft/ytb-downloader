package downloader

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
	"ytb-downloader/internal/handle/request"
)

func DownloadFile(req *request.Request, reportProgress func(*request.Progress)) error {
	filePath := req.FilePath()
	fi, err := os.Stat(filePath)
	if err == nil {
		if fi.IsDir() {
			return fmt.Errorf("path %s exists but is a directory", filePath)
		}
	} else if os.IsNotExist(err) {
		dir := filepath.Dir(filePath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create parent directories for %s: %v", filePath, err)
		}
	} else {
		return fmt.Errorf("error checking path %s: %v", filePath, err)
	}

	// Create or truncate the file
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	resp, err := http.Get(req.RawUrl())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	totalSize := resp.ContentLength
	var downloaded int64
	startTime := time.Now()
	lastReportTime := startTime
	buf := make([]byte, 32*1024)

	for {
		if req.Status() == request.StatusTerminated {
			return errors.New("download terminated by user")
		}

		n, readErr := resp.Body.Read(buf)
		if n > 0 {
			_, err := file.Write(buf[:n])
			if err != nil {
				return err
			}
			downloaded += int64(n)

			now := time.Now()
			if now.Sub(lastReportTime) >= 100*time.Millisecond || readErr != nil {
				progress := createProgress(downloaded, totalSize, startTime)
				reportProgress(progress)
				lastReportTime = now
			}
		}

		if readErr != nil {
			if readErr == io.EOF {
				break
			}
			return readErr
		}
	}

	progress := createProgress(downloaded, totalSize, startTime)
	reportProgress(progress)

	return nil
}

func createProgress(downloaded, totalSize int64, startTime time.Time) *request.Progress {
	progressPercentage := "NA"
	if totalSize > 0 {
		percent := float64(downloaded) / float64(totalSize) * 100
		progressPercentage = fmt.Sprintf("%.2f%%", percent)
	}

	downloadedSize := formatBytes(downloaded)
	totalSizeStr := formatBytes(totalSize)

	elapsed := time.Since(startTime)
	speedBps := 0.0
	if elapsed.Seconds() > 0 {
		speedBps = float64(downloaded) / elapsed.Seconds()
	}
	speedStr := fmt.Sprintf("%s/s", formatBytes(int64(speedBps)))

	etaStr := "NA"
	if totalSize > 0 && downloaded <= totalSize && speedBps > 0 {
		remaining := totalSize - downloaded
		etaSeconds := float64(remaining) / speedBps
		etaDuration := time.Duration(etaSeconds) * time.Second
		etaStr = formatDuration(etaDuration)
	}

	return &request.Progress{
		DownloadProgress:  progressPercentage,
		DownloadedSize:    downloadedSize,
		DownloadTotalSize: totalSizeStr,
		DownloadSpeed:     speedStr,
		DownloadEta:       etaStr,
	}
}

func formatBytes(bytes int64) string {
	if bytes < 0 {
		return "NA"
	}
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func formatDuration(d time.Duration) string {
	d = d.Round(time.Second)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second
	if h > 0 {
		return fmt.Sprintf("%dh%02dm%02ds", h, m, s)
	} else if m > 0 {
		return fmt.Sprintf("%dm%02ds", m, s)
	}
	return fmt.Sprintf("%ds", s)
}
