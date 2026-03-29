package logger

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
)

func NewConsoleWriter() io.Writer {
	out := zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) {
		w.TimeFormat = time.RFC3339
		w.Out = os.Stderr

	})

	return out
}

func SetConsoleMode() {
	SetLogger(zerolog.New(NewConsoleWriter()).Level(GetLevelFromEnv()).With().Timestamp().Logger())
}
