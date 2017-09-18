package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/unitehere/membership-analytics/pkg/handlers"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/goware/cors"
	"github.com/unrolled/secure"
)

func main() {
	secureMiddleware := secure.New(secure.Options{
		FrameDeny:        true,
		BrowserXssFilter: true,
	})

	r := chi.NewRouter()

	cors := cors.New(cors.Options{
		// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	r.Use(cors.Handler)

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(secureMiddleware.Handler)

	r.Route("/search", func(r chi.Router) {
		r.Get("/ssn", handlers.GetSearchSSN)
		r.Post("/ssn", handlers.PostSearchSSN)

		r.Post("/name", handlers.PostSearchName)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}

	fmt.Println("Application initializing on port " + port)
	http.ListenAndServe(":"+port, r)
}
