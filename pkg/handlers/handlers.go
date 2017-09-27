package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/unitehere/membership-analytics/pkg/services/members"
)

var (
	membersService members.Service
	errNoHits      = errors.New("No results were found")
)

// Query implements a Validate method used to
// interact with all other query structs with Validate
type Query interface {
	Validate() error
	Size() string
	From() string
}

// The ResponseValues type describes the structure of the all responses.
type ResponseValues struct {
	Error   string          `json:"error,omitempty"`
	Members *members.Member `json:"members,omitempty"`
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
	var payload ResponseValues

	enc := json.NewEncoder(w)
	enc.SetIndent("", "    ")

	from, size, err := getWindowParamsFromURL(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		payload.Error = err.Error()
		enc.Encode(payload)
		return
	}

	ssnQuery := members.SSNQuery{SSN: (r.URL.Query()["q"][0])}
	err = ssnQuery.Validate()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		payload.Error = err.Error()
		enc.Encode(payload)
		return
	}

	searchResult, err := membersService.SearchSSN(ssnQuery, from, size)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	} else if searchResult.TotalHits > 0 {
		payload.Members = &searchResult
	} else {
		payload.Error = errNoHits.Error()
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	enc.Encode(payload)
}

// PostSearchSSN returns a fuzzy matched array of imis_id given a ssn
// r.Post("/ssn", handlers.PostSearchSSN)
func PostSearchSSN(w http.ResponseWriter, r *http.Request) {
	var (
		ssnQuery     members.SSNQuery
		searchResult members.Member
		payload      ResponseValues
	)
	enc := json.NewEncoder(w)
	enc.SetIndent("", "    ")

	err := decodeAndValidate(r, &ssnQuery)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		payload.Error = err.Error()
		enc.Encode(payload)
		return
	}

	from, size, err := getWindowParamsFromQuery(ssnQuery)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		payload.Error = err.Error()
		enc.Encode(payload)
		return
	}

	searchResult, err = membersService.SearchSSN(ssnQuery, from, size)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	} else if searchResult.TotalHits > 0 {
		payload.Members = &searchResult
	} else {
		payload.Error = errNoHits.Error()
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	enc.Encode(payload)
}

// GetSearchName returns a fuzzy matched array of imis_id given a name
// r.Get("/name", handlers.GetSearchSSN)
func GetSearchName(w http.ResponseWriter, r *http.Request) {
	var (
		nameQuery    members.NameQuery
		searchResult members.Member
		payload      ResponseValues
	)

	enc := json.NewEncoder(w)
	enc.SetIndent("", "    ")

	from, size, err := getWindowParamsFromURL(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		payload.Error = err.Error()
		enc.Encode(payload)
		return
	}

	nameQuery = members.NameQuery{FirstName: (r.URL.Query()["first_name"][0]), LastName: (r.URL.Query()["last_name"][0])}
	err = nameQuery.Validate()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		payload.Error = err.Error()
		enc.Encode(payload)
		return
	}

	searchResult, err = membersService.SearchName(nameQuery, from, size)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	} else if searchResult.TotalHits > 0 {
		payload.Members = &searchResult
	} else {
		payload.Error = errNoHits.Error()
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	enc.Encode(payload)
}

// PostSearchName returns a fuzzy matched array of imis_id given a first name and or last name
// r.Post("/name", handlers.PostSearchName)
func PostSearchName(w http.ResponseWriter, r *http.Request) {
	var (
		nameQuery    members.NameQuery
		searchResult members.Member
		payload      ResponseValues
	)
	enc := json.NewEncoder(w)
	enc.SetIndent("", "    ")

	err := decodeAndValidate(r, &nameQuery)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		payload.Error = err.Error()
		enc.Encode(payload)
		return
	}

	from, size, err := getWindowParamsFromQuery(nameQuery)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		payload.Error = err.Error()
		enc.Encode(payload)
		return
	}

	searchResult, err = membersService.SearchName(nameQuery, from, size)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	} else if searchResult.TotalHits > 0 {
		payload.Members = &searchResult
	} else {
		payload.Error = errNoHits.Error()
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	enc.Encode(payload)
}

func decodeAndValidate(r *http.Request, q Query) error {
	if err := json.NewDecoder(r.Body).Decode(q); err != nil {
		return err
	}
	defer r.Body.Close()
	return q.Validate()
}

func getWindowParamsFromURL(r *http.Request) (from int, size int, err error) {
	err = nil

	values := r.URL.Query()
	if fromValues, found := values["from"]; found {
		from, err = strconv.Atoi(fromValues[0])
		if err != nil {
			return 0, 0, err
		}
	} else {
		from = 0
	}
	if sizeValues, found := values["size"]; found {
		size, err = strconv.Atoi(sizeValues[0])
		if err != nil {
			return 0, 0, err
		}
	} else {
		size = 10
	}

	return from, size, err
}

func getWindowParamsFromQuery(q Query) (from int, size int, err error) {
	err = nil

	if q.From() != "" {
		from, err = strconv.Atoi(q.From())
		if err != nil {
			return 0, 0, err
		}
	} else {
		from = 0
	}
	if q.Size() != "" {
		size, err = strconv.Atoi(q.Size())
		if err != nil {
			return 0, 0, err
		}
	} else {
		size = 10
	}

	return from, size, err
}

func writeError(w http.ResponseWriter, status int, err error) {
	w.WriteHeader(status)
	w.Write([]byte(err.Error()))
}
