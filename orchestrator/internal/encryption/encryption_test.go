package encryption

import "testing"

func TestEncryption(t *testing.T) {
	var password1 string
	var err error
	password1, err = Generate("12345678")
	if err != nil {
		t.Fatal("Unexpected error occured")
	}
	err = Compare(password1, "1234567")
	if err == nil {
		t.Fatal("Error expected")
	}

	//русская "о"
	password1, err = Generate("sоmething")
	if err != nil {
		t.Fatal("Unexpected error occured")
	}
	//английская "o"
	err = Compare(password1, "something")
	if err == nil {
		t.Fatal("Error expected")
	}
}
