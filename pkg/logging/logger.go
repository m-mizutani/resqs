package logging

import (
	"strings"

	"github.com/sirupsen/logrus"
)

// Logger is common logger interface of resqs
var Logger *logrus.Logger

func init() {
	Logger = logrus.New()
}

func SetLogLevel(logLevel string) {
	switch strings.ToUpper(logLevel) {
	case "DEBUG":
		Logger.SetLevel(logrus.DebugLevel)
	case "INFO":
		Logger.SetLevel(logrus.InfoLevel)
	case "WARN":
		Logger.SetLevel(logrus.WarnLevel)
	case "ERROR":
		Logger.SetLevel(logrus.ErrorLevel)

	default:
		Logger.WithField("level", logLevel).Warn("Invalid LogLevel, set INFO")
	}
}
