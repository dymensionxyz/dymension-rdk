package app_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dymensionxyz/rollapp/app"
	"github.com/dymensionxyz/rollapp/utils"
)

func TestLogLevel(t *testing.T) {
	var buf bytes.Buffer

	logger := app.NewLogger("", "error", nil)
	logger.SetOutput(&buf)

	logger.Debug("debug msg")
	logger.Info("info msg")
	logger.Error("error msg")
	msg := strings.TrimSpace(buf.String())

	assert.Contains(t, msg, "error msg", "error msg wasn't logged")
	assert.NotContains(t, msg, "info")
	assert.NotContains(t, msg, "debug")
}

func TestMultipleLoggerWithMethod(t *testing.T) {
	var buf bytes.Buffer

	logger := app.NewLogger("", "info", nil)
	logger.SetOutput(&buf)

	logger1 := logger.With("module", "logger1")
	logger2 := logger.With("module", "logger2")

	logger1.Info("testing")
	msg := strings.TrimSpace(buf.String())
	assert.Contains(t, msg, "logger1", "logger1 didn't have the module name as expected")
	assert.NotContains(t, msg, "logger2", "logger1 shouldn't have logger2 in it's log")

	buf.Reset()
	logger2.Info("testing")
	msg = strings.TrimSpace(buf.String())
	assert.Contains(t, msg, "logger2", "logger2 didn't have the module name as expected")
	assert.NotContains(t, msg, "logger1", "logger2 shouldn't have logger1 in it's log")
}

func TestMultipleWithCalls(t *testing.T) {
	// t.Skip("nested With calls not supported")
	var buf bytes.Buffer

	logger := app.NewLogger("", "info", nil)
	logger.SetOutput(&buf)

	logger1 := logger.With("module", "module1")
	logger2 := logger1.With("arg", "custom_arg") //uses logger1 as context

	logger1.Info("testing")
	msg := strings.TrimSpace(buf.String())
	assert.Contains(t, msg, "module1", "logger1 didn't have the module name as expected")
	assert.NotContains(t, msg, "custom_arg", "logger1 shouldn't have logger2 in it's log")

	buf.Reset()
	logger2.Info("testing")
	msg = strings.TrimSpace(buf.String())
	assert.Contains(t, msg, "module1", "logger2 didn't have the module name as expected")
	assert.Contains(t, msg, "custom_arg", "logger2 didn't have the module name as expected")
}

func TestModuleOverrideLevel(t *testing.T) {
	var buf bytes.Buffer

	moduleOverrides := utils.ConvertStringToStringMap("module2:info", ",", ":")
	logger := app.NewLogger("", "error", moduleOverrides)
	logger.SetOutput(&buf)

	logger1 := logger.With("module", "module1")

	logger2 := logger.With("module", "module2").(app.Logger)
	logger2.SetOutput(&buf)

	/* -------------------------- check module 1 logger ------------------------- */
	logger1.Debug("testing debug")
	logger1.Info("testing info")
	logger1.Error("testing error")
	msg := strings.TrimSpace(buf.String())
	assert.Contains(t, msg, "testing error", "logger didn't have the module name as expected")
	assert.NotContains(t, msg, "testing info", "logger didn't have the module name as expected")
	assert.NotContains(t, msg, "testing debug", "logger didn't have the module name as expected")

	buf.Reset()

	/* -------------------------- check module 2 logger ------------------------- */

	logger2.Debug("testing debug")
	logger2.Info("testing info")
	logger2.Info("testing error")
	msg = strings.TrimSpace(buf.String())
	assert.Contains(t, msg, "testing error", "logger didn't have the module name as expected")
	assert.Contains(t, msg, "testing info", "logger didn't have the module name as expected")
	assert.NotContains(t, msg, "testing debug", "logger didn't have the module name as expected")
}

// // NewLog creates a new Log struct with the given persistent fields
// func NewLogger(path string, level string, moduleOverrideLevel ...map[string]string) Logger {
// 	var logger = Logger{
// 		Fields:              map[string]interface{}{},
// 		moduleOverrideLevel: map[string]string{},
// 		logFilePath:         path,
// 	}
// 	if len(moduleOverrideLevel) > 0 {
// 		logger.moduleOverrideLevel = moduleOverrideLevel[0]
// 	}
// 	logger.Logger = logger.setupLogger(level)

// 	return logger
// }

// func (l Logger) setupLogger(level string) *log.Logger {
// 	logger := log.New()
// 	// Set level
// 	if level != "" {
// 		level, err := log.ParseLevel(level)
// 		if err != nil {
// 			l.Error("failed to parse log level", err)
// 		} else {
// 			logger.SetLevel(level)
// 		}
// 	}
// 	// Set log file path
// 	if l.logFilePath != "" {
// 		logger.SetOutput(&lumberjack.Logger{
// 			Filename:   l.logFilePath,
// 			MaxSize:    defaultLogSizeBytes, // megabytes
// 			MaxBackups: defaultMaxBackups,
// 			MaxAge:     defaultMaxAgeDays, //days
// 			Compress:   true,              // disabled by default
// 		})
// 	}
// 	logger.SetFormatter(&log.TextFormatter{})
// 	return logger
// }

// func (l Logger) Debug(msg string, keyvals ...interface{}) {
// 	l.Logger.WithFields(l.Fields).Debug(msg, keyvals)
// }

// func (l Logger) Info(msg string, keyvals ...interface{}) {
// 	l.Logger.WithFields(l.Fields).Info(msg, keyvals)
// }
// func (l Logger) Error(msg string, keyvals ...interface{}) {
// 	l.Logger.WithFields(l.Fields).Error(msg, keyvals)
// }

// func (l Logger) With(keyvals ...interface{}) tmlog.Logger {
// 	//FIXME: allow support of multiple With assignment
// 	// if len(l.Fields) > 0 {
// 	// 	l.Error("only single With assignment supported for logging. Cant assign new keyvals", keyvals...)
// 	// }
// 	fields := log.Fields{}
// 	logger := l.Logger

// 	for i := 0; i < len(keyvals); i += 2 {
// 		key, ok := keyvals[i].(string)
// 		if !ok {
// 			return l
// 		}
// 		// Check if the key is a module and if it has an override level
// 		if key == moduleKey {
// 			if val, ok := l.moduleOverrideLevel[keyvals[i+1].(string)]; ok {
// 				logger = l.setupLogger(val)
// 			}
// 		}
// 		value := keyvals[i+1]
// 		fields[key] = value
// 	}

// 	return Logger{
// 		Logger: logger,
// 		Fields: fields,
// 	}
// }
