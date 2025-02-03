package logger

import (
	"fmt"
	"github.com/rs/zerolog"
	"os"
	"sync"
	"ytb-downloader/internal/settings"
)

var loggerOnce sync.Once
var Queue zerolog.Logger
var Downloader zerolog.Logger

func InitLogger() {
	loggerOnce.Do(func() {
		file, err := os.OpenFile(settings.Get().GetLogPath(), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			panic(fmt.Sprintf("Error creating log file: %v\n", err))
		}
		consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout}
		multi := zerolog.MultiLevelWriter(consoleWriter, file)
		logger := zerolog.New(multi).With().Timestamp().Logger()

		Queue = logger.With().Str("scope", "queue").Logger()
		Downloader = logger.With().Str("scope", "downloader").Logger()
	})
}
