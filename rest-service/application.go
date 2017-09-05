package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/goware/cors"
	"github.com/unrolled/secure"
	"gopkg.in/olivere/elastic.v5"
)

// The ResponseValues type describes the structure of the all responses.
type ResponseValues struct {
	Values map[string]interface{}
}

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
		r.Get("/ssn", GetSearchSSN)
		// r.Post("/ssn", PostSearchSSN)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "4444"
	}

	fmt.Println("Application initializing on port " + port)
	http.ListenAndServe(":"+port, r)
}

// GetSearchSSN ...[comment here]
func GetSearchSSN(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	client, err := elastic.NewClient(
		elastic.SetURL("https://elasticsearch.unitehere.org:9200"),
		elastic.SetSniff(false),
	)
	if err != nil {
		// Handle error
		panic(err)
	}
	query := elastic.NewMatchQuery("demographics.ssn", r.URL.Query()["q"][0])

	searchResult, err := client.Search().
		Query(query).
		Pretty(true).
		FetchSourceContext(elastic.NewFetchSourceContext(true).Include("imis_id")).
		Do(ctx)

	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	payload := ResponseValues{}
	if searchResult.Hits.TotalHits > 0 {
		for _, hit := range searchResult.Hits.Hits {
			var data map[string]interface{}
			err := json.Unmarshal(*hit.Source, &data)
			if err != nil {
				// Deserialization failed
			}
			payload.Values = data
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(payload)
	return
}

func writeError(w http.ResponseWriter, status int, err error) {
	w.WriteHeader(status)
	w.Write([]byte(err.Error()))
}
