package oauth2

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

const (
	codeVerifierLength  = 32
	stateLength         = 80
	codeChallengeMethod = "S256"
)

type Session struct {
	State               string
	CodeVerifier        string
	CodeChallenge       string
	CodeChallengeMethod string
}

func newSession() *Session {
	r := getRandomBytes(stateLength)
	state := base64.RawURLEncoding.EncodeToString(r)

	r = getRandomBytes(codeVerifierLength)
	codeVerifier := base64.RawURLEncoding.EncodeToString(r)

	h := sha256.Sum256([]byte(codeVerifier))
	codeChallenge := base64.RawURLEncoding.EncodeToString(h[:])

	fmt.Println(len(codeChallenge))

	return &Session{
		State:               state,
		CodeVerifier:        codeVerifier,
		CodeChallenge:       codeChallenge,
		CodeChallengeMethod: codeChallengeMethod,
	}
}

func getRandomBytes(l int) []byte {
	b := make([]byte, l)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return b
}
