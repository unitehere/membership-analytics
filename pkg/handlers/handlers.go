package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/unitehere/membership-analytics/pkg/services/members"
)

var (
	membersService members.Service
)

// Query implements a Validate method used to
// interact with all other query structs with Validate
type Query interface {
	Validate() error
}

// The ResponseValues type describes the structure of the all responses.
type ResponseValues struct {
	Values []map[string]interface{} `json:"values"`
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
	var (
		payload = ResponseValues{}
	)
	ssnQuery := members.SSNQuery{SSN: (r.URL.Query()["q"][0])}
	err := ssnQuery.Validate()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		payload.Error = err.Error()
		json.NewEncoder(w).Encode(payload)
		return
	}
	searchResult, err := membersService.SearchSSN(ssnQuery.SSN)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	payload.Values = searchResult

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(payload)
	if err != nil {
		panic(err)
	}
	return
}

// PostSearchSSN returns a fuzzy matched array of imis_id given a ssn
// r.Post("/ssn", handlers.PostSearchSSN)
func PostSearchSSN(w http.ResponseWriter, r *http.Request) {
	var (
		ssnQuery     members.SSNQuery
		searchResult map[string]members.Member
		payload      ResponseValues
	)
	err := decodeAndValidate(r, &ssnQuery)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		payload.Error = err.Error()
		json.NewEncoder(w).Encode(payload)
		return
	}
	searchResult, err = membersService.SearchSSN(ssnQuery.SSN)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	payload.Values = searchResult
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(payload)
}

// PostSearchName returns a fuzzy matched array of imis_id given a first name and or last name
// r.Post("/name", handlers.PostSearchName)
func PostSearchName(w http.ResponseWriter, r *http.Request) {
	var (
		nameQuery    members.NameQuery
		searchResult map[string]members.Member
		payload      ResponseValues
	)
	err := decodeAndValidate(r, &nameQuery)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		payload.Error = err.Error()
		json.NewEncoder(w).Encode(payload)
		return
	}
	searchResult, err = membersService.SearchName(nameQuery)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	payload.Values = searchResult
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(payload)
}

func decodeAndValidate(r *http.Request, q Query) error {

	if err := json.NewDecoder(r.Body).Decode(q); err != nil {
		return err
	}
	defer r.Body.Close()
	return q.Validate()
}

// PostSearchName returns a fuzzy matched array of imis_id given a first name and or last name
// r.Post("/name", handlers.PostSearchName)
func PostSearchName(w http.ResponseWriter, r *http.Request) {
	var (
		nameQuery    members.NameQuery
		searchResult []map[string]interface{}
	)
	err := decodeAndValidate(r, &nameQuery)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	searchResult, err = membersService.SearchName(nameQuery)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	payload := ResponseValues{searchResult, ""}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(payload)
}

func decodeAndValidate(r *http.Request, q Query) error {
	if err := json.NewDecoder(r.Body).Decode(q); err != nil {
		return err
	}
	defer r.Body.Close()
	return q.Validate(r)
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
