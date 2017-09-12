package HttpHandlers

import (
	"context"
	"encoding/json"
	"gopkg.in/olivere/elastic.v5"
	"net/http"

	"membership-analytics/pkg/configuration"
)

// The ResponseValues type describes the structure of the all responses.
type ResponseValues struct {
	Values map[string]interface{}
}

// GetSearchSSN ...[comment here]
func GetSearchSSN(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	client, err := elastic.NewSimpleClient(
		elastic.SetURL("https://elasticsearch.unitehere.org:9200"),
		elastic.SetBasicAuth(config.Values.ElasticUsername, config.Values.ElasticPassword))
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
