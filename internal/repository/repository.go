package repository

type Repository interface {
	CreateRequest() error
	CreateRespose() error
	GetRequest() error
	GetAllRequests() error
}
