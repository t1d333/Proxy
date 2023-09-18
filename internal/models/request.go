package models

import (
	"bytes"
	"fmt"
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
	Scheme     string              `bson:"scheme" json:"scheme"`
}

func (r *Request) ParseRequest(req *http.Request) {
	r.Method = req.Method
	r.Path = req.URL.Path
	r.Host = req.Host
	r.Scheme = req.URL.Scheme
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

func (r *Request) ConvertToHttpRequest() *http.Request {
	body := []byte{}
	url := ""

	url += r.Scheme + "://" + r.Host + r.Path

	if len(r.GetParams) > 0 {
		url += "?"
		for k, v := range r.GetParams {
			for _, i := range v {
				url += fmt.Sprintf("%s=%s&", k, i)
			}
		}
	}

	if len(r.PostParams) > 0 {
		for k, v := range r.PostParams {
			for _, i := range v {
				body = append(body, []byte(k+"="+i+"&")...)
			}
		}
	} else {
		body = []byte(r.Body)
	}

	res, err := http.NewRequest(r.Method, url, bytes.NewReader(body))
	if err != nil {
		return res
	}

	for k, v := range r.Headers {
		res.Header.Add(k, v)
	}

	for k, v := range r.Cookies {
		res.AddCookie(&http.Cookie{Name: k, Value: v})
	}

	return res
}
