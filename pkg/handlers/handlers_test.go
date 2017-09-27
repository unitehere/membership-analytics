package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
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
		{"", "TestGetSearchSSN/invalidinput_response.json", 400},
		{"123456", "TestGetSearchSSN/invalidinput_response.json", 400},
		{"123456789", "TestGetSearchSSN/onefound_response.json", 200},
		{"555555555", "TestGetSearchSSN/notfound_response.json", 200},
		{"404040404", "TestGetSearchSSN/multiplefound_response.json", 200},
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

		assert.Equalf(t, expectedResult, response, "Check out this case: %s", testCase.expectedResultPath)
	}
}

func TestPostSearchName(t *testing.T) {
	cases := []struct {
		requestPath        string
		expectedResultPath string
		status             int
	}{
		{"TestSearchName/invalid_request1.json", "TestSearchName/invalidinput_response.json", 400},
		{"TestSearchName/invalid_request2.json", "TestSearchName/invalidinput_response.json", 400},
		{"TestSearchName/onefound_request.json", "TestSearchName/onefound_response.json", 200},
		{"TestSearchName/notfound_request.json", "TestSearchName/notfound_response.json", 200},
		// TODO multiple found
		// {"404040404", "TestSearchName/multiplefound_response.json", 200},
	}
	membersService = mockService{}

	for _, testCase := range cases {
		req, err := http.NewRequest("POST", "search/name", nil)
		if err != nil {
			t.Fatal(err)
		}
		requestBytes, err := ioutil.ReadFile(testCase.requestPath) // read request as bytes
		if err != nil {
			t.Fatal(err)
		}
		req.Body = ioutil.NopCloser(strings.NewReader(string(requestBytes))) // convert from bytes to string, then set as body

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(PostSearchName)
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

		assert.Equalf(t, expectedResult, response, "Check out %s", testCase.expectedResultPath)
	}
}

type mockService struct {
}

func (s mockService) SearchSSN(ssnQuery members.SSNQuery, from int, size int) (members.Member, error) {
	var result []map[string]interface{}
	if ssnQuery.SSN == "123456789" { // for handlers test, if one found
		result = []map[string]interface{}{map[string]interface{}{"imis_id": "5962"}}
	} else if ssnQuery.SSN == "404040404" {
		result = []map[string]interface{}{map[string]interface{}{"imis_id": "5962"}, map[string]interface{}{"imis_id": "5965"}}
	} else { // for handlers test, if none found
		result = []map[string]interface{}{}
	}
	return members.Member{Data: result, TotalHits: int64(len(result))}, nil
}

func (s mockService) SearchName(nameQuery members.NameQuery, from int, size int) (members.Member, error) {
	var result []map[string]interface{}
	if nameQuery.FirstName == "Alberto" { // handlers test, one found
		result = []map[string]interface{}{map[string]interface{}{"imis_id": "18775", "first_name": "Alberto", "last_name": "Monteiro"}}
	} else {
		result = []map[string]interface{}{}
	}
	return members.Member{Data: result, TotalHits: int64(len(result))}, nil
}
