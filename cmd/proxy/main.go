package main

import (
	"context"
	"net/http"
	"time"

	"github.com/t1d333/proxyhw/internal/db/mongo"
	"github.com/t1d333/proxyhw/internal/proxy"
	rep "github.com/t1d333/proxyhw/internal/repository/mongo"
	"go.uber.org/zap"
)

const timeout = 10 * time.Second

func main() {
	logger, _ := zap.NewDevelopment()
	sugar := logger.Sugar()
	defer logger.Sync()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	conn := mongo.InitDB(ctx, sugar)

	rep := rep.NewRepository(conn, sugar)
	proxy := proxy.NewForwardProxy(sugar, rep)
	sugar.Info("starting proxy server on port 8080...")
	if err := http.ListenAndServe(":8080", proxy); err != nil {
		sugar.Fatalw("failed to start server", "err", err)
	}
}
