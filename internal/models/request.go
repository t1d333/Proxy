package models

import (
	"net/http"
	"net/url"
)

type Request struct {
	Method     string              `bson:"method" json:"method"`
	Path       string              `bson:"path" json:"path"`
	Headers    map[string][]string `bson:"headers" json:"headers"`
	GetParams  map[string][]string `bson:"get_params" json:"get_params"`
	Cookies    map[string]string   `bson:"cookies" json:"cookies"`
	PostParams map[string]string   `bson:"post_params" json:"post_params"`
}

func (r *Request) ParseRequst(req *http.Request) {
	r.Method = req.Method
	r.Path = req.URL.Path
	r.Headers = req.Header
	for _, c := range req.Cookies() {
		r.Cookies[c.Name] = c.Value
	}

	tmp, _ := url.ParseQuery(req.URL.RawQuery)
	for k, v := range tmp {
		r.GetParams[k] = v
	}

	// req.ParseForm
	// TODO: разобраться с post параметрами
}
