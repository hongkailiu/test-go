package rand

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

type Option struct {
	Base64 bool
	Out    string
}

func Run(o Option, num int) error {
	log.WithField("base64", o.Base64).
		WithField("out", o.Out).
		WithField("num", num).
		Debug("rand")

	token := make([]byte, num)
	_, err := rand.Read(token)
	if err != nil {
		return fmt.Errorf("failed to read token: %w", err)
	}

	encoded := base64.StdEncoding.EncodeToString(token)
	if err := os.WriteFile(o.Out, []byte(encoded), 0644); err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	log.WithField("file", o.Out).Info("Saved the output")
	return nil
}
