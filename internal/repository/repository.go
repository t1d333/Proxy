package repository

import (
	"net/http"

	"github.com/t1d333/proxyhw/internal/models"
)

type Repository interface {
	CreateRequestResponsePair(req *http.Request, res *http.Response) error
	GetRequestResponsePair(id string) (models.RequestResponsePair, error)
	GetAllPairs() ([]models.RequestResponsePair, error)
}
