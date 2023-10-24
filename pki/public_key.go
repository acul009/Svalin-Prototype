package pki

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
)

type PublicKey ecdsa.PublicKey

func (pub *PublicKey) MarshalJSON() ([]byte, error) {
	bytes, err := pub.BinaryEncode()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal public key: %w", err)
	}
	return json.Marshal(bytes)
}

func (pub *PublicKey) UnmarshalJSON(data []byte) error {
	pubBytes := make([]byte, 0, len(data))
	err := json.Unmarshal(data, &pubBytes)
	if err != nil {
		return fmt.Errorf("failed to unmarshal certificate: %w", err)
	}

	newPub, err := PublicKeyFromBinary(pubBytes)
	if err != nil {
		return fmt.Errorf("failed to decode certificate: %w", err)
	}

	*pub = *newPub

	return nil
}

func (pub *PublicKey) BinaryEncode() ([]byte, error) {
	bytes, err := x509.MarshalPKIXPublicKey(pub.ToEcdsa())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal public key: %w", err)
	}

	return bytes, nil
}

func PublicKeyFromBinary(bytes []byte) (*PublicKey, error) {
	pub, err := x509.ParsePKIXPublicKey(bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	return ImportPublicKey(pub)
}

func (pub *PublicKey) PemEncode() ([]byte, error) {
	bytes, err := pub.BinaryEncode()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal public key: %w", err)
	}
	return pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: bytes}), nil
}

func PublicKeyFromPem(certPEM []byte) (*PublicKey, error) {
	block, _ := pem.Decode(certPEM)
	if block == nil {
		return nil, fmt.Errorf("failed to decode certificate PEM")
	}

	return PublicKeyFromBinary(block.Bytes)
}

func (pub *PublicKey) Base64Encode() (string, error) {
	bytes, err := pub.BinaryEncode()
	if err != nil {
		return "", fmt.Errorf("failed to marshal public key: %w", err)
	}
	return base64.StdEncoding.EncodeToString(bytes), nil
}

func (pub *PublicKey) ToEcdsa() *ecdsa.PublicKey {
	ec := ecdsa.PublicKey(*pub)
	return &ec
}

func ImportPublicKey(pub any) (*PublicKey, error) {
	switch typed := pub.(type) {
	case *ecdsa.PublicKey:
		pubRef := PublicKey(*typed)
		return &pubRef, nil

	case ecdsa.PublicKey:
		pubRef := PublicKey(typed)
		return &pubRef, nil

	default:
		return nil, fmt.Errorf("public key is not of type *ecdsa.PublicKey")
	}
}

func (pub *PublicKey) Equal(compare *PublicKey) bool {
	return pub.ToEcdsa().Equal(compare.ToEcdsa())
}
