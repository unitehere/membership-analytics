package dq

import (
  "encoding/json"
  // "io/ioutil"
  "net/http"
  "net/http/httptest"
  "strings"
  "testing"
  // "reflect"
  "fmt"
  // "strconv"

  "github.com/stretchr/testify/assert"
)

type Case struct {
  Name string `json"name"`
  Details caseDetails `json"details"`
}

type caseDetails struct {
  Request map[string]interface{} `json"request"`
  Response caseDetailsResponse `json"response"`
}

type caseDetailsResponse struct {
  Status int `json"status"`
  Data map[string]interface{} `json"data"`
}

type ServiceResponse struct {
  Name string `json"name"`
  Data []map[string]interface{} `json"data"`
  Total_hits int `json"total_hits"`
}

type mockElasticService struct {}
var currentTestCase string
var serviceResults []ServiceResponse

func TestCases(t *testing.T) {
  elasticService = mockElasticService{}
  caseJSONBytes, err := loadJSONFileBytes("./Testing/test_data.json")
	if err != nil {
		fmt.Println(err.Error())
	}
  var cases []Case
  json.Unmarshal(caseJSONBytes, &cases)

  serviceResultsBytes, err := loadJSONFileBytes("./Testing/elasticservice_response_data.json")
  if err != nil {
    fmt.Println(err.Error())
  }
  json.Unmarshal(serviceResultsBytes, &serviceResults)

  for _, testCase := range cases {
    currentTestCase = testCase.Name

    var requestBytes []byte
    var marshalErr error
    if testCase.Name == "requestInvalid_DataNotArray" {
      requestBytes, marshalErr = json.Marshal(testCase.Details.Request["data"])
    } else {
      requestBytes, marshalErr = json.Marshal(testCase.Details.Request["data"])
    }
    if marshalErr != nil {
      fmt.Println(marshalErr)
    }

    bodyReader := strings.NewReader(string(requestBytes))
    req, err := http.NewRequest("POST", "/search/member", bodyReader)
    if err != nil {
      fmt.Println(err)
    }
    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(SearchMember)
    handler.ServeHTTP(rr, req)
    actualStatus := rr.Code

    var actualResponse map[string]interface{}
    err = json.Unmarshal([]byte(rr.Body.String()), &actualResponse)

    fmt.Println(rr.Body.String())
    // fmt.Println(strconv.Itoa(actualStatus) + " == " + strconv.Itoa(testCase.Details.Response.Status))
    assert.Equal(t, testCase.Details.Response.Status, actualStatus, "should be equal")

    assert.Equal(t, actualResponse, testCase.Details.Response.Data, "should be equal")
  }

}

func (mock mockElasticService) QueryElasticService(queryBody string) (ResponseValues, error) {
  var res ResponseValues

  var responseResults ServiceResponse
  for _, responseRow := range serviceResults {
    if responseRow.Name == currentTestCase {
      responseResults = responseRow
    }
  }

  res.Members.TotalHits = int64(responseResults.Total_hits)
  res.Members.Data = responseResults.Data

  return res, nil
}
