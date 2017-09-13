package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	// "gopkg.in/olivere/elastic.v5"
)

// r.Get("/ssn", handlers.GetSearchSSN)
func TestGetSearchSSN(t *testing.T) {
	req, err := http.NewRequest("GET", "/search/ssn?q=killme", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetSearchSSN)

	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	status := rr.Code
	assert.Equal(t, status, 200, "they should be equal")

	// Check the response body is what we expect.
	// expected := `{"some": thing}`
	// if rr.Body.String() != expected {
	// 	t.Errorf("handler returned unexpected body: got %v want %v",
	// 		rr.Body.String(), expected)
	// }
}
