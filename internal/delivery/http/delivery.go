package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type delivery struct {
	logger *zap.Logger
}

func (d *delivery) initHandlers(router chi.Router) {
	router.Get("/requests", func(w http.ResponseWriter, r *http.Request) {})
	router.Get("/requests/:id", func(w http.ResponseWriter, r *http.Request) {})
	router.Post("/repeat/:id", func(w http.ResponseWriter, r *http.Request) {})
}
