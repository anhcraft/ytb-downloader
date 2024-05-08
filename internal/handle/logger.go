package handle

import (
	"github.com/rs/zerolog"
	"sync"
	"ytb-downloader/internal/settings"
)

var loggerOnce sync.Once
var queueLogger zerolog.Logger
var downloaderLogger zerolog.Logger

func InitLogger() {
	loggerOnce.Do(func() {
		queueLogger = settings.Get().GetLogger().With().Str("scope", "queue").Logger()
		downloaderLogger = settings.Get().GetLogger().With().Str("scope", "downloader").Logger()
	})
}
