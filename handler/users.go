package handler

import (
	"backend/data"
	"encoding/json"
	"io"
	"net/http"

	"go.mongodb.org/mongo-driver/v2/bson"
	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Use POST Method", http.StatusBadRequest)
		return
	}

	payload := struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&payload)

	switch {
	case err == io.EOF:
		http.Error(w, "No Request Body", http.StatusBadRequest)
		return
	case err != nil:
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if payload.Username == "" || payload.Email == "" || payload.Password == "" {
		http.Error(w, "No UserId", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), 14)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	newUser := data.User{Username: payload.Username, Email: payload.Email, Password: string(hashedPassword)}

	_, err = h.db.Collection("users").InsertOne(ctx, newUser)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Insert Success"))

	// mongo.IsDuplicateKeyError()
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		http.Error(w, "Use GET Method", http.StatusBadRequest)
		return
	}

	payload := struct {
		UserId string `json:"userId"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&payload)

	switch {
	case err == io.EOF:
		http.Error(w, "No Request Body", http.StatusBadRequest)
		return
	case err != nil:
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if payload.UserId == "" {
		http.Error(w, "No UserId", http.StatusBadRequest)
		return
	}

	objID, err := bson.ObjectIDFromHex(payload.UserId)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = h.db.Collection("users").DeleteOne(ctx, bson.D{{Key: "_id", Value: objID}})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Delete Success"))

}

func (h *Handler) FindUserById(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Use GET Method", http.StatusBadRequest)
		return
	}

	payload := struct {
		UserId string `json:"userId"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&payload)

	switch {
	case err == io.EOF:
		http.Error(w, "No Request Body", http.StatusBadRequest)
		return
	case err != nil:
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if payload.UserId == "" {
		http.Error(w, "No UserId", http.StatusBadRequest)
		return
	}

	objID, err := bson.ObjectIDFromHex(payload.UserId)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	var result data.User

	err = h.db.Collection("users").FindOne(ctx, bson.M{"_id": objID}).Decode(&result)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	data := []struct {
		Username string
		Email    string
	}{{result.Username, result.Email}}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
