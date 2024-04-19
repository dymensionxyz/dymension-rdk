package logger

import (
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	tmlog "github.com/tendermint/tendermint/libs/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	defaultMaxBackups = 3
	defaultMaxAgeDays = 28
	moduleKey         = "module"
)

type Logger struct {
	*log.Logger
	Fields              log.Fields
	moduleOverrideLevel map[string]string
	customLogLevel      log.Level
}

// NewLog creates a new Log struct with the given persistent fields
func NewLogger(path string, maxSize int, level string, moduleOverrideLevel ...map[string]string) Logger {
	logger := Logger{
		Fields:              map[string]interface{}{},
		moduleOverrideLevel: map[string]string{},
	}
	if len(moduleOverrideLevel) > 0 {
		logger.moduleOverrideLevel = moduleOverrideLevel[0]
	}
	logger.Logger = logger.setupLogger(path, maxSize, level)
	logger.customLogLevel = logger.GetLevel()
	return logger
}

func (l Logger) setupLogger(path string, maxSize int, level string) *log.Logger {
	logger := log.New()
	// Set log file path
	if path != "" {
		logger.SetOutput(&lumberjack.Logger{
			Filename:   path,
			MaxSize:    maxSize, // megabytes
			MaxBackups: defaultMaxBackups,
			MaxAge:     defaultMaxAgeDays, // days
			Compress:   true,              // disabled by default
		})
	} else {
		logger.SetOutput(os.Stdout)
	}
	logLevel, err := log.ParseLevel(level)
	if err != nil {
		logger.Error("failed to parse log level", "error", err, "level", level)
	} else {
		logger.SetLevel(logLevel)
	}
	logger.SetFormatter(&log.TextFormatter{})
	return logger
}

func (l Logger) Debug(msg string, keyvals ...interface{}) {
	if l.customLogLevel < log.DebugLevel {
		return
	}
	l.Logger.WithFields(l.Fields).Debug(msg, keyvals)
}

func (l Logger) Info(msg string, keyvals ...interface{}) {
	if l.customLogLevel < log.InfoLevel {
		return
	}
	l.Logger.WithFields(l.Fields).Info(msg, keyvals)
}

func (l Logger) Error(msg string, keyvals ...interface{}) {
	if l.customLogLevel < log.ErrorLevel {
		return
	}
	l.Logger.WithFields(l.Fields).Error(msg, keyvals)
}

func (l Logger) With(keyvals ...interface{}) tmlog.Logger {
	// Make deep copy of the current fields
	fields := map[string]interface{}{}
	for k, v := range l.Fields {
		fields[k] = v
	}

	logger := l.Logger
	customLogLevel := l.customLogLevel

	for i := 0; i < len(keyvals); i += 2 {
		key, ok := keyvals[i].(string)
		if !ok {
			return l
		}
		// Check if the key is a module and if it has an override level
		if key == moduleKey {
			if val, ok := l.moduleOverrideLevel[keyvals[i+1].(string)]; ok {
				newLogLevel, err := log.ParseLevel(val)
				if err != nil {
					logger.Error("failed to parse log level", "error", err, "level", val)
				}
				if newLogLevel > customLogLevel {
					logger.Error("can't increase log level for a module")
				} else {
					customLogLevel = newLogLevel
				}
			}
		}
		value := keyvals[i+1]
		fields[key] = value
	}

	return Logger{
		Logger:         logger,
		Fields:         fields,
		customLogLevel: customLogLevel,
	}
}

// ConvertStringToStringMap converts a string to a map[string]string.
// The input string is expected to be for example "key1:value1,key2:value2" where seperators can be specified.
func ConvertStringToStringMap(input string, mapSeperator string, kvSeperator string) map[string]string {
	resultMap := make(map[string]string)
	if input != "" {
		for _, keyValue := range strings.Split(input, mapSeperator) {
			kv := strings.Split(keyValue, kvSeperator)
			resultMap[kv[0]] = kv[1]
		}
	}
	return resultMap
}
