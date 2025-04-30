package handler

import (
	"backend/data"
	"backend/utils"
	"fmt"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (h *Handler) CreateNote(w http.ResponseWriter, r *http.Request) {
	if !utils.IsHTTPMethodCorrect(w, r, "POST") {
		return
	}

	payload := struct {
		UserId string `json:"userId"`
		Title  string `json:"title"`
		Text   string `json:"text"`
	}{}

	if err := utils.DecodeRequestBody(w, r, &payload); err != nil {
		return
	}

	if utils.HasEmptyField(w, payload) {
		return
	}

	var createdAt = time.Now().String()
	var updatedAt = createdAt

	newNote := data.Note{
		UserId:    payload.UserId,
		Title:     payload.Title,
		Text:      payload.Text,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	_, err := h.db.Collection("notes").InsertOne(ctx, newNote)

	if err != nil {
		utils.JSONResponse(w, R{Message: fmt.Sprintf("Server error: %v", err.Error())}, http.StatusInternalServerError)
		return
	}

	utils.JSONResponse(w, R{Message: "Note created"}, http.StatusCreated)
}

func (h *Handler) DeleteNote(w http.ResponseWriter, r *http.Request) {
	if !utils.IsHTTPMethodCorrect(w, r, "DELETE") {
		return
	}

	payload := struct {
		NoteId string `json:"noteId"`
	}{}

	if err := utils.DecodeRequestBody(w, r, &payload); err != nil {
		return
	}

	if utils.HasEmptyField(w, payload) {
		return
	}

	objID, err := primitive.ObjectIDFromHex(payload.NoteId)

	if err != nil {
		utils.JSONResponse(w, R{Message: fmt.Sprintf("Server error: %v", err.Error())}, http.StatusInternalServerError)
		return
	}

	var deletedNote data.Note
	err = h.db.Collection("notes").FindOneAndDelete(ctx, bson.D{{Key: "_id", Value: objID}}).Decode(&deletedNote)
	if err == mongo.ErrNoDocuments {
		utils.JSONResponse(w, R{Message: "Data not found"}, http.StatusNotFound)
		return
	}
	if err != nil {
		utils.JSONResponse(w, R{Message: fmt.Sprintf("Server error: %v", err.Error())}, http.StatusInternalServerError)
		return
	}

	utils.JSONResponse(w, R{Message: "Note deleted successfully", Data: deletedNote}, http.StatusOK)
}

func (h *Handler) UpdateNote(w http.ResponseWriter, r *http.Request) {
	if !utils.IsHTTPMethodCorrect(w, r, "PATCH") {
		return
	}

	payload := struct {
		NoteId string `json:"noteId"`
		Title  string `json:"title"`
		Text   string `json:"text"`
	}{}

	if err := utils.DecodeRequestBody(w, r, &payload); err != nil {
		return
	}

	if utils.HasEmptyField(w, payload) {
		return
	}

	var updatedAt = time.Now().String()

	objID, err := primitive.ObjectIDFromHex(payload.NoteId)

	if err != nil {
		utils.JSONResponse(w, R{Message: fmt.Sprintf("Server error: %v", err.Error())}, http.StatusInternalServerError)
		return
	}

	filter := bson.D{{Key: "_id", Value: objID}}
	update := bson.D{
		{Key: "$set",
			Value: bson.D{
				{Key: "title", Value: payload.Title},
				{Key: "text", Value: payload.Text},
				{Key: "updatedAt", Value: updatedAt},
			}},
	}

	var updatedNote data.Note
	err = h.db.Collection("notes").FindOneAndUpdate(ctx, filter, update).Decode(&updatedNote)

	if err == mongo.ErrNoDocuments {
		utils.JSONResponse(w, R{Message: "Data not found"}, http.StatusNotFound)
		return
	}

	if err != nil {
		utils.JSONResponse(w, R{Message: fmt.Sprintf("Server error: %v", err.Error())}, http.StatusInternalServerError)
		return
	}

	utils.JSONResponse(w, R{Message: "Note updated successfully", Data: updatedNote}, http.StatusOK)
}

func (h *Handler) FindNoteById(w http.ResponseWriter, r *http.Request) {
	if !utils.IsHTTPMethodCorrect(w, r, "GET") {
		return
	}

	payload := struct {
		NoteId string `json:"noteId"`
	}{}

	if err := utils.DecodeRequestBody(w, r, &payload); err != nil {
		return
	}

	if utils.HasEmptyField(w, payload) {
		return
	}

	objID, err := primitive.ObjectIDFromHex(payload.NoteId)

	if err != nil {
		utils.JSONResponse(w, R{Message: fmt.Sprintf("Server error: %v", err.Error())}, http.StatusInternalServerError)
		return
	}

	var result data.Note

	err = h.db.Collection("notes").FindOne(ctx, bson.M{"_id": objID}).Decode(&result)

	if err == mongo.ErrNoDocuments {
		utils.JSONResponse(w, R{Message: "Data not found"}, http.StatusNotFound)
		return
	}

	if err != nil {
		utils.JSONResponse(w, R{Message: fmt.Sprintf("Server error: %v", err.Error())}, http.StatusInternalServerError)
		return
	}

	data := struct {
		UserId    string `json:"userId"`
		Title     string `json:"title"`
		Text      string `json:"text"`
		CreatedAt string `json:"createdAt"`
		UpdatedAt string `json:"updatedAt"`
	}{result.UserId, result.Title, result.Text, result.CreatedAt, result.UpdatedAt}

	utils.JSONResponse(w, R{Message: "Data fetched successfully", Data: data}, http.StatusOK)
}

func (h *Handler) FindUserNotes(w http.ResponseWriter, r *http.Request) {
	if !utils.IsHTTPMethodCorrect(w, r, "GET") {
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

	filter := bson.D{{Key: "userId", Value: payload.UserId}}

	cursor, err := h.db.Collection("notes").Find(ctx, filter)

	if err != nil {
		utils.JSONResponse(w, R{Message: fmt.Sprintf("Server error: %v", err.Error())}, http.StatusInternalServerError)
		return
	}

	var results []data.Note

	if err = cursor.All(ctx, &results); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.JSONResponse(w, R{Message: "Data fetched successfully", Data: results}, http.StatusOK)
}
