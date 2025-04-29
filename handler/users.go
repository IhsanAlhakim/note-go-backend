package handler

import (
	"backend/data"
	"backend/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {

	payload := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&payload)

	switch {
	case err == io.EOF:
		utils.JSONResponse(w, map[string]interface{}{
			"message": "Request body must not be empty",
		}, http.StatusBadRequest)
		return
	case err != nil:
		utils.JSONResponse(w, map[string]interface{}{
			"message": fmt.Sprintf("Error decode response body : %v", err.Error()),
		}, http.StatusInternalServerError)
		return
	}

	if payload.Username == "" || payload.Password == "" {
		utils.JSONResponse(w, map[string]interface{}{
			"message": "Missing Credentials",
		}, http.StatusBadRequest)
		return
	}

	var result data.User

	err = h.db.Collection("users").FindOne(ctx, bson.M{"username": payload.Username}).Decode(&result)

	if err != nil {
		utils.JSONResponse(w, map[string]interface{}{
			"message": fmt.Sprintf("Database error: %v", err.Error()),
		}, http.StatusInternalServerError)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(payload.Password))

	if err != nil {
		http.Error(w, "Unauthorized User", http.StatusForbidden)
		return
	}

	// data := []struct {
	// 	Username string
	// 	Email    string
	// 	password string
	// }{{result.Username, result.Email}}

	// sessionCookieValue := utils.GenerateRandomString(10)
	// fmt.Println(sessionCookieValue)

	session, _ := h.store.Get(r, data.SESSION_ID)
	// session, _ := h.store.Get(r, result.Username)

	session.Values["userID"] = result.UserId.Hex()
	session.Save(r, w)
	// data.SetCookie(w, "notego", data.M{"CookieValue": sessionCookieValue})
	// fmt.Println("Session and Cookie Created")
	w.Write([]byte("Logged In"))
}

func (h *Handler) GetAuthenticatedUser(w http.ResponseWriter, r *http.Request) {
	// data, _ := data.GetCookie(r, "Sandalman")
	// fmt.Println(data)

	session, _ := h.store.Get(r, data.SESSION_ID)

	fmt.Println(session.Values["userID"])

	if session.Values["userID"] == nil {
		http.Error(w, "User No Authenticated", http.StatusBadRequest)
		return
	}

	// if len(session.Values) == 0 {
	// 	http.Error(w, "empty result", http.StatusOK)
	// }
	// fmt.Println(session.Values)
	w.Write([]byte("User Authenticated"))
}

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

	// Kalau pakai mongo v2 pakai ini kalau cari data berdasarkan id, idnya ubah ke obj
	// objID, err := bson.ObjectIDFromHex(payload.UserId)

	// Kalau pakai mongo v1
	objID, err := primitive.ObjectIDFromHex(payload.UserId)

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

	// Kalau pakai mongo v2 pakai ini kalau cari data berdasarkan id, idnya ubah ke obj
	// objID, err := bson.ObjectIDFromHex(payload.UserId)

	// Kalau pakai mongo v1
	objID, err := primitive.ObjectIDFromHex(payload.UserId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	var result data.User

	fmt.Println("Cari User")
	err = h.db.Collection("users").FindOne(ctx, bson.M{"_id": objID}).Decode(&result)

	fmt.Println(result)
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
