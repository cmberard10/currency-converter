package server

import (
	"fmt"
	"stylight/internal/services"
	"net/http"
	"log"
	"stylight/internal/config"
)

//Server is a server ;-)
type Server struct {
	Configuration *config.WebServerConfig
	Router        *Router
}

//NewServer creates a new server
func NewServer(config *config.WebServerConfig) *Server {
	server := &Server{
		Configuration: config,
		Router:        NewRouter(),
	}

	return server
}

//RunServer initializes the server
func RunServer() (err error) {
	webServerConfig, err := config.FromEnv()
	if err != nil {
		return err
	}

	log.Printf("Starting HTTPS server on port %s", webServerConfig.Port)

	err = services.Initialize(webServerConfig.Service)
	if err != nil {
		log.Printf("an error occurred while initializing services: %s", err.Error())
		return err
	}

	server := NewServer(webServerConfig)
	server.Router.InitializeRoutes(webServerConfig)

	if err := http.ListenAndServe(
		fmt.Sprintf("%v:%v", "", webServerConfig.Port),
		*server.Router,
	); err != nil {
		panic(err)
	}

	return nil
}
