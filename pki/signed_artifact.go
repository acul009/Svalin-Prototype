package pki

import (
	"encoding/json"
	"fmt"
	"rahnit-rmm/util"
)

type sArtifactPayload interface {
	MayPublish(*Certificate) bool
	Revokeable() bool
}

type SignedArtifact[T sArtifactPayload] struct {
	creator *PublicKey
	payload *artifactPayload[T]
	raw     []byte
}

type artifactPayload[T sArtifactPayload] struct {
	Timestamp int64
	Nonce     util.Nonce
	Payload   T
}

func NewSignedArtifact[T sArtifactPayload](credentials Credentials, payload T) (*SignedArtifact[T], error) {

	raw, err := MarshalAndSign(payload, credentials)
	if err != nil {
		return nil, err
	}

	return LoadSignedArtifact[T](raw)
}

func LoadSignedArtifact[T sArtifactPayload](raw []byte) (*SignedArtifact[T], error) {
	// TODO
	return nil, fmt.Errorf("not implemented")
}

func (s *SignedArtifact[T]) Payload() T {
	return s.payload.Payload
}

func (s *SignedArtifact[T]) Creator() *PublicKey {
	return s.creator
}

func (s *SignedArtifact[T]) Timestamp() int64 {
	return s.payload.Timestamp
}

func (s *SignedArtifact[T]) Nonce() util.Nonce {
	return s.payload.Nonce
}

func (s *SignedArtifact[T]) Raw() []byte {
	return s.raw
}

func (s *SignedArtifact[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.raw)
}

func (s *SignedArtifact[T]) UnmarshalJSON(data []byte) error {
	raw := make([]byte, 0, len(data))
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return fmt.Errorf("failed to unmarshal raw artifact: %w", err)
	}

	artifact, err := LoadSignedArtifact[T](raw)
	if err != nil {
		return fmt.Errorf("failed to load signed artifact: %w", err)
	}

	*s = *artifact
	return nil
}
