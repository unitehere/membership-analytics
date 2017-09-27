package dq

import (
  "fmt"
  "io/ioutil"
  "log"
  "encoding/json"
  "encoding/base64"
  "path/filepath"
  "github.com/Jeffail/gabs"
  "net/http"
  "strings"
  "github.com/unitehere/membership-analytics/config"
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
  Error string
  Members resultMembers
}

type resultMembers struct {
	Data      []map[string]interface{} `json:"data"`
	TotalHits int64                    `json:"total_hits"`
}

// Handles the POST request for member search
func SearchMember(w http.ResponseWriter, r *http.Request) {
	var data []SearchRequest

  // Decodes the request body
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		writeCustomError(w, http.StatusBadRequest, "Invalid request data format")
		return
	}

  validated := validateRequestData(data)
  if !validated {
    writeCustomError(w, http.StatusBadRequest, "Invalid request data format")
    return
  }

  enc := json.NewEncoder(w)
  enc.SetIndent("", "    ")

  // Uses request body data to build the elastic query body
	elasticQueryBody, err := buildElasticQuery(data)
	if err != nil {
		writeGenericError(w, http.StatusBadRequest, err)
		return
	}

  // Queries the elastic service
  payload, err := queryElasticService(elasticQueryBody)
  if err != nil {
    writeGenericError(w, http.StatusBadRequest, err)
    return
  }

  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(http.StatusOK)
  enc.Encode(payload)

	return
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

func validateRequestData(req []SearchRequest) (validated bool) {
  // Checks if empty array is sent
  if len(req) == 0 {
    return false
  }

  configurationInitial, _ := loadSearchTermConfigurations()
  configurationBytes, _ := json.Marshal(configurationInitial)
  configuration, _ := gabs.ParseJSON(configurationBytes)

  for _, reqRow := range req {
    // Checks if value or field_name is an empty string
    // This also "indirectly" checks that the JSON fields were spelled correctly
    if strings.TrimSpace(reqRow.Field_name) == "" || strings.TrimSpace(reqRow.Value) == "" {
      return false
    }

    // Ensures that each "Field_name" matches a field in the query-config.json file
    if !configuration.Exists(reqRow.Field_name) {
      return false
    }

  }

  return true
}

func queryElasticService(queryBody string) (res ResponseValues, err error) {

  elasticHttp := &http.Client{}
  bodyReader := strings.NewReader(queryBody)

  elasticUrl := config.Values.ElasticURL + "/" + config.Values.Index + "/_search"

  auth := []byte(config.Values.ElasticUsername + ":" + config.Values.ElasticPassword)
	authHeader := "Basic " + base64.StdEncoding.EncodeToString(auth)

  req, err := http.NewRequest("POST", elasticUrl, bodyReader)
  req.Header.Add("Content-Type", "application/json")
  req.Header.Add("Authorization", authHeader)
  resp, err := elasticHttp.Do(req)

  defer resp.Body.Close()
  contents, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    return ResponseValues{}, err
  }

  var contentsJson map[string]interface{}
  json.Unmarshal(contents, &contentsJson)

  results, err := transformSearchResults(contentsJson)
  if err != nil {
    return ResponseValues{}, err
  }

  var payload ResponseValues
  payload.Members = results
  payload.Error = ""

  return payload, nil

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

	rawJSONBytes, err := loadJSONFileBytes("./dynamic-querying/config/query-config.json")
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

func loadJSONFileBytes(relFilePath string) (rawFile []byte, err error) {
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

func transformSearchResults(searchResult map[string]interface{}) (resultMembers, error) {

  hitsJsonBytes, err := json.Marshal(searchResult["hits"])
  var hitsJson map[string]interface{}
  json.Unmarshal(hitsJsonBytes, &hitsJson)

  secondHitsJsonBytes, err := json.Marshal(hitsJson["hits"])
  var secondHitsJson []map[string]interface{}
  json.Unmarshal(secondHitsJsonBytes, &secondHitsJson)

  resultLength := len(secondHitsJson)
  totalHits := int64(hitsJson["total"].(float64))

	result := make([]map[string]interface{}, resultLength, resultLength)
	for i, hit := range secondHitsJson {
    hitBytes, errJson := json.Marshal(hit["_source"].(map[string]interface{}))
    if errJson != nil {
      fmt.Println(errJson)
    }

		var data map[string]interface{}
		err := json.Unmarshal(hitBytes, &data)
		if err != nil {
			return resultMembers{}, err
		}
		result[i] = data
	}
	member := resultMembers{Data: result, TotalHits: totalHits}

	return member, err
}

func writeGenericError(w http.ResponseWriter, status int, err error) {
	w.WriteHeader(status)
	w.Write([]byte(err.Error()))
}

func writeCustomError(w http.ResponseWriter, status int, errString string) {
  w.WriteHeader(status)
  w.Write([]byte(errString))
}
