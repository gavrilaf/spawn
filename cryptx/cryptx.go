package cryptx

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"github.com/gavrilaf/go-auth/errx"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/pbkdf2"
	"io"
)

const (
	ErrScope = "cpyptx"
)

var (
	InvalidSignature = errx.New(ErrScope, "invalid-signature")
)

/*
 * Generate key based on the seed and random salt. Return salted key in hex coding
 */
func GenerateSaltedKey(seed string) (string, error) {
	salt := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return "", err
	}

	key := pbkdf2.Key([]byte(seed), salt, 4096, 32, sha1.New)
	return hex.EncodeToString(key), nil
}

/*
 * Generate signature for the message
 *  - message
 *  - key (string in hex coding)
 * Return signature in hex coding
 */
func GenerateSignature(message string, key string) (string, error) {
	key2, err := hex.DecodeString(key)
	if err != nil {
		return "", nil
	}

	mac := hmac.New(sha256.New, key2)
	mac.Write([]byte(message))
	return hex.EncodeToString(mac.Sum(nil)), nil
}

/*
 * Check message signature
 *  - message
 *  - signature (string in hex coding)
 *  - key (hex)
 * Return true is signature == HMAC(message, key)
 */
func CheckSignature(message string, sign string, key string) error {
	sign2, err := hex.DecodeString(sign)
	if err != nil {
		return err
	}

	key2, err := hex.DecodeString(key)
	if err != nil {
		return err
	}

	mac := hmac.New(sha256.New, key2)
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
