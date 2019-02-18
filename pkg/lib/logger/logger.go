package logger

/*import (
	"os"

	"github.com/sirupsen/logrus"
	lSyslog "github.com/sirupsen/logrus/hooks/syslog"
	"log/syslog"
)

// Logger with syslog hook
var Logger = logrus.New()

func init() {
	Logger.Formatter = &logrus.TextFormatter{FullTimestamp: true}
	Logger.Out = os.Stdout
	Logger.Level = logrus.DebugLevel

	hook, err := lSyslog.NewSyslogHook("", "", syslog.LOG_INFO, "")

	if err != nil {
		panic(err.Error())
	}

	Logger.Hooks.Add(hook)
}*/

import (
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

var (
	logger *logrus.Logger
)

func NewLogger(logFilePath string) *logrus.Logger {
	if logger != nil {
		return logger
	}

	pathMap := lfshook.PathMap{
		logrus.InfoLevel:  logFilePath,
		logrus.ErrorLevel: logFilePath,
	}

	logger = logrus.New()

	logger.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})

	logger.Hooks.Add(lfshook.NewHook(
		pathMap,
		//&logrus.JSONFormatter{},
		&logrus.TextFormatter{},
	))
	return logger
}
