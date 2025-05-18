package middleware

import (
	"backend/data"
	"backend/utils"
	"net/http"
)


func (m *Middleware) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := m.store.Get(r, data.SESSION_ID)

	if session.Values["userID"] == nil {
		utils.JSONResponse(w, R{Message: "User not authenticated"}, http.StatusUnauthorized)
		return
	}

	next.ServeHTTP(w,r)
	})
}