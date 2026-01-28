package utils

import (
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func IsPasswordCorrect(w http.ResponseWriter, hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	if err == bcrypt.ErrMismatchedHashAndPassword {
		JSONResponse(w, R{Message: "Incorrect login credentials"}, http.StatusUnauthorized)
		return false
	}

	if err != nil {
		JSONResponse(w, R{Message: fmt.Sprintf("Server error: %v", err.Error())}, http.StatusInternalServerError)
		return false
	}

	return true
}
