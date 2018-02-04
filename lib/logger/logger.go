package logger

import (
	"os"
	"github.com/sirupsen/logrus"
)

var Logger = logrus.New()

func init() {
	Logger.Formatter = &logrus.TextFormatter{FullTimestamp:true}
	Logger.Out = os.Stdout
	Logger.Level = logrus.DebugLevel
}



