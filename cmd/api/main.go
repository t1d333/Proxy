package main

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/t1d333/proxyhw/internal/db/mongo"
	del "github.com/t1d333/proxyhw/internal/delivery/http"
	rep "github.com/t1d333/proxyhw/internal/repository/mongo"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewDevelopment()
	sugar := logger.Sugar()
	defer logger.Sync()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	conn := mongo.InitDB(ctx, sugar)

	rep := rep.NewRepository(conn, sugar)
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	del.InitHandlers(router, sugar, rep)
	sugar.Info("starting proxy server on port 8000...")
	if err := http.ListenAndServe(":8000", router); err != nil {
		sugar.Fatalw("failed to start server", "err", err)
	}
}
