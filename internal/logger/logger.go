package logger

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.Logger
	debug bool
}

func New(logDir string, debug bool) (*Logger, error) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		log.Fatalf("failed to create log directory: %v", err)
	}

	currentDate := time.Now().Format("2006-01-02")
	logFile := filepath.Join(logDir, currentDate+".log")

	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	var core zapcore.Core
	if debug {
		consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
		fileEncoder := zapcore.NewJSONEncoder(encoderConfig)

		consoleCore := zapcore.NewCore(
			consoleEncoder,
			zapcore.AddSync(os.Stdout),
			zap.DebugLevel,
		)
		fileCore := zapcore.NewCore(
			fileEncoder,
			zapcore.AddSync(file),
			zap.DebugLevel,
		)
		core = zapcore.NewTee(consoleCore, fileCore)
	} else {
		fileEncoder := zapcore.NewJSONEncoder(encoderConfig)
		core = zapcore.NewCore(
			fileEncoder,
			zapcore.AddSync(file),
			zap.InfoLevel,
		)
	}

	zapLogger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))
	return &Logger{
		Logger: zapLogger,
		debug:  debug,
	}, nil
}

func (l *Logger) DebugLog(message string, content ...string) {
	if l.debug {
		if len(content) > 0 {
			l.Debug(message,
				zap.String("content", content[0]),
			)
		} else {
			l.Debug(message)
		}
	}
}
