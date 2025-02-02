package utils

import (
	"crypto/rand"
	"log"
)

func GenerateRandomBytes() [20]byte {
	var b [20]byte
	_, err := rand.Read(b[:])
	if err != nil {
		log.Fatalf("Failed to generate random bytes: %v", err)
	}
	return b
}
