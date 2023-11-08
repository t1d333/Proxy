package mongo

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/t1d333/proxyhw/internal/models"
	repModule "github.com/t1d333/proxyhw/internal/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

var ErrNotFound = errors.New("request not found")

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

	if err := reqModel.ParseRequest(req); err != nil {
		rep.logger.Error("failed to parse request", err)
		return fmt.Errorf("failed to parse request: %w", err)
	}

	if err := resModel.ParseResponse(res); err != nil {
		rep.logger.Error("failed to parse response", err)
		return fmt.Errorf("failed to parse reponse: %w", err)
	}

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

func (rep *repository) GetRequestResponsePair(hid string) (models.RequestResponsePair, error) {
	ctx := context.Background()

	res := models.RequestResponsePair{}
	id, _ := primitive.ObjectIDFromHex(hid)

	c := rep.conn.Database("proxy").Collection("requests").FindOne(ctx, bson.D{bson.E{Key: "_id", Value: id}})

	if err := c.Decode(&res); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return res, ErrNotFound
		}
		rep.logger.Error("failed to find request response pair by id", zap.Error(err))

		return res, fmt.Errorf("failed to find request response pair by id: %w", err)
	}
	return res, nil
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
