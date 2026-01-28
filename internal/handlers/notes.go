package handlers

import (
	"backend/internal/database"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (h *Handler) CreateNote(w http.ResponseWriter, r *http.Request) {
	payload := struct {
		Title string `json:"title"`
		Text  string `json:"text"`
	}{}

	if err := BindJSON(r, &payload); err != nil {
		if err == ErrEmptyBody {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if payload.Text == "" && payload.Title == "" {
		http.Error(w, "Note text and title cannot both be empty. Only one of them", http.StatusBadRequest)
		return
	}

	session, _ := h.store.Get(r, h.cfg.SessionID)

	id := session.Values["userID"].(string)

	var createdAt = time.Now().Format(time.RFC3339)
	var updatedAt = createdAt

	newNote := database.Note{
		UserId:    id,
		Title:     payload.Title,
		Text:      payload.Text,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	_, err := h.db.Collection("notes").InsertOne(ctx, newNote)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	RespondJSON(w, R{Message: "Note created"}, http.StatusCreated)
}

func (h *Handler) DeleteNote(w http.ResponseWriter, r *http.Request) {
	noteId := r.PathValue("id")

	if noteId == "" {
		http.Error(w, "Missing required parameter: noteId", http.StatusBadRequest)
		return
	}

	objID, err := primitive.ObjectIDFromHex(noteId)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = h.db.Collection("notes").FindOneAndDelete(ctx, bson.D{{Key: "_id", Value: objID}}).Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "data not found", http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	RespondJSON(w, R{Message: "Note deleted successfully"}, http.StatusOK)
}

func (h *Handler) UpdateNote(w http.ResponseWriter, r *http.Request) {
	noteId := r.PathValue("id")

	payload := struct {
		Title string `json:"title"`
		Text  string `json:"text"`
	}{}

	if err := BindJSON(r, &payload); err != nil {
		if err == ErrEmptyBody {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if noteId == "" {
		http.Error(w, "Missing required parameter: noteId", http.StatusBadRequest)
		return
	}

	var updatedAt = time.Now().Format(time.RFC3339)

	objID, err := primitive.ObjectIDFromHex(noteId)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

	if err = h.db.Collection("notes").FindOneAndUpdate(ctx, filter, update).Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Data not found", http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	RespondJSON(w, R{Message: "Note updated successfully"}, http.StatusOK)
}

func (h *Handler) FindNoteById(w http.ResponseWriter, r *http.Request) {
	noteId := r.PathValue("id")

	if noteId == "" {
		http.Error(w, "Missing required parameter: noteId", http.StatusBadRequest)
		return
	}

	objID, err := primitive.ObjectIDFromHex(noteId)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var result database.Note

	if err = h.db.Collection("notes").FindOne(ctx, bson.M{"_id": objID}).Decode(&result); err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Data not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	RespondJSON(w, R{Message: "Data fetched successfully", Data: result}, http.StatusOK)
}

func (h *Handler) FindUserNotes(w http.ResponseWriter, r *http.Request) {
	session, _ := h.store.Get(r, h.cfg.SessionID)

	id := session.Values["userID"].(string)

	filter := bson.D{{Key: "userId", Value: id}}

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "_id", Value: -1}})

	cursor, err := h.db.Collection("notes").Find(ctx, filter, findOptions)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var results []database.Note

	if err = cursor.All(ctx, &results); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	RespondJSON(w, R{Message: "Data fetched successfully", Data: results}, http.StatusOK)
}
