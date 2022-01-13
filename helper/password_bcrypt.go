package helper

import "golang.org/x/crypto/bcrypt"

// HashPassword help to encrypt raw passwords and then insert to database
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}
