package pkey

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

type Option struct {
	PubOut bool
	Text   bool
	NoOut  bool
	In     string
	Out    string
}

func Run(o Option) error {
	log.WithField("pubout", o.PubOut).
		WithField("noout", o.NoOut).
		WithField("text", o.Text).
		WithField("in", o.In).
		WithField("out", o.Out).
		Debug("pkey")

	bytes, err := os.ReadFile(o.In)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}
	block, _ := pem.Decode(bytes)
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse the private key: %w", err)
	}

	k := &privateKey.PublicKey
	publicKey := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(k),
	}

	pubKeyBytes := pem.EncodeToMemory(publicKey)
	log.Debug("Public key generated")

	if err := os.WriteFile(o.Out, pubKeyBytes, 0644); err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	log.WithField("file", o.Out).Info("Saved the public key")
	return nil
}
