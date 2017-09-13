package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"gopkg.in/olivere/elastic.v5"
	"net/http"
	"sync"

	"membership-analytics/config"
)

var (
	clientInit sync.Once
	client     *elastic.Client
)

// The ResponseValues type describes the structure of the all responses.
type ResponseValues struct {
	Values []map[string]interface{}
}

// Client inits a new client on initial call, and returns the initialized client subsequently
func Client() *elastic.Client {
	clientInit.Do(func() {
		c, err := elastic.NewClient(
			elastic.SetURL("https://elasticsearch.unitehere.org:9200"),
			elastic.SetBasicAuth(config.Values.ElasticUsername, config.Values.ElasticPassword),
			elastic.SetSniff(false),
			elastic.SetHealthcheck(false))
		if err != nil {
			panic(err)
		}
		client = c
	})
	return client
}

// GetSearchSSN returns a fuzzy matched array of imis_id given a ssn
// r.Get("/ssn", handlers.GetSearchSSN)
func GetSearchSSN(w http.ResponseWriter, r *http.Request) {
	ssnQuery := r.URL.Query()["q"]

	if len(ssnQuery) == 0 || len(ssnQuery[0]) < 7 {
		writeError(w, http.StatusBadRequest, errors.New("You need to pass in a ssn string of at least 7 digits as a q param"))
		return
	}
	ctx := context.Background()

	query := elastic.NewMatchQuery("demographics.ssn", ssnQuery[0]).Fuzziness("Auto")

	searchResult, err := Client().Search().
		Index(config.Values.Index).
		Query(query).
		Pretty(true).
		FetchSourceContext(elastic.NewFetchSourceContext(true).Include("imis_id", "demographics.ssn")).
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
				panic("Could not read data from api response")
			}
			payload.Values = append(payload.Values, data)
		}
	} else {
		writeError(w, http.StatusNotFound, errors.New("Did not match any ssns"))
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
