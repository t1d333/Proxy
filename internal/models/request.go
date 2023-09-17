package models

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
)

type Request struct {
	Method     string              `bson:"method" json:"method"`
	Path       string              `bson:"path" json:"path"`
	Host       string              `bson:"host" json:"host"`
	Headers    map[string]string   `bson:"headers" json:"headers"`
	GetParams  map[string][]string `bson:"get_params" json:"get_params"`
	Cookies    map[string]string   `bson:"cookies" json:"cookies"`
	PostParams map[string][]string `bson:"post_params" json:"post_params"`
	Body       string              `bson:"body" json:"body"`
}

func (r *Request) ParseRequest(req *http.Request) {
	r.Method = req.Method
	r.Path = req.URL.Path
	r.Host = req.Host
	r.Cookies = map[string]string{}
	r.Headers = map[string]string{}
	r.PostParams = map[string][]string{}
	r.GetParams = map[string][]string{}
	r.Body = ""

	for k, v := range req.Header {
		if k == "Cookie" {
			continue
		}

		r.Headers[k] = v[0]

	}

	for _, c := range req.Cookies() {
		r.Cookies[c.Name] = c.Value
	}

	tmp, _ := url.ParseQuery(req.URL.RawQuery)

	for k, v := range tmp {
		r.GetParams[k] = v
	}

	if req.Body == nil {
		return
	}

	if req.Header.Get("Content-Type") == "application/x-www-form-urlencoded" {
		if err := req.ParseForm(); err != nil {
			return
		}

		for k, v := range req.PostForm {
			r.PostParams[k] = v
		}

		return
	}

	body := bytes.NewBuffer([]byte{})
	io.Copy(body, req.Body)
	r.Body = body.String()
}
