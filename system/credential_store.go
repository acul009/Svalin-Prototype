package system

import (
	"errors"
	"fmt"

	"github.com/rahn-it/svalin/db"
	"github.com/rahn-it/svalin/pki"
	"go.etcd.io/bbolt"
)

type CredentialStore struct {
	scope db.Scope
}

func OpenCredentialStore(scope db.Scope) *CredentialStore {
	return &CredentialStore{
		scope: scope,
	}
}

func (cs *CredentialStore) LoadCredentials(name string, password []byte) (*pki.PermanentCredentials, error) {

	var raw []byte
	err := cs.scope.View(func(b *bbolt.Bucket) error {
		raw = b.Get([]byte(name))
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to load credentials: %w", err)
	}

	if raw == nil {
		return nil, errors.New("credentials not found")
	}

	credentials, err := pki.CredentialsFromPem(raw, password)
	if err != nil {
		return nil, fmt.Errorf("failed to parse credentials: %w", err)
	}

	return credentials, nil
}

func (cs *CredentialStore) SaveCredentials(name string, credentials *pki.PermanentCredentials, password []byte) error {

	raw, err := credentials.PemEncode(password)
	if err != nil {
		return fmt.Errorf("failed to encode credentials: %w", err)
	}

	key := []byte(name)

	err = cs.scope.Update(func(b *bbolt.Bucket) error {
		return b.Put(key, raw)
	})

	if err != nil {
		return fmt.Errorf("failed to save credentials: %w", err)
	}

	return nil
}

func (cs *CredentialStore) List() []string {
	names := make([]string, 16)
	err := cs.scope.View(func(b *bbolt.Bucket) error {
		return b.ForEach(func(k, v []byte) error {
			names = append(names, string(k))
			return nil
		})
	})

	if err != nil {
		panic(err)
	}

	return names
}
