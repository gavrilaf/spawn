package cryptx

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	msgs := []string{"12345", "This is a very long key", "", "client_test", "client_ios", "client_testd1test-user"}

	key := []byte("client_test_key")

	var signs []string
	for _, s := range msgs {
		sn := GenerateSignature(s, key)
		signs = append(signs, sn)
		fmt.Printf("Signatire for %v is %v\n", s, sn)
	}

	for i, s := range msgs {
		err := CheckSignature(s, signs[i], key)
		assert.Nil(t, err)
	}

	err := CheckSignature(msgs[0]+"1", signs[0], key)
	assert.NotNil(t, err)
}

func TestPasswords(t *testing.T) {
	psws := []string{"111111", "This is a very long password", "test", "", "id1-password", "id2-password"}

	var hashes []string
	for _, s := range psws {
		hash, err := GenerateHashedPassword(s)
		require.Nil(t, err)
		fmt.Printf("Pasword hash for %v is %v\n", s, hash)

		hashes = append(hashes, hash)
	}

	for i, s := range psws {
		err := CheckPassword(s, hashes[i])
		assert.Nil(t, err)
	}

	err := CheckPassword(psws[0]+"1", hashes[0])
	assert.NotNil(t, err)
}

func TestGenerateConfirm(t *testing.T) {
	s1 := GenerateConfimCode()
	s2 := GenerateConfimCode()
	s3 := GenerateConfimCode()

	assert.Equal(t, ConfirmCodeLength, len(s1))
	assert.Equal(t, ConfirmCodeLength, len(s2))
	assert.Equal(t, ConfirmCodeLength, len(s3))
}
