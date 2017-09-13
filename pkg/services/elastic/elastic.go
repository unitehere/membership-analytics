package elastic

import (
	"context"
	"gopkg.in/olivere/elastic.v5"
	"sync"

	"membership-analytics/config"
)

var (
	clientInit sync.Once
	client     *elastic.Client
)

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

// SearchSSN takes in a ssn as a string and returns an *elastic.SearchResult or error
func SearchSSN(ssn string) (*elastic.SearchResult, error) {
	ctx := context.Background()

	query := elastic.NewMatchQuery("demographics.ssn", ssn).Fuzziness("Auto")

	searchResult, err := Client().Search().
		Index(config.Values.Index).
		Query(query).
		Pretty(true).
		FetchSourceContext(elastic.NewFetchSourceContext(true).Include("imis_id", "demographics.ssn")).
		Do(ctx)

	return searchResult, err
}
