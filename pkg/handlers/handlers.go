package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/unitehere/membership-analytics/pkg/services/members"
)

var (
	membersService members.Service
)

// The ResponseValues type describes the structure of the all responses.
type ResponseValues struct {
	Values []map[string]interface{} `json:"values,omitempty"`
	Error  string                   `json:"error,omitempty"`
}

func init() {
	var err error
	membersService, err = members.Client()
	if err != nil {
		panic(err)
	}
}

// GetSearchSSN returns a fuzzy matched array of imis_id given a ssn
// r.Get("/ssn", handlers.GetSearchSSN)
func GetSearchSSN(w http.ResponseWriter, r *http.Request) {
	ssnQuery := r.URL.Query()["q"]

	if len(ssnQuery) == 0 || len(ssnQuery[0]) < 7 {
		payload := ResponseValues{nil, "You need to pass in a ssn string of at least 7 digits as a q param"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(payload)
		return
	}

	searchResult, err := membersService.SearchSSN(ssnQuery[0])
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	payload := ResponseValues{}
	payload.Values = searchResult
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(payload)
	return
}

// PostSearchSSN returns a fuzzy matched array of imis_id given a ssn
// r.Post("/ssn", handlers.PostSearchSSN)
func PostSearchSSN(w http.ResponseWriter, r *http.Request) {
	query, err := getValuesFromBody(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if ssn, ok := query["ssn"]; ok && len(ssn.(string)) >= 7 {
		searchResult, err := membersService.SearchSSN(ssn.(string))
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}

		payload := ResponseValues{}
		payload.Values = searchResult
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(payload)
	} else {
		payload := ResponseValues{nil, "You need to pass in a ssn string of at least 7 digits as a q param"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(payload)
	}
}

func writeError(w http.ResponseWriter, status int, err error) {
	w.WriteHeader(status)
	w.Write([]byte(err.Error()))
}

func getValuesFromBody(r *http.Request) (map[string]interface{}, error) {
	decoder := json.NewDecoder(r.Body)
	var requestValues map[string]interface{}
	err := decoder.Decode(&requestValues)
	defer r.Body.Close()
	return requestValues, err
}
