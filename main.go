package main

import (
	"os"
	"os/signal"

	"github.com/landonpoch/lib-api/application"
	"github.com/landonpoch/lib-api/data"
	"go.uber.org/zap"
)

func main() {
	// Structured logging provides indexing for microservices
	// Usually through something like logging pipeline with ELK
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	zap.ReplaceGlobals(logger)

	log := zap.S()
	log.Info("Application lib-api starting up...")

	// Config usually read in from external config service like consul
	// Currently using port 8080 with standard HTTP
	// HTTPS termination and authentication is usually handled
	// with an API gateway of some sort that sits in front of
	// the microservice.
	cfg := Config{ServerAddr: ":8080"}
	app := bootstrap(cfg)
	app.Startup()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	<-c
	log.Info("Interrupt received. Application lib-api shutting down...")

	app.Close() // Graceful shutdown of connections w/ timeout
}

func bootstrap(cfg Config) *application.App {
	// Wire up dependencies
	// In memory repo could be replaced with something more concrete like
	// cassandra repo assuming it implements the same repo interface.
	repo := data.NewInMemBookRepository()
	routes := application.NewRoutes(repo)
	return application.NewApp(cfg.ServerAddr, routes.Router)
}

type Config struct {
	ServerAddr string // May require serialization tags if using service to populate
}
