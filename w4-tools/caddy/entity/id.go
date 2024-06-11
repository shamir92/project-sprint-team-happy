package entity

import (
	"crypto/rand"
	"encoding/base64"
	"log"
)

func GeneratePregeneratedId() string {
	id := make([]byte, 16)
	_, err := rand.Read(id)
	if err != nil {
		log.Printf("Failed to generate random ID: %v", err)
		return ""
	}
	return base64.URLEncoding.EncodeToString(id)
}
