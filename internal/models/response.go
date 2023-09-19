package models

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
)

type Response struct {
	Code    int               `bson:"code" json:"code"`
	Message string            `bson:"message" json:"message"`
	Headers map[string]string `bson:"headers" json:"headers"`
	Body    string            `bson:"body" json:"body"`
}

func (r *Response) ParseResponse(res *http.Response) {
	r.Code = res.StatusCode
	r.Message = http.StatusText(r.Code)
	r.Headers = map[string]string{}

	for k, v := range res.Header {
		r.Headers[k] = v[0]
	}

	buffer := bytes.NewBuffer([]byte{})
	io.Copy(buffer, res.Body)
	res.Body = io.NopCloser(bytes.NewReader(buffer.Bytes()))

	if res.Header.Get("Content-Encoding") == "gzip" {
		reader, err := gzip.NewReader(buffer)
		if err != nil {
			return
		}
		defer reader.Close()

		b, err := io.ReadAll(reader)
		if err != nil {
			return
		}
		r.Body = string(b)
	} else {
		r.Body = buffer.String()
	}
}
