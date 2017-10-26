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
	fmt.Println(signs)

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
