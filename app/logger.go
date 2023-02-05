package app

import (
	"path/filepath"

	log "github.com/sirupsen/logrus"
	tmlog "github.com/tendermint/tendermint/libs/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Logger struct {
	*log.Logger
	Fields log.Fields
}

// NewLog creates a new Log struct with the given persistent fields
func NewLogger(path string) Logger {
	var logger = Logger{
		Logger: log.New(),
		Fields: map[string]interface{}{},
	}

	if path != "" {
		logger.SetOutput(&lumberjack.Logger{
			Filename:   filepath.Join(path, "log/rollapp.log"),
			MaxSize:    1000, // megabytes
			MaxBackups: 3,
			MaxAge:     28,   //days
			Compress:   true, // disabled by default
		})
	}
	logger.SetFormatter(&log.TextFormatter{})
	logger.SetLevel(log.DebugLevel)

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
	//FIXME: allow support of multiple With assignment
	// if len(l.Fields) > 0 {
	// 	l.Error("only single With assignment supported for logging. Cant assign new keyvals", keyvals...)
	// }
	fields := log.Fields{}
	for i := 0; i < len(keyvals); i += 2 {
		key, ok := keyvals[i].(string)
		if !ok {
			return l
		}
		value := keyvals[i+1]
		fields[key] = value
	}

	return Logger{
		Logger: l.Logger,
		Fields: fields,
	}
}
