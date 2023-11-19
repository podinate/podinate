package errorhandler

import (
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
