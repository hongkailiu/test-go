package file

import (
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

var (
	logger *logrus.Logger
)

// NewLoggerWithLFSHook returns a logrus logger which writes both on standard error and a file
// https://github.com/Sirupsen/logrus/issues/230
// if we do not need the color on terminal, the following will work as described here
// https://github.com/Sirupsen/logrus/issues/230#issuecomment-381639138
// mw := io.MultiWriter(os.Stdout, logFile)
// logrus.SetOutput(mw)
func NewLoggerWithLFSHook(logFilePath string) *logrus.Logger {
	pathMap := lfshook.PathMap{
		logrus.InfoLevel:  logFilePath,
		logrus.ErrorLevel: logFilePath,
	}

	logger = logrus.New()

	logger.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})

	logger.Hooks.Add(lfshook.NewHook(
		pathMap,
		&logrus.TextFormatter{},
	))
	return logger
}
