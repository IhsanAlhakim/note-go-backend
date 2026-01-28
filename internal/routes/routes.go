package routes

import (
	"backend/internal/handlers"
	middleware "backend/internal/middlewares"
	"backend/internal/mux"
	"net/http"
)

func Register(mux *mux.Mux, m *middleware.Middleware, h *handlers.Handler) {
	mux.Handle("POST /sessions", m.CheckContentType(http.HandlerFunc(h.Login)))
	mux.Handle("DELETE /sessions", m.Auth(http.HandlerFunc(h.Logout)))
	mux.Handle("GET /users", m.Auth(http.HandlerFunc(h.GetLoggedInUser)))
	mux.Handle("POST /users", m.CheckContentType(http.HandlerFunc(h.CreateUser)))
	mux.Handle("DELETE /users", m.Auth(http.HandlerFunc(h.DeleteUser)))
	mux.Handle("GET /notes", m.Auth(http.HandlerFunc(h.FindUserNotes)))
	mux.Handle("GET /notes/{id}", m.Auth(http.HandlerFunc(h.FindNoteById)))
	mux.Handle("POST /notes", m.Auth(m.CheckContentType(http.HandlerFunc(h.CreateNote))))
	mux.Handle("DELETE /notes/{id}", m.Auth(http.HandlerFunc(h.DeleteNote)))
	mux.Handle("PUT /notes/{id}", m.Auth(m.CheckContentType(http.HandlerFunc(h.UpdateNote))))
}
