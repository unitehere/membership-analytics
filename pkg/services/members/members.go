package members

import (
	"context"
	"encoding/json"
	"sync"

	"gopkg.in/olivere/elastic.v5"

	"github.com/unitehere/membership-analytics/config"
)

var (
	clientInit sync.Once
)

// Service interface for all simple member searches
type Service interface {
	SearchSSN(ssn string) ([]map[string]interface{}, error)
}

type service struct {
	client *elastic.Client
}

// Client inits a new client on initial call, and returns the initialized client subsequently
func Client() (Service, error) {
	var client *elastic.Client
	var err error
	clientInit.Do(func() {
		client, err = elastic.NewClient(
			elastic.SetURL(config.Values.ElasticURL),
			elastic.SetBasicAuth(config.Values.ElasticUsername, config.Values.ElasticPassword),
			elastic.SetSniff(false),
			elastic.SetHealthcheck(false))
	})
	if err != nil {
		return nil, err
	}
	return &service{client: client}, nil
}

// SearchSSN takes in a ssn as a string and returns an *elastic.SearchResult or error
func (s *service) SearchSSN(ssn string) ([]map[string]interface{}, error) {
	ctx := context.Background()

	query := elastic.NewMatchQuery("demographics.ssn", ssn).Fuzziness("Auto")

	searchResult, err := s.client.Search().
		Index(config.Values.Index).
		Query(query).
		Pretty(true).
		FetchSourceContext(elastic.NewFetchSourceContext(true).Include("imis_id")).
		Do(ctx)

	if err != nil {
		return nil, err
	}

	resultLength := len(searchResult.Hits.Hits)
	result := make([]map[string]interface{}, resultLength, resultLength)
	for i, hit := range searchResult.Hits.Hits {
		var data map[string]interface{}
		err := json.Unmarshal(*hit.Source, &data)
		if err != nil {
			return nil, err
		}
		result[i] = data
	}

	return result, err
}
