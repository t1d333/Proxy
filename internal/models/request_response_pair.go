package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type RequestResponsePair struct {
	ID       primitive.ObjectID `bson:"_id" json:"id"`
	Request  Request            `bson:"request" json:"request"`
	Response Response           `bson:"response" json:"response"`
}
