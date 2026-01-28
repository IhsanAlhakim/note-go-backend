package routes

import (
	"backend/internal/handlers"
	middleware "backend/internal/middlewares"
	"backend/internal/mux"
	"net/http"
)

func Register(mux *mux.Mux, m *middleware.Middleware, h *handlers.Handler) {
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Connection Ok"))
	})

	mux.Handle("/login", m.CheckContentType(http.HandlerFunc(h.Login)))
	mux.Handle("/logout", m.Auth(http.HandlerFunc(h.Logout)))
	mux.Handle("/user", m.Auth(http.HandlerFunc(h.GetLoggedInUser)))
	mux.Handle("/create/user", m.CheckContentType(http.HandlerFunc(h.CreateUser)))
	mux.Handle("/delete/user", m.Auth(http.HandlerFunc(h.DeleteUser)))
	mux.Handle("/notes", m.Auth(http.HandlerFunc(h.FindUserNotes)))
	mux.Handle("/note", m.Auth(http.HandlerFunc(h.FindNoteById)))
	mux.Handle("/create/note", m.Auth(m.CheckContentType(http.HandlerFunc(h.CreateNote))))
	mux.Handle("/delete/note", m.Auth(http.HandlerFunc(h.DeleteNote)))
	mux.Handle("/update/note", m.Auth(m.CheckContentType(http.HandlerFunc(h.UpdateNote))))
}
