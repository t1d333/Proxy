package http

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"

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

	if _, err := w.Write(b); err != nil {
		d.logger.Error("failed to write reponse", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (d *delivery) repeatRequest(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	pair, err := d.rep.GetRequestResponsePair(id)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, mongo.ErrNotFound) {
			status = http.StatusNotFound
		}
		http.Error(w, err.Error(), status)
		return
	}
	req := pair.Request.ConvertToHTTPRequest()

	proxyURL, _ := url.Parse("http://proxy:8080")

	client := http.Client{
		Transport: &http.Transport{Proxy: http.ProxyURL(proxyURL)},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	response, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer response.Body.Close()

	for k, vv := range response.Header {
		for _, v := range vv {
			w.Header().Add(k, v)
		}
	}

	w.WriteHeader(response.StatusCode)
	if _, err := io.Copy(w, response.Body); err != nil {
		d.logger.Error("failed to copy data from response", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (d *delivery) getRequest(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	pair, err := d.rep.GetRequestResponsePair(id)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, mongo.ErrNotFound) {
			status = http.StatusNotFound
		}

		http.Error(w, err.Error(), status)
		return
	}

	b, _ := json.Marshal(pair)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(b); err != nil {
		d.logger.Error("failed to write reponse", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
	}
}
