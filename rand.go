package gotwtr_oauth

import "crypto/rand"

func GetRandomString(l int) []byte {
	b := make([]byte, l)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return b
}
