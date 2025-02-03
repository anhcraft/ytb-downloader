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
	YtdlpPath           string `json:"ytdlpPath,omitempty"`
	FfmpegPath          string `json:"ffmpegPath,omitempty"`
	ConcurrentDownloads uint32 `json:"concurrentDownloads,omitempty"`
	ConcurrentFragments uint32 `json:"concurrentFragments,omitempty"`
	LogPath             string `json:"logPath,omitempty"`
	ExtraYtdlpOptions   string `json:"extraYtdlpOptions,omitempty"`
}

func NewSettings() *Settings {
	return &Settings{
		Format:              format.Default,
		EmbedThumbnail:      thumbnail.AudioOnly,
		DownloadFolder:      "",
		YtdlpPath:           "",
		FfmpegPath:          "",
		ConcurrentDownloads: 1,
		ConcurrentFragments: 3,
		LogPath:             "",
		ExtraYtdlpOptions:   "",
	}
}

func (s *Settings) Normalize() {
	s.YtdlpPath = strings.TrimSpace(s.YtdlpPath)
	s.FfmpegPath = strings.TrimSpace(s.FfmpegPath)
	s.LogPath = strings.TrimSpace(s.LogPath)
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

func (s *Settings) GetFormat() string {
	return s.Format
}

func (s *Settings) SetFormat(format string) {
	s.Format = format
}

func (s *Settings) GetEmbedThumbnail() string {
	return s.EmbedThumbnail
}

func (s *Settings) SetEmbedThumbnail(embedThumbnail string) {
	s.EmbedThumbnail = embedThumbnail
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

func (s *Settings) SetDownloadFolder(downloadFolder string) {
	s.DownloadFolder = downloadFolder
}

func (s *Settings) GetYtdlpPath() string {
	if s.YtdlpPath == "" {
		return "./yt-dlp.exe"
	}
	_, err := os.Stat(s.YtdlpPath)
	if err != nil {
		return "./yt-dlp.exe"
	}
	return s.YtdlpPath
}

func (s *Settings) SetYtdlpPath(ytdlpPath string) {
	s.YtdlpPath = ytdlpPath
}

func (s *Settings) GetFfmpegPath() string {
	if s.FfmpegPath == "" {
		return "./ffmpeg.exe"
	}
	_, err := os.Stat(s.FfmpegPath)
	if err != nil {
		return "./ffmpeg.exe"
	}
	return s.FfmpegPath
}

func (s *Settings) SetFfmpegPath(ffmpegPath string) {
	s.FfmpegPath = ffmpegPath
}

func (s *Settings) GetConcurrentDownloads() uint32 {
	return s.ConcurrentDownloads
}

func (s *Settings) SetConcurrentDownloads(concurrentDownloads uint32) {
	s.ConcurrentDownloads = concurrentDownloads
}

func (s *Settings) GetConcurrentFragments() uint32 {
	return s.ConcurrentFragments
}

func (s *Settings) SetConcurrentFragments(concurrentFragments uint32) {
	s.ConcurrentFragments = concurrentFragments
}

func (s *Settings) GetLogPath() string {
	if s.LogPath == "" {
		return "./log.txt"
	}
	_, err := os.Stat(s.LogPath)
	if err != nil {
		return "./log.txt"
	}
	return s.LogPath
}

func (s *Settings) SetLogPath(logPath string) {
	s.LogPath = logPath
}

func (s *Settings) GetExtraYtdlpOptions() string {
	return s.ExtraYtdlpOptions
}

func (s *Settings) ExtraYtdlpOptionsAsArray() []string {
	return strings.Fields(s.ExtraYtdlpOptions)
}

func (s *Settings) SetExtraYtdlpOptions(extraYtdlpOptions string) {
	s.ExtraYtdlpOptions = extraYtdlpOptions
}
