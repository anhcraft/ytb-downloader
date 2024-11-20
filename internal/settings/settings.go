package settings

import (
	"fmt"
	"github.com/rs/zerolog"
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
	LogPath             string `json:"logPath,omitempty"`
	ExtraYtdlpOptions   string `json:"extraYtdlpOptions,omitempty"`
	globalLogger        *zerolog.Logger
}

func (s *Settings) GetLogger() *zerolog.Logger {
	if s.globalLogger == nil {
		file, err := os.OpenFile(s.GetLogPath(), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			panic(fmt.Sprintf("Error creating log file: %v\n", err))
		}
		consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout}
		multi := zerolog.MultiLevelWriter(consoleWriter, file)
		logger := zerolog.New(multi).With().Timestamp().Logger()
		s.globalLogger = &logger
	}
	return s.globalLogger
}

func (s *Settings) Normalize() {
	s.YTdlpPath = strings.TrimSpace(s.YTdlpPath)
	s.FFmpegPath = strings.TrimSpace(s.FFmpegPath)
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

func (s *Settings) GetExtraYtdlpOptions() []string {
	return strings.Split(s.ExtraYtdlpOptions, " ")
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
		LogPath:             "",
		ExtraYtdlpOptions:   "",
	}
}
