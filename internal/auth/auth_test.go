package auth

import "testing"

func TestHashPassword(t *testing.T) {
	password := "mysecretpassword"

	// Test hashing the password
	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Error hashing password: %v", err)
	}

	// The hashed password should not equal the original password
	if hashedPassword == password {
		t.Error("Hashed password should not be the same as the original password")
	}

	// Test verifying the correct password
	err = CheckPasswordHash(password, hashedPassword)
	if err != nil {
		t.Errorf("Expected password to match hash, got error: %v", err)
	}

	// Test verifying an incorrect password
	wrongPassword := "wrongpassword"
	err = CheckPasswordHash(wrongPassword, hashedPassword)
	if err == nil {
		t.Error("Expected error when verifying wrong password, but got nil")
	}
}
