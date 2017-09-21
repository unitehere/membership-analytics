package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/unitehere/membership-analytics/pkg/services/members"
)

// r.Get("/ssn", handlers.GetSearchSSN)
func TestGetSearchSSN(t *testing.T) {
	cases := []struct {
		query              string
		expectedResultPath string
		status             int
	}{
<<<<<<< HEAD
		{"", "TestGetSearchSSN/invalidinput_response.json", 400},
		{"123456", "TestGetSearchSSN/invalidinput_response.json", 400},
		{"123456789", "TestGetSearchSSN/onefound_response.json", 200},
		{"555555555", "TestGetSearchSSN/notfound_response.json", 200},
=======
	// {"", "TestSearchSSN/invalidinput_response.json", 400},
	// {"123456", "TestSearchSSN/invalidinput_response.json", 400},
	// {"123456789", "TestSearchSSN/onefound_response.json", 200},
	// {"555555555", "TestSearchSSN/notfound_response.json", 200},
>>>>>>> Fix tests to adhere to new response format
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
		expectedResultBytes, err := ioutil.ReadFile(testCase.expectedResultPath)
		assert.NoError(t, err)
		var expectedResult map[string]interface{}
		err = json.Unmarshal(expectedResultBytes, &expectedResult)
		assert.NoError(t, err)

		var response map[string]interface{}
		err = json.Unmarshal(rr.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, expectedResult, response)
	}
}

type mockService struct {
}

func (s mockService) SearchSSN(placeholder string) (map[string]members.Member, error) {
	return map[string]members.Member{}, nil
}

func (s mockService) SearchName(placeholder members.NameQuery) (map[string]members.Member, error) {
	return map[string]members.Member{}, nil
}
