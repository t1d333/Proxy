package main

import (
	"net/http"

	"github.com/t1d333/proxyhw/internal/proxy"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewDevelopment()
	sugar := logger.Sugar()
	defer logger.Sync()
	proxy := proxy.NewForwardProxy(sugar)

	http.ListenAndServe(":8080", proxy)
}
