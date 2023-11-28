package settings

import (
	"os"
	"path/filepath"
	"strings"
	"ytb-downloader/internal/format"
	"ytb-downloader/internal/thumbnail"
)

type Settings struct {
	Format              string `json:"format,omitempty"`
	EmbedThumbnail      string `json:"embedThumbnail,omitempty"`
	DownloadFolder      string `json:"downloadFolder,omitempty"`
	YTdlpPath           string `json:"ytdlpPath,omitempty"`
	FFmpegPath          string `json:"ffmpegPath,omitempty"`
	ConcurrentDownloads int    `json:"concurrentDownloads,omitempty"`
	ConcurrentFragments int    `json:"concurrentFragments,omitempty"`
}

func (s *Settings) Normalize() {
	s.YTdlpPath = strings.TrimSpace(s.YTdlpPath)
	s.FFmpegPath = strings.TrimSpace(s.FFmpegPath)
	s.DownloadFolder = strings.TrimSpace(s.DownloadFolder)
	if !format.IsValid(s.Format) {
		s.Format = format.Default
	}
	if !thumbnail.IsValid(s.EmbedThumbnail) {
		s.EmbedThumbnail = thumbnail.AudioOnly
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

func (s *Settings) GetYTdlpPath() string {
	if s.YTdlpPath == "" {
		return "./yt-dlp.exe"
	}
	_, err := os.Stat(s.YTdlpPath)
	if err != nil {
		return "./yt-dlp.exe"
	}
	return s.YTdlpPath
}

func (s *Settings) GetFFmpegPath() string {
	if s.FFmpegPath == "" {
		return "./ffmpeg.exe"
	}
	_, err := os.Stat(s.FFmpegPath)
	if err != nil {
		return "./ffmpeg.exe"
	}
	return s.FFmpegPath
}

func NewSettings() *Settings {
	return &Settings{
		Format:              format.Default,
		EmbedThumbnail:      thumbnail.AudioOnly,
		DownloadFolder:      "",
		YTdlpPath:           "",
		FFmpegPath:          "",
		ConcurrentDownloads: 1,
		ConcurrentFragments: 3,
	}
}
