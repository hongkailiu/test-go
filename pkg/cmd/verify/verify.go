package verify

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
)

type Option struct {
	CAFile string
}

func Run(o Option, certificates []string) error {
	log.WithField("CAfile", o.CAFile).
		WithField("certificates", certificates).
		Debug("verify")

	bytes, err := os.ReadFile(o.CAFile)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}
	block, _ := pem.Decode(bytes)
	caCertificate, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse the certificate: %w", err)
	}

	bytes, err = os.ReadFile(certificates[0])
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}
	block, _ = pem.Decode(bytes)
	certificate, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse the certificate: %w", err)
	}

	roots := x509.NewCertPool()
	roots.AddCert(caCertificate)

	_, err = certificate.Verify(x509.VerifyOptions{
		Roots:     roots,
		KeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	})
	if err != nil {
		return fmt.Errorf("failed to verify the certificate: %w", err)
	}
	log.Info("Certificate is verified")
	return nil
}
