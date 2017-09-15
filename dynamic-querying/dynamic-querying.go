package dq

import (
  "io/ioutil"
  "log"
  "sort"
  "strconv"
  "strings"
  // "fmt"
  "encoding/json"
  "path/filepath"
)

// Option 1 is much more configurable and flexible, but also more complex.
type dynamicQueryResult struct {
  Options dynamicQueryOpt1
  ElasticQuery string
}

type dynamicQueryFieldsOpt1 struct {
  Field_name string `json:"field_name"`
  Apply bool        `json:"apply"`
  Field_type string `json:"field_type"`
  Boost int         `json:"boost"`
  Fuzziness float64 `json:"fuzziness"`
  Analyzer string   `json:"analyzer"`
}

type dynamicQueryOpt1 struct {
  Id int `json:"id"`
  Min_score int `json:"min_score"`
  Size int `json:"size"`
  Fields []dynamicQueryFieldsOpt1 `json:"fields"`
}

// ***NOT USED***
// Option 2 is much simpler, but not configurable for different
// combinations of search terms
// type dynamicQueryOpt2 struct {
//   field_name string
//   field_type string
//   boost int
//   fuzziness float64
//   analyzer string
// }

type SearchRequest struct {
  Field_name string
  Value_text string
  Value_number float64
}

type queryText struct {
  Field_name string
  Query_text string
}

func DynamicQuery(req []SearchRequest) (elasticQueryBody string, err error) {
  queryOpts, err := findQueryConfig(req)
  if err != nil {
    return "", err
  }

  elasticQuery, err := buildElasticQuery(queryOpts, req)
  if err != nil {
    return "", err
  }

  return elasticQuery, nil
}

func findQueryConfig(req []SearchRequest) (res dynamicQueryOpt1, err error) {
  configurations, err := loadSearchTermConfigurations()
  if err != nil {
    return dynamicQueryOpt1{}, err
  }

  queryOptions := dynamicQueryOpt1{}

  // Searches each config row looking for the appropriate query settings
  // using the search request data.
  for _, configRow := range configurations {

    // Assumes current config row is the correct one, until:
    //  - Request contains a search term not in the config row
    //  - Request does not contain a search term the config row expects
    configRowFound := true

    for _, configField := range configRow.Fields {
      // Assumes field in configuration row is not in req row, until proven otherwise.
      configFieldFound := false

      reqRowIndex := sort.Search(len(req), func(i int) bool { return req[i].Field_name == configField.Field_name })
      if reqRowIndex != len(req) {
        // Config field present in the client search request
        configFieldFound = true
        if !configField.Apply {
          // Config field not applied in this config row
          configRowFound = false
        }

      }

      if !configFieldFound && configField.Apply {
        // Config field not present in the client search request,
        // but the config row expects it to be applied.
        configRowFound = false
      }
    }

    if configRowFound {
      queryOptions = configRow
    }
  }

  return queryOptions, nil
}

func buildElasticQuery(queryOpts dynamicQueryOpt1, req []SearchRequest) (elasticQueryString string, err error) {
  queryTextConfig, err := loadQueryTextConfigurations()
  if err != nil {
    log.Fatal(err)
  }
  initialQueryString := "{ \"min_score\": " + strconv.Itoa(queryOpts.Min_score) + ", \"size\": " + strconv.Itoa(queryOpts.Size) + " \"query\": { \"bool\": { \"should\": [ ${QUERY} ]} } } "
  intermediateQueryString := ""

  for _, reqField := range req {
    for _, configField := range queryOpts.Fields {
      if reqField.Field_name == configField.Field_name {
        for _, queryTextConfigField := range queryTextConfig {
          if queryTextConfigField.Field_name == reqField.Field_name {
            queryString := queryTextConfigField.Query_text

            switch configField.Field_type {
              case "string":
                queryString = strings.Replace(queryString, "${VALUE}", reqField.Value_text, -1)
              case "number":
                queryString = strings.Replace(queryString, "${VALUE}", strconv.FormatFloat(reqField.Value_number, 'E', -1, 64), -1)
            }
            queryString = strings.Replace(queryString, "${BOOST}", strconv.Itoa(configField.Boost), -1)
            queryString = strings.Replace(queryString, "${FUZZINESS}", strconv.FormatFloat(configField.Fuzziness, 'E', -1, 64), -1)
            intermediateQueryString = intermediateQueryString + queryString
          }
        }
      }
    }
  }

  finalQueryString := strings.Replace(initialQueryString, "${QUERY}", intermediateQueryString, -1)

  return finalQueryString, nil
}

func loadSearchTermConfigurations() (res []dynamicQueryOpt1, err error) {

	rawJSONBytes, err := loadJSONFileBytes("./dynamic-querying/config/query-config.json")
	if err != nil {
		log.Fatal(err)
	}

  var jsonOutput []dynamicQueryOpt1
  json.Unmarshal(rawJSONBytes, &jsonOutput)

  return jsonOutput, nil
}

func loadQueryTextConfigurations() (res []queryText, err error) {
  rawJSONBytes, err := loadJSONFileBytes("./dynamic-querying/config/query-text.json")
  if err != nil {
    log.Fatal(err)
  }

  var queryTextConfig []queryText
  json.Unmarshal(rawJSONBytes, &queryTextConfig)

  return queryTextConfig, nil
}

func loadJSONFileBytes(relFilePath string) (rawFile []byte, err error) {
  absPath, err := filepath.Abs(relFilePath)
  if err != nil {
    log.Fatal(err)
  }

  rawJSONBytes, err := ioutil.ReadFile(absPath)
  if err != nil {
    log.Fatal(err)
  }

  return rawJSONBytes, nil
}
