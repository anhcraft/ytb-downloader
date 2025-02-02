package request

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

type Progress struct {
	DownloadProgress  string
	DownloadedSize    string
	DownloadTotalSize string
	DownloadSpeed     string
	DownloadEta       string
}

func ExtractProgress(input string) (Progress, bool) {
	if input == "" || !strings.HasPrefix(input, "[[PROGRESS]]") {
		return Progress{}, false
	}

	input = strings.TrimPrefix(input, "[[PROGRESS]]")
	input = strings.TrimSpace(input)

	parts := strings.Split(input, ",")
	if len(parts) != 5 {
		return Progress{}, false
	}

	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}

	return Progress{
		DownloadProgress:  parts[0],
		DownloadedSize:    parts[1],
		DownloadTotalSize: parts[2],
		DownloadSpeed:     parts[3],
		DownloadEta:       parts[4],
	}, true
}

func hash(link string) string {
	sha := sha256.New()
	sha.Write([]byte(link))
	return hex.EncodeToString(sha.Sum(nil))
}
