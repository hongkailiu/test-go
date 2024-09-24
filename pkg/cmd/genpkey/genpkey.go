package genpkey

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

type Option struct {
	Algorithm string
	PKeyOpt   []string
	Out       string
}

func Run(o Option) error {
	log.WithField("algorithm", o.Algorithm).
		WithField("pKeyOpt", o.PKeyOpt).
		WithField("out", o.Out).
		Debug("genpkey")
	bitSize := 1024
	for _, pKeyOpt := range o.PKeyOpt {
		if strings.HasPrefix(pKeyOpt, "rsa_keygen_bits:") {
			size, err := strconv.Atoi(strings.Replace(pKeyOpt, "rsa_keygen_bits:", "", 1))
			if err != nil {
				return fmt.Errorf("failed to parse the option: %w", err)
			}
			bitSize = size
			log.WithField("bitSize", bitSize).Debug("Found bitSize from option")
		}
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		return fmt.Errorf("failed to generate key: %w", err)
	}

	log.Debug("Private key generated")

	// Get ASN.1 DER format
	privDER := x509.MarshalPKCS1PrivateKey(privateKey)

	// pem.Block
	privBlock := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   privDER,
	}

	// Private key in PEM format
	privateKeyBytes := pem.EncodeToMemory(&privBlock)

	if err := os.WriteFile(o.Out, privateKeyBytes, 0600); err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	log.WithField("file", o.Out).Info("Saved the private key")
	return nil
}
