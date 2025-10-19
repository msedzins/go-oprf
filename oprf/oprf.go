package oprf

import (
	"crypto/rand"
	"crypto/sha512"
	"errors"
	"io"

	"github.com/gtank/ristretto255"
)

// Blind: client blinds input value x -> returns blinded element + blinding scalar.
func Blind(x []byte) (*ristretto255.Scalar, *ristretto255.Element, error) {
	var rBytes [64]byte
	if _, err := io.ReadFull(rand.Reader, rBytes[:]); err != nil {
		return nil, nil, err
	}

	r, err := new(ristretto255.Scalar).SetUniformBytes(rBytes[:])
	if err != nil {
		return nil, nil, err
	}

	P, err := hashToGroup(x)
	if err != nil {
		return nil, nil, err
	}

	blinded := new(ristretto255.Element).ScalarMult(r, P)
	return r, blinded, nil
}

// BlindEvaluate: server multiplies blinded element by secret key.
func BlindEvaluate(sk *ristretto255.Scalar, blinded *ristretto255.Element) *ristretto255.Element {
	return new(ristretto255.Element).ScalarMult(sk, blinded)
}

// Finalize: client unblinds server’s response.
func Finalize(r *ristretto255.Scalar, evaluated *ristretto255.Element) []byte {
	rInv := new(ristretto255.Scalar).Invert(r)
	unblinded := new(ristretto255.Element).ScalarMult(rInv, evaluated)

	unblindedBytes := unblinded.Bytes()
	hash := sha512.Sum512(unblindedBytes)
	return hash[:32] // 256-bit output
}

// hashToGroup maps arbitrary bytes to a Ristretto255 point.
func hashToGroup(x []byte) (*ristretto255.Element, error) {
	h := sha512.Sum512(x)
	elem, err := new(ristretto255.Element).SetUniformBytes(h[:])
	if err != nil {
		return nil, errors.New("failed to set uniform bytes")
	}
	return elem, nil
}
