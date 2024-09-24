package req

import (
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

type Option struct {
	New    bool
	X509   bool
	NoEnc  bool
	Sha256 bool
	Days   int
	Key    string
	Out    string
}

// https://oidref.com/1.2.840.113549.1.9.1
var oidEmailAddress = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 1}

func Run(o Option) error {
	log.WithField("new", o.New).
		WithField("noout", o.Key).
		WithField("sha256", o.Sha256).
		WithField("x509", o.X509).
		WithField("noenc", o.NoEnc).
		WithField("days", o.Days).
		WithField("out", o.Out).
		Debug("pkey")

	bytes, err := os.ReadFile(o.Key)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}
	block, _ := pem.Decode(bytes)
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse the private key: %w", err)
	}

	subject := readSubject(o.X509)

	subj := pkix.Name{
		CommonName:         subject.commonName,
		Country:            []string{subject.country},
		Province:           []string{subject.province},
		Locality:           []string{subject.locality},
		Organization:       []string{subject.organization},
		OrganizationalUnit: []string{subject.organizationalUnit},
		ExtraNames: []pkix.AttributeTypeAndValue{
			{
				Type: oidEmailAddress,
				Value: asn1.RawValue{
					Tag:   asn1.TagIA5String,
					Bytes: []byte(subject.email),
				},
			},
		},
	}

	if o.X509 {
		// generate a random serial number (a real cert authority would have some logic behind this)
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
			IsCA:                  true,
			KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature,
		}

		certificateBytes, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &privateKey.PublicKey, privateKey)
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
	} else {

		template := x509.CertificateRequest{
			Subject:            subj,
			SignatureAlgorithm: x509.SHA256WithRSA,
		}

		csrBytes, err := x509.CreateCertificateRequest(rand.Reader, &template, privateKey)
		if err != nil {
			return fmt.Errorf("failed to create certificate request: %w", err)
		}
		certificateReqestbytes := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE REQUEST", Bytes: csrBytes})

		log.Debug("Certificate request generated")

		if err := os.WriteFile(o.Out, certificateReqestbytes, 0644); err != nil {
			return fmt.Errorf("failed to write to file: %w", err)
		}

		log.WithField("file", o.Out).Info("Saved the certificate request")
		return nil
	}
}

type subject struct {
	commonName         string
	country            string
	province           string
	locality           string
	organization       string
	organizationalUnit string
	email              string
}

func readSubject(x509 bool) subject {
	// TODO read from stin
	return subject{
		commonName:   "www.hliu.ca",
		country:      "CA",
		province:     "Quebec",
		locality:     "Montreal",
		organization: "hliu.ca",
	}
}
