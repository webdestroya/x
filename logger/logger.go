package logger

import (
	"context"
	"log" //nolint:depguard

	"os"
	"time"

	"github.com/rs/zerolog"
)

var (
	Logger = zerolog.New(os.Stderr).With().Timestamp().Logger()
)

// https://docs.aws.amazon.com/lambda/latest/dg/monitoring-cloudwatchlogs-advanced.html
// AWS_LAMBDA_LOG_FORMAT

func Disable() {
	SetLogger(zerolog.Nop())
}

func init() {
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.TimestampFunc = timestampFunc
	SetLevelFromEnv()
}

func SetLogger(lgr zerolog.Logger) {
	Logger = lgr
	zerolog.DefaultContextLogger = &lgr
	log.SetFlags(0)
	log.SetOutput(Logger)
}

func SetLevelFromEnv() {
	level := zerolog.InfoLevel
	if val, ok := os.LookupEnv(`AWS_LAMBDA_LOG_LEVEL`); ok { //nolint:forbidigo
		if lvl, err := zerolog.ParseLevel(val); err == nil {
			level = lvl
		}
	} else if val, ok := os.LookupEnv(`LOG_LEVEL`); ok { //nolint:forbidigo
		if lvl, err := zerolog.ParseLevel(val); err == nil {
			level = lvl
		}
	}

	SetLogger(Logger.Level(level))
}

func Trace() *zerolog.Event {
	return Logger.Trace()
}

func Debug() *zerolog.Event {
	return Logger.Debug()
}

func Info() *zerolog.Event {
	return Logger.Info()
}

func Warn() *zerolog.Event {
	return Logger.Warn()
}

func Error() *zerolog.Event {
	return Logger.Error()
}

func Fatal() *zerolog.Event {
	return Logger.Fatal()
}

func With() zerolog.Context {
	return Logger.With()
}

func timestampFunc() time.Time {
	return time.Now().UTC()
}

func WithContext(ctx context.Context) context.Context {
	return Logger.WithContext(ctx)
}

func Ctx(ctx context.Context) *zerolog.Logger {
	return zerolog.Ctx(ctx)
}
