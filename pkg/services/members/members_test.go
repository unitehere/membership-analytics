package members

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	elastic "gopkg.in/olivere/elastic.v5"
)

func TestSearchSSN(t *testing.T) {
	cases := []struct {
		testDataPath string
		want         []map[string]interface{}
	}{
		{"testdata/notfound.json", []map[string]interface{}{}},
		{"testdata/onefound.json", []map[string]interface{}{map[string]interface{}{"imis_id": "5962"}}},
		{"testdata/multiplefound.json", []map[string]interface{}{map[string]interface{}{"demographics": map[string]interface{}{"ssn": "123456789"}, "imis_id": "5962"}, map[string]interface{}{"imis_id": "5965", "demographics": map[string]interface{}{"ssn": "123456788"}}}},
	}

	for _, testCase := range cases {
		handler := http.NotFound
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handler(w, r)
		}))
		defer ts.Close()

		resp, err := ioutil.ReadFile(testCase.testDataPath)
		assert.NoError(t, err)

		handler = func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(resp))
		}

		s, err := MockService(ts.URL)
		assert.NoError(t, err)

		actualResult, err := s.SearchSSN("123456789")
		assert.NoError(t, err)

		assert.Equal(t, testCase.want, actualResult)
	}
}

func MockService(url string) (Service, error) {
	client, err := elastic.NewSimpleClient(elastic.SetURL(url))
	if err != nil {
		return nil, err
	}
	return &service{client: client}, nil
}
