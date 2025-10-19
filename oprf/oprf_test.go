package oprf

import (
	"bytes"
	"crypto/sha512"
	"testing"

	"github.com/gtank/ristretto255"
)

func TestOPRFTwoDifferentKeys(t *testing.T) {
	key, err := NewKeyPair()
	if err != nil {
		t.Fatal(err)
	}

	input := []byte("password123")

	r, blinded, err := Blind(input)
	if err != nil {
		t.Fatal(err)
	}

	evaluated := BlindEvaluate(key.Private, blinded)
	out1 := Finalize(r, evaluated)

	// simulate a different secret -> should produce different output
	key2, _ := NewKeyPair()
	evaluated2 := BlindEvaluate(key2.Private, blinded)
	out2 := Finalize(r, evaluated2)

	if bytes.Equal(out1, out2) {
		t.Error("Outputs should differ for different server keys")
	}

	if len(out1) != 32 {
		t.Errorf("Expected 32-byte output, got %d", len(out1))
	}
}

func TestOPRFEndToEnd(t *testing.T) {
	key, err := NewKeyPair()
	if err != nil {
		t.Fatal(err)
	}

	input := []byte("testinput")

	// Hash input to group element
	P, err := hashToGroup(input)
	if err != nil {
		t.Fatal(err)
	}

	// Direct calculation: input * secret key
	direct := new(ristretto255.Element).ScalarMult(key.Private, P)
	directBytes := direct.Bytes()
	directHash := sha512.Sum512(directBytes)
	expected := directHash[:32]

	// OPRF calculation
	r, blinded, err := Blind(input)
	if err != nil {
		t.Fatal(err)
	}
	evaluated := BlindEvaluate(key.Private, blinded)
	out := Finalize(r, evaluated)

	// Compare direct and OPRF output
	if !bytes.Equal(out, expected) {
		t.Error("Direct calculation and OPRF output differ")
	}
}
