package utils

import (
	"net/http"
)

func IsHTTPMethodCorrect(w http.ResponseWriter, r *http.Request, Method string) bool {
	if r.Method != Method {
		JSONResponse(w, R{Message: "Wrong HTTP request method"}, http.StatusBadRequest)
		return false
	}
	return true
}
