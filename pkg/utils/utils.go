package utils

import "golang.org/x/crypto/bcrypt"

// Hashing generates a bcrypt hash for the given string.
// It returns the hashed string or an error if the hashing fails.
func Hashing(s string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(s), 10)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// CheckHashing compares a plaintext string with a bcrypt hash.
// It returns true if the string matches the hash, otherwise false.
func CheckHashing(s, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(s))
	return err == nil
}
