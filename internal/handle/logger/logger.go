package logger

import (
	"github.com/rs/zerolog"
	"sync"
	"ytb-downloader/internal/settings"
)

var loggerOnce sync.Once
var Queue zerolog.Logger
var Downloader zerolog.Logger

func InitLogger() {
	loggerOnce.Do(func() {
		Queue = settings.Get().GetLogger().With().Str("scope", "queue").Logger()
		Downloader = settings.Get().GetLogger().With().Str("scope", "downloader").Logger()
	})
}
