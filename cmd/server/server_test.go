package main

import "testing"

func TestGenerateChallenge(t *testing.T) {
	challenge, err := generateChallenge()
	if err != nil {
		t.Errorf("Error generating challenge: %v", err)
	}
	if len(challenge) != 32 {
		t.Errorf("Expected challenge length of 32, got %d", len(challenge))
	}
}

func TestVerifyResponse(t *testing.T) {
	challenge := "a1b2c3d4e5f6g7h8"
	response := "103579"
	if !verifyResponse(challenge, response, 4) {
		t.Error("Expected response to be invalid")
	}
}

func TestVerifyResponseInvalid(t *testing.T) {
	challenge := "a1b2c3d4e5f6g7h8"
	response := "10"
	if verifyResponse(challenge, response, 4) {
		t.Error("Expected response to be invalid")
	}
}
