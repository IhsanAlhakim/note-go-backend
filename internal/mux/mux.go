package mux

import "net/http"

type Mux struct {
	http.ServeMux
	middlewares []func(next http.Handler) http.Handler
}

func (m *Mux) RegisterMiddleware(next func(next http.Handler) http.Handler) {
	m.middlewares = append(m.middlewares, next)
}

func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var current http.Handler = &m.ServeMux

	for _, next := range m.middlewares {
		current = next(current)
	}

	current.ServeHTTP(w, r)
}

func New() *Mux {
	return &Mux{}
}
