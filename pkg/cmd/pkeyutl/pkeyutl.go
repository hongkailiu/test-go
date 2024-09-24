package pkeyutl

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

type Option struct {
	Encrypt bool
	Decrypt bool
	In      string
	PubIn   bool
	InKey   string
	Out     string
}

func Run(o Option) error {
	log.WithField("encrypt", o.Encrypt).
		WithField("decrypt", o.Decrypt).
		WithField("in", o.In).
		WithField("pubin", o.PubIn).
		WithField("inkey", o.InKey).
		WithField("out", o.Out).
		Debug("pkeyutl")

	inBytes, err := os.ReadFile(o.In)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	if o.Encrypt {
		bytes, err := os.ReadFile(o.InKey)
		if err != nil {
			return fmt.Errorf("failed to read file: %w", err)
		}
		block, _ := pem.Decode(bytes)
		publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return fmt.Errorf("failed to parse the public key: %w", err)
		}

		ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey.(*rsa.PublicKey), inBytes)
		if err != nil {
			return fmt.Errorf("failed to encrypt: %w", err)
		}

		if err := os.WriteFile(o.Out, ciphertext, 0644); err != nil {
			return fmt.Errorf("failed to write to file: %w", err)
		}
		log.WithField("out", o.Out).Info("Saved output")
	}

	if o.Decrypt {
		bytes, err := os.ReadFile(o.InKey)
		if err != nil {
			return fmt.Errorf("failed to read file: %w", err)
		}
		block, _ := pem.Decode(bytes)
		privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return fmt.Errorf("failed to parse the private key: %w", err)
		}

		plaintext, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey.(*rsa.PrivateKey), inBytes)
		if err != nil {
			return fmt.Errorf("failed to decrypt: %w", err)
		}

		if err := os.WriteFile(o.Out, plaintext, 0644); err != nil {
			return fmt.Errorf("failed to write to file: %w", err)
		}
		log.WithField("out", o.Out).Info("Saved output")
	}
	return nil
}
