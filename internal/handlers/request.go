package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

var ErrEmptyBody = errors.New("request body is empty")

func BindJSON(r *http.Request, payload any) error {
	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		if err == io.EOF {
			return ErrEmptyBody
		}
		return err
	}

	return nil
}
