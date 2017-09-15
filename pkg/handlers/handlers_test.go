package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// r.Get("/ssn", handlers.GetSearchSSN)
func TestGetSearchSSN(t *testing.T) {
	cases := []struct {
		query    string
		expected ResponseValues
		status   int
	}{
		{"", ResponseValues{nil, "You need to pass in a ssn string of at least 7 digits as a q param"}, 400},
		{"123456", ResponseValues{nil, "You need to pass in a ssn string of at least 7 digits as a q param"}, 400},
		{"123456789", ResponseValues{[]map[string]interface{}{map[string]interface{}{"imis_id": "5962"}}, ""}, 200},
		{"555555555", ResponseValues{nil, ""}, 200},
	}
	membersService = mockService{}

	for _, testCase := range cases {
		req, err := http.NewRequest("GET", "/search/ssn?q="+testCase.query, nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(GetSearchSSN)
		handler.ServeHTTP(rr, req)

		// Check the status code is what we expect.
		actualStatus := rr.Code
		assert.Equal(t, testCase.status, actualStatus, "they should be equal")

		// Check the response body is what we expect.
		// response := rr.Body.String()
		response := ResponseValues{}
		json.Unmarshal(rr.Body.Bytes(), &response)
		assert.Equal(t, testCase.expected, response)
	}
}

type mockService struct {
}

func (s mockService) SearchSSN(ssn string) ([]map[string]interface{}, error) {
	if ssn == "123456789" { // pretend this ssn is in the ES service
		return []map[string]interface{}{map[string]interface{}{"imis_id": "5962"}}, nil
	}
	return []map[string]interface{}{}, nil // else it found nothing
}
