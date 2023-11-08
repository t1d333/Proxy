package models

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
)

type Response struct {
	Code    int               `bson:"code" json:"code"`
	Message string            `bson:"message" json:"message"`
	Headers map[string]string `bson:"headers" json:"headers"`
	Body    string            `bson:"body" json:"body"`
}

func (r *Response) ParseResponse(res *http.Response) error {
	r.Code = res.StatusCode
	r.Message = http.StatusText(r.Code)
	r.Headers = map[string]string{}

	for k, v := range res.Header {
		r.Headers[k] = v[0]
	}

	buffer := bytes.NewBuffer([]byte{})
	if _, err := io.Copy(buffer, res.Body); err != nil {
		return fmt.Errorf("failed to copy response body: %w", err)
	}

	res.Body = io.NopCloser(bytes.NewReader(buffer.Bytes()))

	if res.Header.Get("Content-Encoding") == "gzip" {
		reader, err := gzip.NewReader(buffer)
		if err != nil {
			return fmt.Errorf("failed to create gzip encoder: %w", err)
		}

		defer reader.Close()

		b, err := io.ReadAll(reader)
		if err != nil {
			return fmt.Errorf("failed to encode body: %w", err)
		}
		r.Body = string(b)
	} else {
		r.Body = buffer.String()
	}

	return nil
}
