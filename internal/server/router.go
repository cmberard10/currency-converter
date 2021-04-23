package server

import (
	"github.com/gorilla/mux"
	"stylight/internal/config"
	"stylight/internal/server/handlers"
	"net/http"
)

//Router is the core router struct
type Router struct {
	*mux.Router
}

//NewRouter created a new router
func NewRouter() *Router {
	return &Router{mux.NewRouter()}
}

//InitializeRoutes does exactly what you think it does
func (r *Router) InitializeRoutes(serverConfig *config.WebServerConfig) {
	stylightHandler := r.Router.PathPrefix("/stylight").Subrouter()

	stylightHandler.HandleFunc("/healthcheck", handlers.HealthCheckHandler).
		Methods(http.MethodGet).
		Name("healthcheck")

	stylightHandler.HandleFunc("/conversion", handlers.CurrencyConversionHandler).
		Methods(http.MethodPost).
		Name("conversion")
}
