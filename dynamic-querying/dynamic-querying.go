package dq

import (
  "fmt"
  "io/ioutil"
  "log"
  "strings"
  "encoding/json"
  "path/filepath"
  "github.com/Jeffail/gabs"
)

type queryConfigField struct {
  Min_score_adjustment int `json:"min_score_adjustment"`
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

  // Passes the found query configuration and the SearchRequest array
  // to a function that returns the fully-built "elastic query" as a string.
  elasticQuery, err := buildElasticQuery(req)
  if err != nil {
    return "", err
  }

  // application.go uses this as the body to pass to the elasticsearch service.
  return elasticQuery, nil
}

func buildElasticQuery(req []SearchRequest) (elasticQueryString string, err error) {
  queryConfiguration, err := loadSearchTermConfigurations()
  if err != nil {
    return "", err
  }

  queryFieldJson := gabs.New()

  for _, reqField := range req {
    configField, err := parseFieldConfiguration(queryConfiguration[reqField.Field_name].(map[string]interface{}))
    if err != nil {
      fmt.Println(err)
    }

    err = incorporateJSON(&queryFieldJson, configField, reqField)
    if err != nil {
      log.Fatal(err)
      return "", err
    }

  }

  finalElasticQuery, err := gabs.ParseJSON([]byte(queryConfiguration["base_json"].(string)))
  if err != nil {
    log.Fatal(err)
    return "", err
  }
  finalElasticQuery.Set(queryFieldJson.Data(), strings.Split(queryConfiguration["base_query_path"].(string),";")...)

  return finalElasticQuery.String(), nil
}

func incorporateJSON(finalJsonObject **gabs.Container, configField queryConfigField, reqField SearchRequest) (err error) {
  for _, queryTextRow := range configField.Query_text {

    queryTextJson, err := gabs.ParseJSON([]byte(queryTextRow.Json))
    if err != nil {
      log.Fatal(err)
      return err
    }

    children, _ := queryTextJson.ChildrenMap()

    for key, child := range children {
      if !(*finalJsonObject).Exists(key) {
        (*finalJsonObject).Array(key)
      }

      jsonToAdd, err := gabs.ParseJSON([]byte(child.Index(0).String()))
      if err != nil {
        fmt.Println(err)
      }
      jsonToAdd.Set(reqField.Value, strings.Split(queryTextRow.Value_path, ";")...)
      (*finalJsonObject).ArrayAppend(jsonToAdd.Data(), key)

    }
  }
  return
}

func loadSearchTermConfigurations() (res map[string]interface{}, err error) {

	rawJSONBytes, err := LoadJSONFileBytes("./dynamic-querying/config/query-config.json")
	if err != nil {
		log.Fatal(err)
    return nil, err
	}

  var jsonOutput map[string]interface{}
  json.Unmarshal(rawJSONBytes, &jsonOutput)

  return jsonOutput, nil
}

func parseFieldConfiguration(req map[string]interface{}) (res queryConfigField, err error) {

  jsonInput, err := json.Marshal(req)

  var jsonOutput queryConfigField
  json.Unmarshal(jsonInput, &jsonOutput)

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
