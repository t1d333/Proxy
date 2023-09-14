package proxy

import (
	"crypto/rand"
	"crypto/tls"
	"io"
	"math/big"
	"net"
	"net/http"
	"os/exec"
	"time"

	"go.uber.org/zap"
)

var HopHeaders = []string{
	"Connection",
	"Keep-Alive",
	"Proxy-Connection",
	"Proxy-Authenticate",
	"Proxy-Authorization",
	"TE",
	"Trailer",
	"Transfer-Encoding",
	"Upgrade",
}

type ForwardProxy struct{ logger *zap.SugaredLogger }

func NewForwardProxy(logger *zap.SugaredLogger) *ForwardProxy {
	return &ForwardProxy{logger}
}

func (p *ForwardProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.logger.Infow("new request", "uri", r.RequestURI)
	if r.Method == http.MethodConnect {
		p.handleHttps(w, r)
	} else {
		p.handleHttp(w, r)
	}
}

func (p *ForwardProxy) handleHttp(w http.ResponseWriter, r *http.Request) {
	client := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	r.RequestURI = ""

	for _, header := range HopHeaders {
		r.Header.Del(header)
	}
	response, err := client.Do(r)
	if err != nil {
		p.logger.Error("failed to send request", zap.Error(err))
		return
	}

	defer response.Body.Close()

	for k, vv := range response.Header {
		for _, v := range vv {
			w.Header().Add(k, v)
		}
	}

	w.WriteHeader(response.StatusCode)
	io.Copy(w, response.Body)
}

func (p *ForwardProxy) handleHttps(w http.ResponseWriter, r *http.Request) {
	hj, ok := w.(http.Hijacker)
	if !ok {
		p.logger.Error("failed to convert hijacker")
	}

	conn, _, err := hj.Hijack()
	if err != nil {
		p.logger.Error("failed to hijack conn", zap.Error(err))
	}

	host, _, err := net.SplitHostPort(r.Host)
	if err != nil {
		p.logger.Error("failed to splitting host/port", zap.Error(err))
	}

	p.CreateCert(host)

	if _, err := conn.Write([]byte("HTTP/1.0 200 Connection established")); err != nil {
		p.logger.Error("failed to writing status", zap.Error(err))
	}
}

func (p *ForwardProxy) CreateCert(host string) (tls.Certificate, error) {
	serial, _ := rand.Int(rand.Reader, big.NewInt(time.Now().Unix()))
	cmd := exec.Command("/bin/sh", "/scripts/gen_cert.sh", host, serial.String())

	if b, err := cmd.CombinedOutput(); err != nil {
		p.logger.Error("failed to generate new cert", zap.Error(err), zap.String("output", string(b)))
	}

	return tls.Certificate{}, nil
}
