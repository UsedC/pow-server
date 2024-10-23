package main

import (
	"testing"
)

func TestComputeNonce(t *testing.T) {
	expectedNonce := 103579
	nonce := computeNonce("a1b2c3d4e5f6g7h8", 4)
	if nonce != expectedNonce {
		t.Errorf("Expected nonce %d, got %d", expectedNonce, nonce)
	}
}
