package app

import (
	log "github.com/sirupsen/logrus"
	tmlog "github.com/tendermint/tendermint/libs/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	defaultLogSizeBytes = 1000
	defaultMaxBackups   = 3
	defaultMaxAgeDays   = 28
	moduleKey           = "module"
)

type Logger struct {
	*log.Logger
	Fields              log.Fields
	moduleOverrideLevel map[string]string
	logFilePath         string
}

// NewLog creates a new Log struct with the given persistent fields
func NewLogger(path string, level string, moduleOverrideLevel ...map[string]string) Logger {
	var logger = Logger{
		Fields:              map[string]interface{}{},
		moduleOverrideLevel: map[string]string{},
		logFilePath:         path,
	}
	if len(moduleOverrideLevel) > 0 {
		logger.moduleOverrideLevel = moduleOverrideLevel[0]
	}
	logger.Logger = logger.setupLogger(level)

	return logger
}

func (l Logger) setupLogger(level string) *log.Logger {
	logger := log.New()
	// Set level
	if level != "" {
		level, err := log.ParseLevel(level)
		if err != nil {
			l.Error("failed to parse log level", err)
		} else {
			logger.SetLevel(level)
		}
	}
	// Set log file path
	if l.logFilePath != "" {
		logger.SetOutput(&lumberjack.Logger{
			Filename:   l.logFilePath,
			MaxSize:    defaultLogSizeBytes, // megabytes
			MaxBackups: defaultMaxBackups,
			MaxAge:     defaultMaxAgeDays, //days
			Compress:   true,              // disabled by default
		})
	}
	logger.SetFormatter(&log.TextFormatter{})
	return logger
}

func (l Logger) Debug(msg string, keyvals ...interface{}) {
	l.Logger.WithFields(l.Fields).Debug(msg, keyvals)
}

func (l Logger) Info(msg string, keyvals ...interface{}) {
	l.Logger.WithFields(l.Fields).Info(msg, keyvals)
}
func (l Logger) Error(msg string, keyvals ...interface{}) {
	l.Logger.WithFields(l.Fields).Error(msg, keyvals)
}

func (l Logger) With(keyvals ...interface{}) tmlog.Logger {
	//Make deep copy of the current fields
	fields := map[string]interface{}{}
	for k, v := range l.Fields {
		fields[k] = v
	}

	logger := l.Logger

	for i := 0; i < len(keyvals); i += 2 {
		key, ok := keyvals[i].(string)
		if !ok {
			return l
		}
		// Check if the key is a module and if it has an override level
		if key == moduleKey {
			if val, ok := l.moduleOverrideLevel[keyvals[i+1].(string)]; ok {
				logger = l.setupLogger(val)
			}
		}
		value := keyvals[i+1]
		fields[key] = value
	}

	return Logger{
		Logger: logger,
		Fields: fields,
	}
}
