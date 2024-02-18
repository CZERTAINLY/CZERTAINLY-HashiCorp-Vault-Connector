package main

import (
	"CZERTAINLY-HashiCorp-Vault-Connector/cmd/config"
	"CZERTAINLY-HashiCorp-Vault-Connector/cmd/info"
	"CZERTAINLY-HashiCorp-Vault-Connector/cmd/logger"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"log"
	"net/http"
)

var version = "0.0.1"

func main() {
	l := logger.Get()
	c := config.Get()

	l.Info("Starting the server version: " + version)

	infoService := info.NewService()

	// start the server
	router := mux.NewRouter()
	router.Use(logMiddleware)

	info.RegisterRoutes(router, infoService)

	log.Fatal(http.ListenAndServe(":"+c.Server.Port, router))
}

func logMiddleware(next http.Handler) http.Handler {
	l := logger.Get()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l.Info("Request received", zap.String("path", r.URL.Path))
		next.ServeHTTP(w, r)
	})
}
