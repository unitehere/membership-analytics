package members

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	elastic "gopkg.in/olivere/elastic.v5"
)

func TestSearchSSNNoHits(t *testing.T) {
	handler := http.NotFound
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler(w, r)
	}))
	defer ts.Close()

	handler = func(w http.ResponseWriter, r *http.Request) {
		resp := `{
      "took": 0,
      "timed_out": false,
      "_shards": {
        "total": 5,
        "successful": 5,
        "failed": 0
      },
      "hits": {
        "total": 0,
        "max_score": null,
        "hits": []
      }
    }`

		w.Write([]byte(resp))
	}

	s, err := MockService(ts.URL)
	assert.NoError(t, err)

	actualResult, err := s.SearchSSN("123456789")
	assert.NoError(t, err)
	expectedResult := []map[string]interface{}{}
	assert.Equal(t, expectedResult, actualResult)
}

func TestSearchSSNOneHit(t *testing.T) {
	handler := http.NotFound
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler(w, r)
	}))
	defer ts.Close()

	handler = func(w http.ResponseWriter, r *http.Request) {
		resp := `{
      "took": 1,
      "timed_out": false,
      "_shards": {
        "total": 5,
        "successful": 5,
        "failed": 0
      },
      "hits": {
        "total": 1,
        "max_score": 7.411551,
        "hits": [
          {
            "_index": "members-2017.09.14",
            "_type": "member",
            "_id": "AV5_Z3USj5SI64ssYb5t",
            "_score": 7.411551,
            "_source": {
              "imis_id": "5962",
              "demographics": {
                "ssn": "123456789"
              }
            }
          }
        ]
      }
    }`
		w.Write([]byte(resp))
	}

	s, err := MockService(ts.URL)
	assert.NoError(t, err)

	actualResult, err := s.SearchSSN("123456789")
	assert.NoError(t, err)

	assert.Equal(t, "5962", actualResult[0]["imis_id"])
}

func TestSearchSSNMultipleHit(t *testing.T) {
	handler := http.NotFound
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler(w, r)
	}))
	defer ts.Close()

	handler = func(w http.ResponseWriter, r *http.Request) {
		resp := `{
      "took": 1,
      "timed_out": false,
      "_shards": {
        "total": 5,
        "successful": 5,
        "failed": 0
      },
      "hits": {
        "total": 1,
        "max_score": 7.411551,
        "hits": [
          {
            "_index": "members-2017.09.14",
            "_type": "member",
            "_id": "AV5_Z3USj5SI64ssYb5t",
            "_score": 7.411551,
            "_source": {
              "imis_id": "5962",
              "demographics": {
                "ssn": "123456789"
              }
            }
          },
          {
            "_index": "members-2017.09.14",
            "_type": "member",
            "_id": "AV5_Z3USj5SI64ssYb6t",
            "_score": 7.411551,
            "_source": {
              "imis_id": "5965",
              "demographics": {
                "ssn": "123456788"
              }
            }
          }
        ]
      }
    }`

		w.Write([]byte(resp))
	}

	s, err := MockService(ts.URL)
	assert.NoError(t, err)

	actualResult, err := s.SearchSSN("123456789")

	assert.Equal(t, "5962", actualResult[0]["imis_id"])
	assert.Equal(t, "5965", actualResult[1]["imis_id"])
}

func MockService(url string) (Service, error) {
	client, err := elastic.NewSimpleClient(elastic.SetURL(url))
	if err != nil {
		return nil, err
	}
	return &service{client: client}, nil
}
