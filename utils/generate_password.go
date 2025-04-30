package utils

import (
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func GenerateHashPassword(w http.ResponseWriter, password string) (hashedPassword []byte, err error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	if err != nil {
		JSONResponse(w, R{Message: fmt.Sprintf("Server error: %v", err.Error())}, http.StatusInternalServerError)
		return nil, err
	}

	return hash, nil
}
