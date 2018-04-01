package cryptx

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha512"
	"encoding/hex"
	"io"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/pbkdf2"

	"github.com/gavrilaf/spawn/pkg/errx"
)

const (
	ErrScope = "cpyptx"
)

var (
	InvalidSignature = errx.New(ErrScope, "invalid-signature")
)

/*
 * Generate key based on the seed and random salt.
 */
func GenerateSaltedKey(seed string) ([]byte, error) {
	salt := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, err
	}

	return pbkdf2.Key([]byte(seed), salt, 4096, 32, sha1.New), nil
}

/*
 * Generate signature for the message
 *  - message
 *  - key
 * Return signature in hex coding
 */
func GenerateSignature(message string, key []byte) string {
	mac := hmac.New(sha512.New, key)
	mac.Write([]byte(message))
	return hex.EncodeToString(mac.Sum(nil))
}

/*
 * Check message signature
 *  - message
 *  - signature (string in hex coding)
 *  - key
 * Return true is signature == HMAC(message, key)
 */
func CheckSignature(message string, signature string, key []byte) error {
	sign2, err := hex.DecodeString(signature)
	if err != nil {
		return err
	}

	mac := hmac.New(sha512.New, key)
	mac.Write([]byte(message))
	expectedSign := mac.Sum(nil)

	if !hmac.Equal(sign2, expectedSign) {
		return InvalidSignature
	}

	return nil
}

/*
 * Generate pasword hash based on user password
 *  - password
 * Return password hash in hex
 */
func GenerateHashedPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(hashedPassword), nil
}

/*
 * Check if password matches with hash
 * - password
 * - hash (in hex)
 * Return nil if matched
 */
func CheckPassword(password string, hash string) error {
	hash2, err := hex.DecodeString(hash)
	if err != nil {
		return err
	}

	return bcrypt.CompareHashAndPassword(hash2, []byte(password))
}
