package cryptx

import (
	"fmt"
	"testing"
)

func TestSaltedKey(t *testing.T) {
	seeds := []string{"12345", "This is a very long key", "", "client_test", "client_ios"}
	for _, s := range seeds {
		key, err := GenerateSaltedKey(s)
		if err != nil {
			t.Fatal(err.Error())
		}
		fmt.Printf("Seed = %v, key = %v\n", s, key)
	}
}

func TestSignatures(t *testing.T) {
	msgs := []string{"12345", "This is a very long key", "", "client_test", "client_ios"}

	key, _ := GenerateSaltedKey("key-seed")

	var signs []string
	for _, s := range msgs {
		sn, err := GenerateSignature(s, key)
		if err != nil {
			t.Fatal(err.Error())
		}
		signs = append(signs, sn)
	}
	fmt.Printf("Generated signatires: %v\n", signs)

	for i, s := range msgs {
		err := CheckSignature(s, signs[i], key)
		if err != nil {
			t.Fatal(err.Error())
		}
	}

	err := CheckSignature(msgs[0]+"1", signs[0], key)
	if err == nil {
		t.Fatal("Signature should be invalid")
	}
}

func TestPasswords(t *testing.T) {
	psws := []string{"111111", "This is a very long password", "test", ""}

	var hashes []string
	for _, s := range psws {
		hash, err := GenerateHashedPassword(s)
		if err != nil {
			t.Fatal(err.Error())
		}
		hashes = append(hashes, hash)
	}
	fmt.Printf("Hashed passwords: %v\n", hashes)

	for i, s := range psws {
		err := CheckPassword(s, hashes[i])
		if err != nil {
			t.Fatal(err.Error())
		}
	}

	err := CheckPassword(psws[0]+"1", hashes[0])
	if err == nil {
		t.Fatal("Passwords must not match")
	}

}
