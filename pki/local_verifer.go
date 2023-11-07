package pki

import (
	"context"
	"crypto/x509"
	"fmt"
	"rahnit-rmm/config"
	"rahnit-rmm/ent"
	"rahnit-rmm/ent/device"
	"rahnit-rmm/ent/user"
)

type Verifier interface {
	Verify(cert *Certificate) ([]*Certificate, error)
	VerifyPublicKey(pub *PublicKey) ([]*Certificate, error)
}

type localVerify struct {
	rootPool      *x509.CertPool
	intermediates *x509.CertPool
}

func NewLocalVerify() (*localVerify, error) {
	rootCert, err := Root.Get()
	if err != nil {
		return nil, fmt.Errorf("failed to load root certificate: %w", err)
	}

	rootPool := x509.NewCertPool()
	rootPool.AddCert(rootCert.ToX509())

	intermediatePool := x509.NewCertPool()

	db := config.DB()
	users, err := db.User.Query().All(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}

	for _, user := range users {
		userCert, err := CertificateFromPem([]byte(user.Certificate))
		if err != nil {
			return nil, fmt.Errorf("failed to load user certificate: %w", err)
		}

		intermediatePool.AddCert(userCert.ToX509())
	}

	return &localVerify{
		rootPool:      rootPool,
		intermediates: intermediatePool,
	}, nil
}

func (v *localVerify) options() x509.VerifyOptions {

	return x509.VerifyOptions{
		Roots:         v.rootPool,
		Intermediates: v.intermediates,
	}
}

func (v *localVerify) Verify(cert *Certificate) ([]*Certificate, error) {
	if cert == nil {
		return nil, fmt.Errorf("certificate is nil")
	}

	chains, err := cert.ToX509().Verify(v.options())
	if err != nil || len(chains) == 0 {
		return nil, fmt.Errorf("failed to verify certificate: %w", err)
	}

	chain := make([]*Certificate, 0, len(chains[0]))

	for _, cert := range chains[0] {
		workingCert, err := ImportCertificate(cert)
		if err != nil {
			return nil, fmt.Errorf("failed to import certificate: %w", err)
		}

		err = v.checkCertificateInfo(workingCert)
		if err != nil {
			return nil, fmt.Errorf("failed to check certificate info: %w", err)
		}

		chain = append(chain, workingCert)
	}

	return chain, nil
}

func (v *localVerify) checkCertificateInfo(cert *Certificate) error {
	err := v.checkRevoked(cert)
	if err != nil {
		return fmt.Errorf("certificate has been revoked: %w", err)
	}

	return nil
}

func (v *localVerify) checkRevoked(cert *Certificate) error {
	// TODO: check if certificate was revoked
	return nil
}

func (v *localVerify) VerifyPublicKey(pub *PublicKey) ([]*Certificate, error) {
	root, err := Root.Get()
	if err != nil {
		return nil, fmt.Errorf("failed to check if public key is root: %w", err)
	}

	if root.GetPublicKey().Equal(pub) {
		return []*Certificate{root}, nil
	}

	if Upstream.Available() {
		upstream, err := Upstream.Get()
		if err != nil {
			return nil, fmt.Errorf("failed to check if public key is upstream: %w", err)
		}
		return v.Verify(upstream)
	}

	db := config.DB()

	// check for user with public key
	user, err := db.User.Query().Where(user.PublicKeyEQ(pub.Base64Encode())).Only(context.Background())
	if err != nil {
		if !ent.IsNotFound(err) {
			return nil, fmt.Errorf("failed to query user: %w", err)
		}
	} else {
		cert, err := CertificateFromPem([]byte(user.Certificate))
		if err != nil {
			return nil, fmt.Errorf("failed to load user certificate: %w", err)
		}
		return v.Verify(cert)
	}

	// check for device with public key
	device, err := db.Device.Query().Where(device.PublicKeyEQ(pub.Base64Encode())).Only(context.Background())
	if err != nil {
		if !ent.IsNotFound(err) {
			return nil, fmt.Errorf("failed to query device: %w", err)
		}
	} else {
		cert, err := CertificateFromPem([]byte(device.Certificate))
		if err != nil {
			return nil, fmt.Errorf("failed to load device certificate: %w", err)
		}
		return v.Verify(cert)
	}

	return nil, fmt.Errorf("unknown public key")
}