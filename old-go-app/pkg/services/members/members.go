package members

import (
	"context"
	"encoding/json"
	"errors"
	"sync"

	"gopkg.in/olivere/elastic.v5"

	"github.com/asaskevich/govalidator"
	"github.com/unitehere/membership-analytics/config"
)

var (
	clientInit     sync.Once
	errInvalidName = errors.New("You need at least a first or last name")
	errInvalidSSN  = errors.New("You need to pass in a ssn string of at least 7 digits")
)

// Service interface for all simple member searches
type Service interface {
	SearchSSN(ssnQuery SSNQuery, from int, size int) (Member, error)
	SearchName(query NameQuery, from int, size int) (Member, error)
}

type service struct {
	client *elastic.Client
}

// Member represents all responses returned by the service
type Member struct {
	Data      []map[string]interface{} `json:"data"`
	TotalHits int64                    `json:"total_hits"`
}

// NameQuery is used in SearchName
type NameQuery struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	FromValue string `json:"from"`
	SizeValue string `json:"size"`
}

// SSNQuery is used in SearchSSN
type SSNQuery struct {
	SSN       string `json:"ssn"`
	FromValue string `json:"from"`
	SizeValue string `json:"size"`
}

// Validate validates that either firstName or lastName exists
func (t NameQuery) Validate() error {
	if govalidator.IsNull(t.FirstName) && govalidator.IsNull(t.LastName) {
		return errInvalidName
	}
	return nil
}

// From returns the from value of the query
func (t NameQuery) From() string {
	return t.FromValue
}

// Size returns the size value of the query
func (t NameQuery) Size() string {
	return t.SizeValue
}

// Validate validates that an ssn of atleast length of 7 exists
func (t SSNQuery) Validate() error {
	if govalidator.IsNull(t.SSN) || govalidator.StringLength(t.SSN, "0", "7") {
		return errInvalidSSN
	}
	return nil
}

// From returns the from value of the query
func (t SSNQuery) From() string {
	return t.FromValue
}

// Size returns the size value of the query
func (t SSNQuery) Size() string {
	return t.SizeValue
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
func (s *service) SearchSSN(ssnQuery SSNQuery, from int, size int) (Member, error) {
	ctx := context.Background()

	query := elastic.NewBoolQuery()
	query = query.Must(elastic.NewMatchQuery("demographics.ssn", ssnQuery.SSN).Fuzziness("1"))
	query = query.Should(elastic.NewTermQuery("demographics.ssn.raw", ssnQuery.SSN).Boost(10.0))

	searchResult, err := s.client.Search().
		Index(config.Values.Index).
		From(from).
		Size(size).
		Query(query).
		Pretty(true).
		FetchSourceContext(elastic.NewFetchSourceContext(true).Include("imis_id")).
		Do(ctx)

	if err != nil {
		return Member{}, err
	}

	resultLength := len(searchResult.Hits.Hits)
	result := make([]map[string]interface{}, resultLength, resultLength)
	for i, hit := range searchResult.Hits.Hits {
		var data map[string]interface{}
		err := json.Unmarshal(*hit.Source, &data)
		if err != nil {
			return Member{}, err
		}
		result[i] = data
	}

	member := Member{Data: result, TotalHits: searchResult.Hits.TotalHits}

	return member, err
}

// SearchName takes in a ssn as a string and returns an *elastic.SearchResult or error
func (s *service) SearchName(query NameQuery, from int, size int) (Member, error) {
	ctx := context.Background()

	elasticQuery := elastic.NewBoolQuery()
	mapping := map[string]string{
		"first_name": query.FirstName,
		"last_name":  query.LastName,
	}

	for key, value := range mapping {
		if value != "" {
			elasticQuery = elasticQuery.Must(elastic.NewMatchQuery(key, value).Fuzziness("AUTO"))
		}
	}

	searchResult, err := s.client.Search().
		Index(config.Values.Index).
		From(from).
		Size(size).
		Query(elasticQuery).
		FetchSourceContext(elastic.NewFetchSourceContext(true).Include("imis_id", "first_name", "last_name")).
		Do(ctx)

	if err != nil {
		return Member{}, err
	}

	resultLength := len(searchResult.Hits.Hits)
	result := make([]map[string]interface{}, resultLength, resultLength)
	for i, hit := range searchResult.Hits.Hits {
		var data map[string]interface{}
		err := json.Unmarshal(*hit.Source, &data)
		if err != nil {
			return Member{}, err
		}
		result[i] = data
	}
	member := Member{Data: result, TotalHits: searchResult.Hits.TotalHits}

	return member, err
}
