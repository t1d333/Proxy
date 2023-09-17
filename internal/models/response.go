package models

import (
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Response struct {
	ID      primitive.ObjectID  `bson:"_id" json:"id"`
	Code    int                 `bson:"code" json:"code"`
	Message string              `bson:"message" json:"message"`
	Headers map[string][]string `bson:"headers" json:"headers"`
	Body    string              `bson:"body" json:"body"`
}

func (r *Response) method(res *http.Response) {
	r.Code = res.StatusCode
	// r.Message = ???
	r.Headers = res.Header

	body := []byte{}
	// TODO: доавить обработку ошибки
	res.Body.Read(body)
	r.Body = string(body)
}
