package proxy

import (
	"bufio"
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"io"
	"math/big"
	"net"
	"net/http"
	"net/url"
	"os/exec"
	"strings"
	"time"

	"go.uber.org/zap"
)

var HopHeaders = []string{
	// "Connection",
	"Keep-Alive",
	"Proxy-Connection",
	"Proxy-Authenticate",
	"Proxy-Authorization",
	"TE",
	"Trailer",
	"Transfer-Encoding",
	"Upgrade",
}

type ForwardProxy struct {
	logger *zap.SugaredLogger
}

func NewForwardProxy(logger *zap.SugaredLogger) *ForwardProxy {
	return &ForwardProxy{logger}
}

func (p *ForwardProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.logger.Infow("new request", "url", r.URL.String(), "method", r.Method, "scheme", r.URL.Scheme)
	if r.Method == http.MethodConnect {
		p.handleHttps(w, r)
	} else {
		p.handleHttp(w, r)
	}
}

func (p *ForwardProxy) handleHttp(w http.ResponseWriter, r *http.Request) {
	r.RequestURI = ""

	for _, header := range HopHeaders {
		r.Header.Del(header)
	}

	response, err := http.DefaultTransport.RoundTrip(r)
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
	if _, err := io.Copy(w, response.Body); err != nil {
		p.logger.Error("failed to copy data from response", zap.Error(err))
	}
}

func (p *ForwardProxy) handleHttps(w http.ResponseWriter, rawReq *http.Request) {
	hj, ok := w.(http.Hijacker)
	if !ok {
		p.logger.Error("failed to convert hijacker")
		return
	}

	conn, _, err := hj.Hijack()
	
	defer conn.Close()
	if err != nil {
		p.logger.Error("failed to hijack conn", zap.Error(err))
		return
	}

	host, _, err := net.SplitHostPort(rawReq.Host)
	if err != nil {
		p.logger.Error("failed to splitting host/port", zap.Error(err))
		return
	}

	cert, err := p.CreateCert(host)
	if err != nil {
		p.logger.Error("failed go generate certificate", zap.Error(err))
		return
	}

	if _, err := conn.Write([]byte("HTTP/1.0 200 Connection established\r\n\r\n")); err != nil {
		p.logger.Error("failed to writing status", zap.Error(err))
		return
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS13,
	}

	tlsConn := tls.Server(conn, tlsConfig)

	reader := bufio.NewReader(tlsConn)

	var r *http.Request = nil

	for r, err = http.ReadRequest(reader); err != io.EOF; r, err = http.ReadRequest(reader) {
		if err != nil {
			p.logger.Error("failed to read from connection", zap.Error(err))
			return
		}

		p.UpdateURL(r, rawReq.Host)
		p.logger.Infow("tls request", "uri", r.RequestURI, "url", r.URL.String(), "method", r.Method)

		response, err := http.DefaultClient.Do(r)
		if err != nil {
			p.logger.Error("failed to send request", zap.Error(err))
		}

		defer response.Body.Close()

		if err := response.Write(tlsConn); err != nil {
			p.logger.Error("failed to write response to tls connection", zap.Error(err))
		}
	}
	tlsConn.Close()
}

func (p *ForwardProxy) UpdateURL(r *http.Request, host string) {
	if !strings.HasPrefix(r.URL.String(), "https") {
		host = "https://" + host
	}

	newUrl, err := url.Parse(host)
	if err != nil {
		p.logger.Error("failed to update request url", zap.Error(err))
	}

	newUrl.Path = r.URL.Path
	newUrl.RawQuery = r.URL.RawQuery

	r.URL = newUrl
	r.RequestURI = ""
}

func (p *ForwardProxy) CreateCert(host string) (tls.Certificate, error) {
	serial, _ := rand.Int(rand.Reader, big.NewInt(time.Now().Unix()))
	cmd := exec.Command("/bin/sh", "/scripts/gen_cert.sh", host, serial.String())

	if b, err := cmd.CombinedOutput(); err != nil {
		p.logger.Error("failed to generate new cert", zap.Error(err), zap.String("output", string(b)))
	}

	return tls.LoadX509KeyPair(fmt.Sprintf("/certs/%s.crt", host), "/cert.key")
}
