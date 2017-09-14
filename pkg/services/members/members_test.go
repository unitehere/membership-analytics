package members

import (
	"fmt"
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
              "member_type": "M",
              "employers": [
                {
                  "thru_date": "2009-01-31 00:00:00.000",
                  "effective_date": "2004-09-28 00:00:00.000",
                  "employer_name": "Konocti",
                  "imis_seqn": "16409",
                  "employer_id": "13898"
                },
                {
                  "primary_employer": true,
                  "effective_date": "2009-03-24 00:00:00.000",
                  "employer_name": "Unknown",
                  "imis_seqn": "52895",
                  "employer_id": "14004"
                }
              ],
              "addresses": [
                {
                  "zip": "95451",
                  "country": "US",
                  "address_type": "HOME",
                  "city": "Kelseyville",
                  "address_1": "123 Test Ave.",
                  "state_province": "CA",
                  "team": "MBR",
                  "imis_seqn": "22860",
                  "preferred": true
                }
              ],
              "imis_id": "5962",
              "last_name": "Schmoe",
              "phone_numbers": [
                {
                  "optin_text": "UNKNOWN",
                  "phone_type": "HOME",
                  "phone": "(707)277-9419",
                  "team": "MBR",
                  "imis_seqn": "19273"
                }
              ],
              "@timestamp": "2017-09-14T08:00:44.119Z",
              "@version": "1",
              "company": "Unknown",
              "category": "F",
              "org_code": "UNITE",
              "first_name": "Joe",
              "status": "A",
              "demographics": {
                "imis_id": "5962",
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
	fmt.Println(actualResult[0]["imis_id"])
	assert.Equal(t, "5962", actualResult[0]["imis_id"])
}

func MockService(url string) (Service, error) {
	client, err := elastic.NewSimpleClient(elastic.SetURL(url))
	if err != nil {
		return nil, err
	}
	return &service{client: client}, nil
}
