package logger

import (
	"encoding/json"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// mlogger is the global metrics application logger.
var mlogger *zap.Logger

// Init initializes a logger to be used globally
// throughout the metrics collector application.
func Init() {
	// create global logger instance
	rawJSON := []byte(`{
		"level": "debug",
		"encoding": "json",
		"outputPaths": ["stdout", "/tmp/lmc.log"],
		"errorOutputPaths": ["stderr", "/tmp/lmc.log"],
		"encoderConfig": {
		  "messageKey": "message",
		  "levelKey": "level",
			"timeKey": "time",
		  "levelEncoder": "lowercase"
		}
	}`)
	var cfg zap.Config
	if err := json.Unmarshal(rawJSON, &cfg); err != nil {
		panic(err)
	}

	// Add time encoder
	cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.UnixDate)

	l, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	mlogger = l
	defer mlogger.Sync()
}

// Info logs an info level message.
func Info(msg string) {
	mlogger.Info(msg)
}

// Debug logs a debug level message.
func Debug(msg string) {
	mlogger.Debug(msg)
}

// Error logs a debug level error.
func Error(err string) {
	mlogger.Error(err)
}
