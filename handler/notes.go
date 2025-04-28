package handler

import (
	"backend/data"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func (h *Handler) CreateNote(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Use POST Method", http.StatusBadRequest)
		return
	}

	payload := struct {
		UserId string `json:"userId"`
		Title  string `json:"title"`
		Text   string `json:"text"`
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

	if payload.UserId == "" || payload.Title == "" || payload.Text == "" {
		http.Error(w, "No UserId", http.StatusBadRequest)
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

	_, err = h.db.Collection("notes").InsertOne(ctx, newNote)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Insert Note Success"))

}

func (h *Handler) DeleteNote(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		http.Error(w, "Use GET Method", http.StatusBadRequest)
		return
	}

	payload := struct {
		NoteId string `json:"noteId"`
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

	if payload.NoteId == "" {
		http.Error(w, "No NoteId", http.StatusBadRequest)
		return
	}

	objID, err := bson.ObjectIDFromHex(payload.NoteId)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = h.db.Collection("notes").DeleteOne(ctx, bson.D{{Key: "_id", Value: objID}})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Delete Note Success"))

}

func (h *Handler) UpdateNote(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Use POST Method", http.StatusBadRequest)
		return
	}

	payload := struct {
		NoteId string `json:"noteId"`
		Title  string `json:"title"`
		Text   string `json:"text"`
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

	if payload.NoteId == "" || payload.Title == "" || payload.Text == "" {
		http.Error(w, "No UserId", http.StatusBadRequest)
		return
	}

	var updatedAt = time.Now().String()

	objID, _ := bson.ObjectIDFromHex(payload.NoteId)

	filter := bson.D{{Key: "_id", Value: objID}}

	update := bson.D{
		{Key: "$set",
			Value: bson.D{
				{Key: "title", Value: payload.Title},
				{Key: "text", Value: payload.Text},
				{Key: "updatedAt", Value: updatedAt},
			}},
	}
	_, err = h.db.Collection("notes").UpdateOne(ctx, filter, update)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Insert Note Success"))
}

func (h *Handler) FindNoteById(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Use GET Method", http.StatusBadRequest)
		return
	}

	payload := struct {
		NoteId string `json:"noteId"`
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

	if payload.NoteId == "" {
		http.Error(w, "No NoteId", http.StatusBadRequest)
		return
	}

	objID, err := bson.ObjectIDFromHex(payload.NoteId)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	var result data.Note

	err = h.db.Collection("notes").FindOne(ctx, bson.M{"_id": objID}).Decode(&result)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	data := []struct {
		UserId    string
		Title     string
		Text      string
		CreatedAt string
		UpdatedAt string
	}{{result.UserId, result.Title, result.Text, result.CreatedAt, result.UpdatedAt}}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) FindUserNotes(w http.ResponseWriter, r *http.Request) {
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

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	filter := bson.D{{Key: "userId", Value: payload.UserId}}

	cursor, err := h.db.Collection("notes").Find(ctx, filter)

	// fmt.Println(cursor)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var results []data.Note

	if err = cursor.All(ctx, &results); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println(results)

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(results)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
