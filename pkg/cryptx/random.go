package cryptx

import (
	"crypto/rand"
	"io"
)

const (
	ConfirmCodeLength = 6
)

var numbers = []byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

func GenerateConfimCode() string {
	return GenRandomString(numbers, ConfirmCodeLength)
}

func GenRandomString(table []byte, length int) string {
	b := make([]byte, length)
	n, _ := io.ReadAtLeast(rand.Reader, b, length)
	if n != length {
		return ""
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}
