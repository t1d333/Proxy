package models

import (
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
	PostParams map[string]string   `bson:"post_params" json:"post_params"`
}

func (r *Request) ParseRequst(req *http.Request) {
	r.Method = req.Method
	r.Path = req.URL.Path
	r.Host = req.Host
	r.Cookies = map[string]string{}
	r.Headers = map[string]string{}
	r.PostParams = map[string]string{}
	r.GetParams = map[string][]string{}

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
}
