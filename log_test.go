package log

import (
	"fmt"
	"testing"

	"github.com/go-logr/logr"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var logrus = log.NewEntry(log.New())

func TestNewLogger(t *testing.T) {
	var defaultLogger = &logrusLogger{
		logger:       logrus,
		defaultLevel: log.InfoLevel,
		name:         "default",
	}

	newLog := newLogger(defaultLogger)
	assert.Equal(t, newLog, defaultLogger)

	// Change the newLogger and make sure the original isn't modified.
	newLog.name = "newName"
	assert.NotEqual(t, newLog, defaultLogger)
}

func TestEnabled(t *testing.T) {
	logger := InitLogger(log.NewEntry(log.New()))
	assert.True(t, logger.Enabled(3))
	assert.True(t, logger.Enabled(1))
	assert.False(t, logger.Enabled(6))

	l := log.New()
	l.SetLevel(log.TraceLevel)
	logger = InitLogger(log.NewEntry(l))
	assert.True(t, logger.Enabled(6))
}

func TestLevel(t *testing.T) {
	levelTests := []int{-1000, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 100, 1000}
	for _, level := range levelTests {
		t.Run(fmt.Sprintf("defaultLevel-%d", level), func(t *testing.T) {
			logger := InitLogger(log.NewEntry(log.New()))
			assert.NotPanics(t, func() { logger.Info(level, "") })
		})
	}
}

func TestParseFields(t *testing.T) {
	var defaultLogger = &logrusLogger{
		logger:       logrus,
		defaultLevel: log.InfoLevel,
		name:         "default",
	}

	var parseTests = []struct {
		name     string
		args     []string
		expected log.Fields
	}{
		{
			"noArgs",
			[]string{},
			make(map[string]interface{}),
		}, {
			"oneArgs",
			[]string{"key1"},
			make(map[string]interface{}),
		}, {
			"twoArgs",
			[]string{"key1", "val1"},
			map[string]interface{}{
				"key1": "val1",
			},
		}, {
			"threeArgs",
			[]string{"key1", "val1", "key2"},
			map[string]interface{}{
				"key1": "val1",
			},
		}, {
			"fourArgs",
			[]string{"key1", "val1", "key2", "val2"},
			map[string]interface{}{
				"key1": "val1",
				"key2": "val2",
			},
		},
	}
	for _, tt := range parseTests {
		t.Run(tt.name, func(t *testing.T) {
			var args []interface{}
			for _, arg := range tt.args {
				args = append(args, arg)
			}
			assert.Equal(
				t,
				parseFields(defaultLogger.logger, tt.name, args),
				tt.expected,
			)
		})
	}
}

func TestWithLogr(t *testing.T) {
	initLogger := InitLogger(log.NewEntry(log.New()))
	logger := logr.New(initLogger)
	logger.Info("logging")

	// Make sure new loggers do not panic and level propagates

	withVals := logger.WithValues("test", "test")
	withVals.Info("logging with values")
	assert.True(t, withVals.GetSink().Enabled(int(log.InfoLevel)))
	assert.False(t, withVals.GetSink().Enabled(int(log.TraceLevel)))

	withName := withVals.WithName("test-logget")
	withName.Info("logging with values and name")
	assert.True(t, withVals.GetSink().Enabled(int(log.InfoLevel)))
	assert.False(t, withVals.GetSink().Enabled(int(log.TraceLevel)))
}
