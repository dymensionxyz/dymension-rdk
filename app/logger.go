package app

import (
	"gopkg.in/natefinch/lumberjack.v2"
)

var RotatingLogger = lumberjack.Logger{
	Filename:   "",
	MaxSize:    1000, // megabytes
	MaxBackups: 3,
	MaxAge:     28,   //days
	Compress:   true, // disabled by default
}
