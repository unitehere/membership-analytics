package dq

import (
  "fmt"
  "io/ioutil"
  "log"
  "encoding/json"
  "path/filepath"
  "github.com/Jeffail/gabs"
  "net/http"
  "strings"
)

type queryConfigField struct {
  Min_score_adjustment int `json:"min_score_adjustment"`
  Query_text []queryText `json:"query_text"`
}

type queryText struct {
  Json string `json:"json"`
  Value_path []string `json:"value_path"`
}

// Data structure of the search requests (application.go file also uses this Type)
type SearchRequest struct {
  Field_name string `json:"field_name"`
  Value string `json:"value"`
}

type ResponseValues struct {
	Values []map[string]interface{} `json:"hits"`
	Error  string                   `json:"error,omitempty"`
}

// Handles the POST request for member search
func SearchMember(w http.ResponseWriter, r *http.Request) {
	var data []SearchRequest
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

  enc := json.NewEncoder(w)
  enc.SetIndent("", "    ")

	elasticQueryBody, err := DynamicQuery(data)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	elasticHttp := &http.Client{}
	bodyReader := strings.NewReader(elasticQueryBody)

	req, err := http.NewRequest("POST", "https://elasticsearch.unitehere.org:9200/members/_search", bodyReader)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization","Basic bWVtYmVyc2hpcF9hbmFseXRpY3M6ckdjbkJqeHhZZzJvSlg=")
	resp, err := elasticHttp.Do(req)

	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
    writeError(w, http.StatusBadRequest, err)
    return
	}
	fmt.Printf("%s\n", string(contents))

	return
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
  min_score := 0
  queryFieldJson := gabs.New()

  for _, reqField := range req {
    configField, err := parseFieldConfiguration(queryConfiguration[reqField.Field_name].(map[string]interface{}))
    if err != nil {
      fmt.Println(err)
    }

    min_score += configField.Min_score_adjustment
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
  finalElasticQuery.Set(min_score, "min_score")

  baseQueryPath := convertPathToStringArray(queryConfiguration["base_query_path"].([]interface{}))
  finalElasticQuery.Set(queryFieldJson.Data(), baseQueryPath...)

  return finalElasticQuery.String(), nil
}

func convertPathToStringArray(path []interface{}) (res []string) {
  pathStringArray := []string{}
  for _, pathRow := range path {
    pathStringArray = append(pathStringArray, pathRow.(string))
  }

  return pathStringArray
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

      jsonToAdd.Set(reqField.Value, queryTextRow.Value_path...)
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

func writeError(w http.ResponseWriter, status int, err error) {
	w.WriteHeader(status)
	w.Write([]byte(err.Error()))
}
