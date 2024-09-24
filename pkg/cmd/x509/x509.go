package x509

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

type Option struct {
	Req            bool
	Sha256         bool
	Days           int
	CA             string
	CAKey          string
	CACreateSerial bool
	In             string
	Out            string
	Text           bool
	NoOut          bool
	PubKey         bool
}

func Run(o Option) error {
	log.WithField("req", o.Req).
		WithField("sha256", o.Sha256).
		WithField("noout", o.NoOut).
		WithField("text", o.Text).
		WithField("ca", o.CA).
		WithField("cakey", o.CAKey).
		WithField("cacreateserial", o.CACreateSerial).
		WithField("days", o.Days).
		WithField("in", o.In).
		WithField("out", o.Out).
		WithField("pubkey", o.PubKey).
		Debug("x509")

	// TODO Aside: viewing the certificate as text
	// openssl x509 -in Alice.crt -text -noout

	if o.PubKey {
		bytes, err := os.ReadFile(o.In)
		if err != nil {
			return fmt.Errorf("failed to read file: %w", err)
		}
		block, _ := pem.Decode(bytes)
		certificate, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return fmt.Errorf("failed to parse the certificate: %w", err)
		}

		k := certificate.PublicKey.(*rsa.PublicKey)
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

	bytes, err := os.ReadFile(o.CAKey)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}
	block, _ := pem.Decode(bytes)
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	privateKey := key.(*rsa.PrivateKey)
	if err != nil {
		return fmt.Errorf("failed to parse the private key: %w", err)
	}

	bytes, err = os.ReadFile(o.CA)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}
	block, _ = pem.Decode(bytes)
	caCertificate, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse the certificate: %w", err)
	}

	bytes, err = os.ReadFile(o.In)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}
	block, _ = pem.Decode(bytes)
	certificateRequest, err := x509.ParseCertificateRequest(block.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse the certificate request: %w", err)
	}

	subj := pkix.Name{
		CommonName:         certificateRequest.Subject.CommonName,
		Country:            certificateRequest.Subject.Country,
		Province:           certificateRequest.Subject.Province,
		Locality:           certificateRequest.Subject.Locality,
		Organization:       certificateRequest.Subject.Organization,
		OrganizationalUnit: certificateRequest.Subject.OrganizationalUnit,
		ExtraNames:         certificateRequest.Subject.ExtraNames,
	}
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return fmt.Errorf("failed to generate serial number: %w", err)
	}

	tmpl := &x509.Certificate{
		SerialNumber:          serialNumber,
		Subject:               subj,
		SignatureAlgorithm:    x509.SHA256WithRSA,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Duration(o.Days) * time.Hour * 24),
		BasicConstraintsValid: true,
	}

	certificateBytes, err := x509.CreateCertificate(rand.Reader, tmpl, caCertificate, &privateKey.PublicKey, privateKey)
	if err != nil {
		return fmt.Errorf("failed to create certificate: %w", err)
	}

	b := pem.Block{Type: "CERTIFICATE", Bytes: certificateBytes}
	certBytes := pem.EncodeToMemory(&b)
	log.Debug("Certificate generated")

	if err := os.WriteFile(o.Out, certBytes, 0644); err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	log.WithField("file", o.Out).Info("Saved the certificate")
	return nil
}
