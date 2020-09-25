package application

import (
	"context"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type App struct {
	log *zap.SugaredLogger
	srv *http.Server
}

func NewApp(addr string, handler http.Handler) *App {
	return &App{
		log: zap.S(),
		srv: &http.Server{Addr: addr, Handler: handler},
	}
}

func (a *App) Startup() {
	go func() {
		a.log.Info("Starting Http Server...")
		if err := a.srv.ListenAndServe(); err != http.ErrServerClosed {
			a.log.Panicw("Could not start HTTP Server!", "err", err)
		} else {
			a.log.Info("HTTP server shutdown successfully")
		}
	}()
}

func (a *App) Close() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	a.log.Info("Shutting down HTTP server...")
	err := a.srv.Shutdown(ctx)
	if err != nil {
		a.log.Errorw("Graceful shutdown of server failed", "err", err)
	}
}
