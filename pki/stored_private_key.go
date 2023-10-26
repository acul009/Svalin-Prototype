package pki

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
)

type storedPrivateKey struct {
	privateKey    *PrivateKey
	password      []byte
	allowOverride bool
	path          string
}

// Custom go error to indicate that the CA certificate is missing
type privateKeyMissingError struct {
	cause error
}

func (e privateKeyMissingError) Error() string {
	return fmt.Errorf("private key not found: %w", e.cause).Error()
}

func (e privateKeyMissingError) Unwrap() error {
	return e.cause
}

var ErrPrivateKeyMissing = privateKeyMissingError{}

func (s *storedPrivateKey) Available() bool {
	_, err := os.Stat(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		log.Printf("failed to check if certificate exists: %v", err)
		return false
	}
	return true
}

func (s *storedPrivateKey) Get() (*PrivateKey, error) {
	if s.privateKey == nil {
		privateKey, err := loadPrivateKeyFromFile(s.path, s.password)
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				return nil, &certMissingError{
					cause: err,
				}
			}
			return nil, fmt.Errorf("failed to load certificate: %w", err)
		}

		s.privateKey = privateKey
	}

	return s.privateKey, nil
}

func (s *storedPrivateKey) Set(privateKey *PrivateKey) error {
	if !s.allowOverride && s.Available() {
		return errors.New("cannot override certificate")
	}

	err := s.privateKey.saveToFile(s.path, s.password)
	if err != nil {
		return fmt.Errorf("failed to save certificate: %w", err)
	}

	s.privateKey = privateKey

	return nil
}
