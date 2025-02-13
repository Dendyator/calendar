package logger

import (
	"bytes"
	"testing"

	"github.com/sirupsen/logrus"         //nolint:depguard
	"github.com/stretchr/testify/assert" //nolint:depguard
)

func TestLogger_New(t *testing.T) {
	testCases := []struct {
		levelString string

		expectedLevel logrus.Level
	}{
		{"panic", logrus.PanicLevel},
		{"fatal", logrus.FatalLevel},
		{"error", logrus.ErrorLevel},
		{"warn", logrus.WarnLevel},
		{"info", logrus.InfoLevel},
		{"debug", logrus.DebugLevel},
		{"trace", logrus.TraceLevel},
		{"invalid", logrus.InfoLevel},
	}

	for _, tc := range testCases {
		t.Run(tc.levelString, func(t *testing.T) {
			logger := New(tc.levelString)
			assert.Equal(t, tc.expectedLevel, logger.GetLevel())
		})
	}
}

func TestLogger_LogMessage(t *testing.T) {
	logger := New("info")
	var buf bytes.Buffer
	logger.SetOutput(&buf)
	logger.Info("test message")
	assert.Contains(t, buf.String(), "test message")
	assert.Contains(t, buf.String(), "level=info")
}
