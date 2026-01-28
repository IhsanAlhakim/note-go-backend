package handlers

import (
	"backend/internal/auth"
	"backend/internal/database"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {

	payload := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}

	if err := BindJSON(r, &payload); err != nil {
		if err == ErrEmptyBody {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if payload.Username == "" || payload.Password == "" {
		http.Error(w, "missing login credentials", http.StatusBadRequest)
		return
	}

	var result database.User

	err := h.db.Collection("users").FindOne(ctx, bson.M{"username": payload.Username}).Decode(&result)

	if err == mongo.ErrNoDocuments {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := auth.VerifyPassword(result.Password, payload.Password); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			http.Error(w, "Invalid Credentials", http.StatusUnauthorized)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session, _ := h.store.Get(r, h.cfg.SessionID)
	session.Values["userID"] = result.UserId.Hex()
	session.Save(r, w)

	responseData := struct {
		Email    string `json:"email"`
		Username string `json:"username"`
	}{result.Email, result.Username}

	RespondJSON(w, R{Message: "Login successful", Data: responseData}, http.StatusOK)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {

	session, err := h.store.Get(r, h.cfg.SessionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	session.Options.MaxAge = -1
	session.Save(r, w)
	RespondJSON(w, R{Message: "Logout successful"}, http.StatusOK)
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	payload := struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}

	if err := BindJSON(r, &payload); err != nil {
		if err == ErrEmptyBody {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if payload.Email == "" || payload.Username == "" || payload.Password == "" {
		http.Error(w, "missing credentials", http.StatusBadRequest)
		return
	}

	hashedPassword, err := auth.GenerateHashPassword(payload.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	newUser := database.User{Username: payload.Username, Email: payload.Email, Password: string(hashedPassword)}

	_, err = h.db.Collection("users").InsertOne(ctx, newUser)
	if mongo.IsDuplicateKeyError(err) {
		http.Error(w, "Username is already taken", http.StatusConflict)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	RespondJSON(w, R{Message: "User created successfully"}, http.StatusCreated)
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {

	session, err := h.store.Get(r, h.cfg.SessionID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	id := session.Values["userID"].(string) // interface{} -> string

	// Kalau pakai mongo v2 pakai ini kalau cari data berdasarkan id, idnya ubah ke obj
	// objID, err := bson.ObjectIDFromHex(payload.UserId)

	// Kalau pakai mongo v1
	objID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Memakai transaction untuk menghapus data user dan note user
	// jika salah satu gagal, maka seluruh operasi akan dibatalkan (rollback) untuk menjaga integritas data

	mongoSession, err := h.client.StartSession()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer mongoSession.EndSession(ctx)

	var deletedUser database.User
	_, err = mongoSession.WithTransaction(ctx, func(sc mongo.SessionContext) (interface{}, error) {
		// Hapus user
		if err = h.db.Collection("users").FindOneAndDelete(ctx, bson.D{{Key: "_id", Value: objID}}).Decode(&deletedUser); err != nil {
			return nil, err // data di rollback jika proses ini gagal
		}

		// Hapus notes milik user
		deleteNotesFilter := bson.D{{Key: "userId", Value: id}}
		if _, err = h.db.Collection("notes").DeleteMany(ctx, deleteNotesFilter); err != nil {
			return nil, err // di rollback jika gagal
		}

		return nil, nil
	})

	if err == mongo.ErrNoDocuments {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Options.MaxAge = -1
	session.Save(r, w)

	RespondJSON(w, R{
		Message: "Data deleted successfully",
	}, http.StatusOK)

}

func (h *Handler) GetLoggedInUser(w http.ResponseWriter, r *http.Request) {

	session, _ := h.store.Get(r, h.cfg.SessionID)

	id := session.Values["userID"].(string) // interface{} -> string

	// Kalau pakai mongo v2 pakai ini kalau cari data b	err = h.db.Collection("users").FindOne(ctx, bson.M{"_id": objID}).Decode(&result)
	// erdasarkan id, idnya ubah ke obj
	// objID, err := bson.ObjectIDFromHex(payload.UserId)

	// Kalau pakai mongo v1
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var result database.User

	err = h.db.Collection("users").FindOne(ctx, bson.D{{Key: "_id", Value: objID}}).Decode(&result)

	if err == mongo.ErrNoDocuments {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		Username string `json:"username"`
		Email    string `json:"email"`
	}{result.Username, result.Email}

	RespondJSON(w, R{
		Message: "Data fetched successfully", Data: data,
	}, http.StatusOK)
}
