package middleware

import (
	"backend/data"
	"backend/utils"
	"fmt"
	"net/http"
)


func (m *Middleware) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := m.store.Get(r, data.SESSION_ID)

	cookie, _ := r.Cookie(data.SESSION_ID)

	fmt.Println(cookie)

	fmt.Println(session.Values["userID"])

	if session.Values["userID"] == nil {
		utils.JSONResponse(w, R{Message: "User not authenticated"}, http.StatusUnauthorized)
		return
	}

	next.ServeHTTP(w,r)
	})
}