package settings

import (
	"os"
	"path/filepath"
	"ytb-downloader/internal/format"
)

type Settings struct {
	Format              string `json:"format,omitempty"`
	EmbedThumbnail      string `json:"embedThumbnail,omitempty"`
	DownloadFolder      string `json:"downloadFolder,omitempty"`
	FFmpegPath          string `json:"ffmpegPath,omitempty"`
	ConcurrentDownloads int    `json:"concurrentDownloads,omitempty"`
	ConcurrentFragments int    `json:"concurrentFragments,omitempty"`
}

func (s *Settings) Normalize() {
	if !format.IsValid(s.Format) {
		s.Format = format.Default
	}
	if !format.IsValid(s.EmbedThumbnail) {
		s.EmbedThumbnail = format.AudioOnly
	}
	if s.ConcurrentDownloads < 1 {
		s.ConcurrentDownloads = 1
	}
	if s.ConcurrentFragments < 1 {
		s.ConcurrentFragments = 1
	}
}

func (s *Settings) GetDownloadFolder() string {
	fi, err := os.Stat(s.DownloadFolder)
	if err != nil || !fi.IsDir() {
		pwd, err := os.Getwd()
		if err != nil {
			return "./downloads/"
		}
		return filepath.Join(pwd, "downloads")
	}
	return s.DownloadFolder
}

func (s *Settings) GetFFmpegPath() string {
	_, err := os.Stat(s.FFmpegPath)
	if err != nil {
		return ""
	}
	return s.FFmpegPath
}

func NewSettings() *Settings {
	return &Settings{
		Format:              format.Default,
		EmbedThumbnail:      format.AudioOnly,
		DownloadFolder:      "",
		FFmpegPath:          "",
		ConcurrentDownloads: 1,
		ConcurrentFragments: 3,
	}
}
