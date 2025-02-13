package logger

import (
	"os"
	"time"

	"github.com/sirupsen/logrus" //nolint:depguard
)

type Logger struct {
	*logrus.Logger
}

func New(level string) *Logger {
	logger := logrus.New()

	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		lvl = logrus.InfoLevel
	}
	logger.SetLevel(lvl)

	logger.SetOutput(os.Stdout)

	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.RFC3339,
		PadLevelText:    true,
	})

	logger.Hooks.Add(&timeZoneHook{})

	return &Logger{logger}
}

type timeZoneHook struct{}

func (h *timeZoneHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *timeZoneHook) Fire(entry *logrus.Entry) error {
	entry.Time = entry.Time.Local()
	return nil
}
