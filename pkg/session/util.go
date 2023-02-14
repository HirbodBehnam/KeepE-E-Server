package session

import (
	"crypto/rand"
	"encoding/base32"
)

func newSessionID() string {
	b := make([]byte, 20)
	_, _ = rand.Read(b)
	return base32.StdEncoding.EncodeToString(b)
}
