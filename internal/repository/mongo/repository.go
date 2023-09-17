package mongo

import (
	"context"
	"fmt"
	"net/http"

	"github.com/t1d333/proxyhw/internal/models"
	repModule "github.com/t1d333/proxyhw/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type repository struct {
	conn   *mongo.Client
	logger *zap.SugaredLogger
}

func (rep *repository) CreateRequestResponsePair(req *http.Request, res *http.Response) error {
	if req == nil {
		return fmt.Errorf("request is nil")
	}

	if res == nil {
		return fmt.Errorf("response is nil")
	}

	reqModel := models.Request{}
	resModel := models.Response{}

	reqModel.ParseRequst(req)
	resModel.ParseResponse(res)

	pair := models.RequestResponsePair{
		ID:       primitive.NewObjectID(),
		Request:  reqModel,
		Response: resModel,
	}

	if _, err := rep.conn.Database("proxy").Collection("requests").InsertOne(context.Background(), pair); err != nil {
		rep.logger.Error("failed to insert new request", zap.Error(err))
		return fmt.Errorf("failed to insert new request: %w", err)
	}

	return nil
}

func (rep *repository) GetRequestResponsePair(id string) (models.RequestResponsePair, error) {
	panic("unimplemented")
}

func (rep *repository) GetAllPairs() ([]models.RequestResponsePair, error) {
	res := []models.RequestResponsePair{}
	ctx := context.Background()

	c, err := rep.conn.Database("proxy").Collection("requests").Find(ctx, struct{}{})
	if err != nil {
		rep.logger.Error("failed to get all pairs", zap.Error(err))
		return res, fmt.Errorf("failed to get all pairs: %w", err)
	}

	if err := c.All(ctx, &res); err != nil {
		rep.logger.Error("failed to get all pairs", zap.Error(err))
		return res, fmt.Errorf("failed to get all pairs: %w", err)
	}

	return res, nil
}

func NewRepository(conn *mongo.Client, logger *zap.SugaredLogger) repModule.Repository {
	return &repository{conn, logger}
}
