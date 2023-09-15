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

	sugar.Info("starting proxy server on port 8080...")
	if err := http.ListenAndServe(":8080", proxy); err != nil {
		sugar.Fatalw("failed to start server", "err", err)
	}
}
