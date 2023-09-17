package http

import (
	"encoding/json"
	// "fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/t1d333/proxyhw/internal/repository"
	"github.com/t1d333/proxyhw/internal/repository/mongo"
	"go.uber.org/zap"
)

type delivery struct {
	logger *zap.SugaredLogger
	rep    repository.Repository
}

func InitHandlers(router chi.Router, logger *zap.SugaredLogger, rep repository.Repository) {
	d := &delivery{logger, rep}
	router.Get("/api/requests", d.getAllRequests)
	router.Get("/api/requests/{id}", d.getRequest)
	router.Post("/api/requests/{id}", d.repeatRequest)
}

func (d *delivery) getAllRequests(w http.ResponseWriter, r *http.Request) {
	pairs, err := d.rep.GetAllPairs()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	b, _ := json.Marshal(pairs)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func (d *delivery) repeatRequest(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	pair, err := d.rep.GetRequestResponsePair(id)
	if err != nil {
		status := http.StatusInternalServerError
		if err == mongo.ErrNotFound {
			status = http.StatusNotFound
		}
		http.Error(w, err.Error(), status)
	}

	
}

func (d *delivery) getRequest(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	pair, err := d.rep.GetRequestResponsePair(id)
	if err != nil {
		status := http.StatusInternalServerError
		if err == mongo.ErrNotFound {
			status = http.StatusNotFound
		}
		http.Error(w, err.Error(), status)
		return
	}

	b, _ := json.Marshal(pair)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}
