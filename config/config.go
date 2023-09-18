package config

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"time"
)

func getConfigDir() []string {
	if os.Getenv("OS") == "Windows_NT" {
		return []string{os.Getenv("APPDATA"), "rahnit-rmm"}
	}
	return []string{"/etc/rahnit-rmm"}
}

func getConfigFilePath(filePath ...string) string {
	pathParts := append(getConfigDir(), filePath...)
	return filepath.Join(pathParts...)
}

func GetCaCert() (*x509.Certificate, error) {
	// Read the CA certificate file
	caCertPEM, err := os.ReadFile(getConfigFilePath("ca.crt"))

	if err != nil {
		return nil, err
	}

	// Decode the PEM-encoded CA certificate
	block, _ := pem.Decode(caCertPEM)
	if block == nil {
		return nil, fmt.Errorf("failed to decode CA certificate PEM")
	}

	// Parse the CA certificate
	caCert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}

	return caCert, nil
}

const (
	caCertFilePath = "ca.crt"
	caKeyFilePath  = "ca.key"
	validFor       = 10 * 365 * 24 * time.Hour
)

func GenerateRootCert() error {
	// check if the CA certificate already exists
	_, err := GetCaCert()
	if err == nil {
		return fmt.Errorf("CA certificate already exists")
	}

	// Generate a new CA private key
	caPrivateKey, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		return err
	}

	// Create a self-signed CA certificate template
	caTemplate := x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: "Root CA"},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(validFor),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	// Create and save the self-signed CA certificate
	caCertDER, err := x509.CreateCertificate(rand.Reader, &caTemplate, &caTemplate, &caPrivateKey.PublicKey, caPrivateKey)
	if err != nil {
		return err
	}

	caCertFile, err := os.Create(getConfigFilePath(caCertFilePath))
	if err != nil {
		return err
	}

	defer caCertFile.Close()
	err = pem.Encode(caCertFile, &pem.Block{Type: "CERTIFICATE", Bytes: caCertDER})
	if err != nil {
		return err
	}

	// Save the CA private key to a file
	caKeyFile, err := os.Create(getConfigFilePath(caKeyFilePath))
	if err != nil {
		return err
	}
	defer caKeyFile.Close()
	caKeyBytes, err := x509.MarshalECPrivateKey(caPrivateKey)
	if err != nil {
		return err
	}
	err = pem.Encode(caKeyFile, &pem.Block{Type: "EC PRIVATE KEY", Bytes: caKeyBytes})
	if err != nil {
		return err
	}

	return nil
}
