package members

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	elastic "gopkg.in/olivere/elastic.v5"
)

func TestSearchSSN(t *testing.T) {
	cases := []struct {
		testResponsePath   string
		expectedResultPath string
	}{
		{"TestSearchSSN/notfound_response.json", "TestSearchSSN/notfound_expected.json"},
		{"TestSearchSSN/onefound_response.json", "TestSearchSSN/onefound_expected.json"},
		{"TestSearchSSN/multiplefound_response.json", "TestSearchSSN/multiplefound_expected.json"},
	}

	for _, testCase := range cases {
		handler := http.NotFound
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handler(w, r)
		}))
		defer ts.Close()

		resp, err := ioutil.ReadFile(testCase.testResponsePath)
		assert.NoError(t, err)

		handler = func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(resp))
		}

		s, err := MockService(ts.URL)
		assert.NoError(t, err)

		actualResult, err := s.SearchSSN("123456789")
		assert.NoError(t, err)

		expectedResultBytes, err := ioutil.ReadFile(testCase.expectedResultPath)
		assert.NoError(t, err)

		var expectedResult []map[string]interface{}
		err = json.Unmarshal(expectedResultBytes, &expectedResult)
		assert.NoError(t, err)

		assert.Equal(t, expectedResult, actualResult)
	}
}

func MockService(url string) (Service, error) {
	client, err := elastic.NewSimpleClient(elastic.SetURL(url))
	if err != nil {
		return nil, err
	}
	return &service{client: client}, nil
}