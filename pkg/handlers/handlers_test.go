package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// r.Get("/ssn", handlers.GetSearchSSN)
func TestGetSearchSSN(t *testing.T) {
	req, err := http.NewRequest("GET", "/search/ssn?q=somefakessntomocklater", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetSearchSSN)

	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	// expected := `{"some": thing}`
	// if rr.Body.String() != expected {
	// 	t.Errorf("handler returned unexpected body: got %v want %v",
	// 		rr.Body.String(), expected)
	// }
}
