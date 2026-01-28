package auth

import (
	"golang.org/x/crypto/bcrypt"
)

func VerifyPassword(hash, password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return err
	}
	return nil
}

func GenerateHashPassword(password string) (hashedPassword []byte, err error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	if err != nil {
		return nil, err
	}

	return hash, nil
}
