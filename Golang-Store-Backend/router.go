package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/Apouzi/golang-shop/app/api/routes"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

func (app *Config) GetAllowedOrigins() []string{
	allowedHostString := os.Getenv("ALLOWED_HOSTS")
	var AllowedOriginsFromEnv []string
	err := json.Unmarshal([]byte(allowedHostString),&AllowedOriginsFromEnv)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	fmt.Println("GetAllowedOrigins completed")
	fmt.Println("AllowedOriginsFromEnv:",AllowedOriginsFromEnv)
	return AllowedOriginsFromEnv
}
func (app *Config) StartRouter() http.Handler { // Change the receiver to (*Config)

	mux := chi.NewRouter()


	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   app.GetAllowedOrigins(),
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           301,
	}))
	// Test if this is working, maybe for microservice
	mux.Use(middleware.Heartbeat("/ping"))


	//Pass the mux to routes to use.
	routes.RouteDigest(mux, app.DB, app.Redis)
	
	return mux
}
