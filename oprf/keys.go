package oprf

import (
	"crypto/rand"
	"io"

	"github.com/gtank/ristretto255"
)

// KeyPair represents server-side OPRF key material.
type KeyPair struct {
	Private *ristretto255.Scalar
	Public  *ristretto255.Element
}

// NewKeyPair generates a random OPRF keypair.
func NewKeyPair() (*KeyPair, error) {
	var buf [64]byte
	if _, err := io.ReadFull(rand.Reader, buf[:]); err != nil {
		return nil, err
	}

	sk, err := new(ristretto255.Scalar).SetUniformBytes(buf[:])
	if err != nil {
		return nil, err
	}

	pk := new(ristretto255.Element).ScalarBaseMult(sk)
	return &KeyPair{Private: sk, Public: pk}, nil
}
