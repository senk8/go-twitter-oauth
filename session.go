package go_twitter_oauth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
)

const (
	codeVerifierLength = 80
	stateLength        = 80
	codeChallengeMethod = "S256"
)

type Session struct {
	State               string
	CodeVerifier        string
	CodeChallenge       string
	CodeChallengeMethod string
}

func newSession() *Session {
	// state
	r := getRandomString(stateLength)
	state := base64.RawURLEncoding.EncodeToString(r)

	// code_verifier 43~128 string
	r = getRandomString(codeVerifierLength)
	codeVerifier := base64.RawURLEncoding.EncodeToString(r)

	// code_challenge 43~128 string
	h := sha256.Sum256([]byte(codeVerifier))
	codeChallenge := base64.RawURLEncoding.EncodeToString(h[:])

	return &Session{
		State:               state,
		CodeVerifier:        codeVerifier,
		CodeChallenge:       codeChallenge,
		CodeChallengeMethod: codeChallengeMethod,
	}
}

func getRandomString(l int) []byte {
	b := make([]byte, l)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return b
}
