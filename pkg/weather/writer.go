package weather

import (
	"encoding/json"

	"github.com/sirupsen/logrus"
)

type Writer interface {
	Write(response Response) error
}

type Logger struct {
	logger *logrus.Entry
}

func (w *Logger) Write(response Response) error {
	bytes, err := json.Marshal(response)
	if err != nil {
		return nil
	}
	w.logger.Infof("%s", string(bytes))
	return nil
}
