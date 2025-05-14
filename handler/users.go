package handler

import (
	"backend/data"
	"backend/utils"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {

	if !utils.IsHTTPMethodCorrect(w, r, "POST") {
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		utils.JSONResponse(w,R{Message: "Content-Type must application/json"}, http.StatusUnsupportedMediaType)
		return
	}

	payload := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}

	if err := utils.DecodeRequestBody(w, r, &payload); err != nil {
		return
	}

	if utils.HasEmptyField(w, payload) {
		return
	}

	var result data.User

	err := h.db.Collection("users").FindOne(ctx, bson.M{"username": payload.Username}).Decode(&result)

	if err == mongo.ErrNoDocuments {
		utils.JSONResponse(w, R{Message: "Data not found"}, http.StatusNotFound)
		return
	}

	if err != nil {
		utils.JSONResponse(w, R{Message: fmt.Sprintf("Error fetch data: %v", err.Error())}, http.StatusInternalServerError)
		return
	}

	if !utils.IsPasswordCorrect(w, result.Password, payload.Password) {
		return
	}

	session, _ := h.store.Get(r, data.SESSION_ID)
	session.Values["userID"] = result.UserId.Hex()
	session.Save(r, w)

	responseData := struct {
		Email     string `json:"email"`
		Username    string `json:"username"`
	}{result.Email, result.Username}


	utils.JSONResponse(w, R{Message: "Login successful", Data: responseData}, http.StatusOK)
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	if !utils.IsHTTPMethodCorrect(w, r, "POST") {
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		utils.JSONResponse(w,R{Message: "Content-Type must application/json"}, http.StatusUnsupportedMediaType)
		return
	}

	payload := struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}

	if err := utils.DecodeRequestBody(w, r, &payload); err != nil {
		return
	}

	if utils.HasEmptyField(w, payload) {
		return
	}

	hashedPassword, err := utils.GenerateHashPassword(w, payload.Password)
	if err != nil {
		return
	}

	newUser := data.User{Username: payload.Username, Email: payload.Email, Password: string(hashedPassword)}

	_, err = h.db.Collection("users").InsertOne(ctx, newUser)
	if mongo.IsDuplicateKeyError(err) {
		utils.JSONResponse(w, R{
			Message: "Username is already taken",
		}, http.StatusConflict)
		return
	}

	if err != nil {
		utils.JSONResponse(w, R{
			Message: fmt.Sprintf("Error create user: %v", err.Error()),
		}, http.StatusInternalServerError)
		return
	}

	utils.JSONResponse(w, R{Message: "User created successfully"}, http.StatusCreated)
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	if !utils.IsHTTPMethodCorrect(w, r, "DELETE") {
		return
	}

	payload := struct {
		UserId string `json:"userId"`
	}{}

	if err := utils.DecodeRequestBody(w, r, &payload); err != nil {
		return
	}

	if utils.HasEmptyField(w, payload) {
		return
	}

	// Kalau pakai mongo v2 pakai ini kalau cari data berdasarkan id, idnya ubah ke obj
	// objID, err := bson.ObjectIDFromHex(payload.UserId)

	// Kalau pakai mongo v1
	objID, err := primitive.ObjectIDFromHex(payload.UserId)

	if err != nil {
		utils.JSONResponse(w, R{
			Message: fmt.Sprintf("Server error : %v", err.Error()),
		}, http.StatusInternalServerError)
		return
	}

	var deletedUser data.User
	err = h.db.Collection("users").FindOneAndDelete(ctx, bson.D{{Key: "_id", Value: objID}}).Decode(&deletedUser)

	if err == mongo.ErrNoDocuments {
		utils.JSONResponse(w, R{Message: "Data not found"}, http.StatusNotFound)
		return
	}

	if err != nil {
		utils.JSONResponse(w, R{
			Message: fmt.Sprintf("Error delete data : %v", err.Error()),
		}, http.StatusInternalServerError)
		return
	}

	utils.JSONResponse(w, R{
		Message: "Data deleted successfully",
		Data:    deletedUser,
	}, http.StatusOK)

}

func (h *Handler) GetLoggedInUser(w http.ResponseWriter, r *http.Request) {
	if !utils.IsHTTPMethodCorrect(w, r, "GET") {
		return
	}

	session, _ := h.store.Get(r, data.SESSION_ID)

	id := session.Values["userID"].(string) // interface{} -> string


	// Kalau pakai mongo v2 pakai ini kalau cari data b	err = h.db.Collection("users").FindOne(ctx, bson.M{"_id": objID}).Decode(&result)
	// erdasarkan id, idnya ubah ke obj
	// objID, err := bson.ObjectIDFromHex(payload.UserId)

	// Kalau pakai mongo v1
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		utils.JSONResponse(w, R{
			Message: fmt.Sprintf("Server error : %v", err.Error()),
		}, http.StatusInternalServerError)
		return
	}

	var result data.User

	err = h.db.Collection("users").FindOne(ctx, bson.D{{Key: "_id", Value: objID}}).Decode(&result)

	if err != nil {
		utils.JSONResponse(w, R{
			Message: fmt.Sprintf("Error fetch data : %v", err.Error()),
		}, http.StatusInternalServerError)
		return
	}

	data := struct {
		Username string `json:"username"`
		Email    string `json:"email"`
	}{result.Username, result.Email}

	utils.JSONResponse(w, R{
		Message: "Data fetched successfully", Data: data,
	}, http.StatusOK)
}
