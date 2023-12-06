package pkg

import (
	"crypto/sha1"
	"fmt"
)

const (
	salt = "mgfd#gg*049u3tnfdld%"
)

func GetPasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
