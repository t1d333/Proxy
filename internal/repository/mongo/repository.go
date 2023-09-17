package mongo

import "go.mongodb.org/mongo-driver/mongo"

type repository struct {
	conn *mongo.Client
}


func (r *repository) CreateRequest()  {
	
}
