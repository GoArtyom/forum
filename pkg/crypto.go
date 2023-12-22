package pkg

import (
	"crypto/sha1"
	"fmt"
	"math/rand"
	"time"
)

const (
	salt = "mgfd#g5"
)

func GetPasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

func RandString(max int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%d", r.Intn(max+1))
}
