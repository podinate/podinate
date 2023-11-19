package errorhandler

import (
	"context"

	zap "go.uber.org/zap"
)

// Log is the global logger, based on zap sugared logger
var Log *zap.SugaredLogger

func Init() {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync() // flushes buffer, if any
	Log = logger.Sugar()

	Log.Infow("Initialized errorhandler", "config", "development")
}

// Error is an error-level error log
func Error(ctx context.Context, message string, args ...interface{}) {
	args = append(args, "request_id", GetRequestID(ctx))
	Log.Errorw(message, args...)
}

// Debug is a debug-level error log
func Debug(ctx context.Context, message string, args ...interface{}) {
	args = append(args, "request_id", GetRequestID(ctx))
	Log.Debugw(message, args...)
}

// Info is an info-level error log
func Info(ctx context.Context, message string, args ...interface{}) {
	args = append(args, "request_id", GetRequestID(ctx))
	Log.Infow(message, args...)
}

// Warn is a warn-level error log
func Warn(ctx context.Context, message string, args ...interface{}) {
	args = append(args, "request_id", GetRequestID(ctx))
	Log.Warnw(message, args...)
}

// Fatal is a fatal-level error log
func Fatal(ctx context.Context, message string, args ...interface{}) {
	args = append(args, "request_id", GetRequestID(ctx))
	Log.Fatalw(message, args...)
}

// Panic is a panic-level error log
func Panic(ctx context.Context, message string, args ...interface{}) {
	args = append(args, "request_id", GetRequestID(ctx))
	Log.Panicw(message, args...)
}
