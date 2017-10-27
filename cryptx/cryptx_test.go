package cryptx

import (
	//"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSaltedKey(t *testing.T) {
	seeds := []string{"12345", "This is a very long key", "", "client_test", "client_ios"}
	for _, s := range seeds {
		key, err := GenerateSaltedKey(s)

		require.Nil(t, err)
		require.NotEmpty(t, key)

		//fmt.Printf("Seed = %v, key = %v\n", s, key)
	}
}

func TestSignatures(t *testing.T) {
	msgs := []string{"12345", "This is a very long key", "", "client_test", "client_ios"}

	key, err := GenerateSaltedKey("key-seed")
	require.Nil(t, err)

	var signs []string
	for _, s := range msgs {
		sn := GenerateSignature(s, key)
		signs = append(signs, sn)
	}
	//fmt.Printf("Generated signatires: %v\n", signs)

	for i, s := range msgs {
		err := CheckSignature(s, signs[i], key)
		assert.Nil(t, err)
	}

	err = CheckSignature(msgs[0]+"1", signs[0], key)
	assert.NotNil(t, err)
}

func TestPasswords(t *testing.T) {
	psws := []string{"111111", "This is a very long password", "test", ""}

	var hashes []string
	for _, s := range psws {
		hash, err := GenerateHashedPassword(s)
		require.Nil(t, err)

		hashes = append(hashes, hash)
	}
	//fmt.Printf("Hashed passwords: %v\n", hashes)

	for i, s := range psws {
		err := CheckPassword(s, hashes[i])
		assert.Nil(t, err)
	}

	err := CheckPassword(psws[0]+"1", hashes[0])
	assert.NotNil(t, err)
}
