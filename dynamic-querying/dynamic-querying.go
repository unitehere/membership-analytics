package dq

import (
  "io/ioutil"
  "log"
  "strings"
  "encoding/json"
  "path/filepath"
  "github.com/Jeffail/gabs"
)

// Mimics the JSON structure in the 'query-config.json' file.
type dynamicQuery struct {
  Id int `json:"id"`
  Query_base string `json:"query_base"`
  Query_path string `json:"query_path"`
  Fields []dynamicQueryFields `json:"fields"`
}

type dynamicQueryFields struct {
  Field_name string `json:"field_name"`
  Apply bool        `json:"apply"`
  Query_text []queryText `json:"query_text"`
}

type queryText struct {
  Json string `json:"json"`
  Value_path string `json:"value_path"`
}

// Data structure of the search requests (application.go file also uses this Type)
type SearchRequest struct {
  Field_name string `json:"field_name"`
  Value string `json:"value"`
}

// Only library function publicly exposed
// Takes a SearchRequest array parameter, and dynamically
// builds (and returns) the elastic query body.
func DynamicQuery(req []SearchRequest) (elasticQueryBody string, err error) {
  // First uses the request to locate the appropriate
  // "query configuration" to use to build the elastic query.
  queryOpts, err := findQueryConfig(req)
  if err != nil {
    return "", err
  }

  // Passes the found query configuration and the SearchRequest array
  // to a function that returns the fully-built "elastic query" as a string.
  elasticQuery, err := buildElasticQuery(queryOpts, req)
  if err != nil {
    return "", err
  }

  // application.go uses this as the body to pass to the elasticsearch service.
  return elasticQuery, nil
}

func findQueryConfig(req []SearchRequest) (res dynamicQuery, err error) {
  configurations, err := loadSearchTermConfigurations()
  if err != nil {
    return dynamicQuery{}, err
  }

  queryOptions := dynamicQuery{}

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

      for _, reqRow := range req {
        if reqRow.Field_name == configField.Field_name {
          // Config field present in the client search request
          configFieldFound = true
          if !configField.Apply {
            // Config field not applied in this config row
            configRowFound = false
          }
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

func buildElasticQuery(queryOpts dynamicQuery, req []SearchRequest) (elasticQueryString string, err error) {

  queryFieldJson := gabs.New()

  for _, reqField := range req {
    for _, configField := range queryOpts.Fields {
      if configField.Field_name == reqField.Field_name {
        for _, queryTextRow := range configField.Query_text {
          queryTextJson, err := gabs.ParseJSON([]byte(queryTextRow.Json))

          children, _ := queryTextJson.ChildrenMap()
          for key, child := range children {
            arrayIndex := 0

            if queryFieldJson.Exists(key) {
              arrayIndex, err = queryFieldJson.ArrayCount(key)
            }

            if !queryFieldJson.Exists(key) {
              queryFieldJson.Array(key)
            }

            err := incorporateJSON(&queryFieldJson, child.Index(0).String(), key, arrayIndex, queryTextRow.Value_path, reqField)
            if err != nil {
              log.Fatal(err)
              return "", err
            }
          }

          if err != nil {
            log.Fatal(err)
            return "", err
          }
        }
      }
    }
  }

  finalElasticQuery, err := gabs.ParseJSON([]byte(queryOpts.Query_base))
  if err != nil {
    log.Fatal(err)
    return "", err
  }
  finalElasticQuery.Set(queryFieldJson.Data(), strings.Split(queryOpts.Query_path,";")...)

  return finalElasticQuery.String(), nil
}

func incorporateJSON(finalJsonObject **gabs.Container, jsonStringToAdd string, arrayToAddTo string, arrayIndexToAddTo int, path string, reqField SearchRequest) (err error) {
  jsonToAdd, err := gabs.ParseJSON([]byte(jsonStringToAdd))
  if err != nil {
    log.Fatal(err)
    return err
  }

  jsonToAdd.Set(reqField.Value, strings.Split(path, ";")...)
  (*finalJsonObject).ArrayAppend(jsonToAdd.Data(), arrayToAddTo)

  return
}

func loadSearchTermConfigurations() (res []dynamicQuery, err error) {

	rawJSONBytes, err := LoadJSONFileBytes("./dynamic-querying/config/query-config.json")
	if err != nil {
		log.Fatal(err)
    return nil, err
	}

  var jsonOutput []dynamicQuery
  json.Unmarshal(rawJSONBytes, &jsonOutput)

  return jsonOutput, nil
}

func LoadJSONFileBytes(relFilePath string) (rawFile []byte, err error) {
  absPath, err := filepath.Abs(relFilePath)
  if err != nil {
    log.Fatal(err)
    return nil, err
  }

  rawJSONBytes, err := ioutil.ReadFile(absPath)
  if err != nil {
    log.Fatal(err)
    return nil, err
  }

  return rawJSONBytes, nil
}
