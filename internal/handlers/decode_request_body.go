package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func DecodeRequestBody(w http.ResponseWriter, r *http.Request, payload interface{}) error {
	err := json.NewDecoder(r.Body).Decode(payload)

	switch {
	case err == io.EOF:
		JSONResponse(w, R{Message: "Request body must not be empty"}, http.StatusBadRequest)
		return err
	case err != nil:
		JSONResponse(w, R{Message: fmt.Sprintf("Error decode response body : %v", err.Error())}, http.StatusInternalServerError)
		return err
	}

	return nil
}
