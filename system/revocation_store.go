package system

import (
	"crypto"

	"github.com/rahn-it/svalin/db"
	"github.com/rahn-it/svalin/pki"
)

type RevocationStore struct {
	scope db.Scope
	root  *pki.Certificate
}

type revocation struct {
	chain []*pki.Certificate
	hash  []byte
	date  int64
}

func OpenRevocationStore(scope db.Scope, root *pki.Certificate) (*RevocationStore, error) {
	return &RevocationStore{
		scope: scope,
		root:  root,
	}, nil
}

func (rs *RevocationStore) getHashers() map[string]crypto.Hash {
	return map[string]crypto.Hash{
		"sha512_": crypto.SHA512,
	}
}

func (rs *RevocationStore) CheckCertificate(cert *pki.Certificate) error {
	return rs.check(cert.BinaryEncode())
}

func (rs *RevocationStore) check(payload []byte) error {
	hashers := rs.getHashers()

	hashKeys := make([][]byte, 0, len(hashers))
	for hashPrefix, hashAlg := range hashers {

		hash := hashAlg.New().Sum(payload)
		hashKey := make([]byte, len(hashPrefix)+len(hash))
		copy(hashKey, hashPrefix)
		copy(hashKey[len(hashPrefix):], hash)

		hashKeys = append(hashKeys, hashKey)
	}

	return nil
	// TODO

}
