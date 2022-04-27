package gotwtr_oauth

import (
	"crypto/sha256"
	"encoding/base64"
)

const (
	codeVerifierLength = 80
	stateLength        = 80
)

type PKCESession struct {
	State               string
	CodeVerifier        string
	CodeChallenge       string
	CodeChallengeMethod string
}

func CreatePKCESession() *PKCESession {
	r := GetRandomString(stateLength)
	state := base64.RawURLEncoding.EncodeToString(r)

	r = GetRandomString(codeVerifierLength)
	codeVerifier := base64.RawURLEncoding.EncodeToString(r)

	hashed := sha256.Sum256([]byte(codeVerifier))
	codeChallenge := base64.RawURLEncoding.EncodeToString(hashed[:])
	codeChallengeMethod := "S256"

	return &PKCESession{
		State:               state,
		CodeVerifier:        codeVerifier,
		CodeChallenge:       codeChallenge,
		CodeChallengeMethod: codeChallengeMethod,
	}
}
