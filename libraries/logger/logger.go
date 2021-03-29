package logger

import (
	"os"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func syslogTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

// GetLogger returns a *zap.SugaredLogger
func GetLogger(module string) *zap.SugaredLogger {
	logLevel := os.Getenv("LOG_LEVEL")
	upperModule := strings.ToUpper(module)
	if os.Getenv("LOG_LEVEL_"+upperModule) != "" {
		logLevel = os.Getenv("LOG_LEVEL_" + upperModule)
	}

	runEnv := os.Getenv("RUN_ENV")
	var config zap.Config
	if strings.ToUpper(runEnv) == "DEV" {
		config = zap.NewDevelopmentConfig()
	} else {
		config = zap.NewProductionConfig()
	}

	config.Level.UnmarshalText([]byte(logLevel))
	config.EncoderConfig.EncodeTime = syslogTimeEncoder
	log, _ := config.Build()

	return log.Named(module).Sugar()
}
