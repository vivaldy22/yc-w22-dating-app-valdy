package crypto

import (
	"golang.org/x/crypto/bcrypt"
)

// HashAndSalt hash and salt plain text using golang bcrypt
func HashAndSalt(plain string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashed), nil
}

// CheckPasswordHash check password hash with plain password
func CheckPasswordHash(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
