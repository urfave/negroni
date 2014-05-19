package negroni

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_Static(t *testing.T) {
	response := httptest.NewRecorder()
	response.Body = new(bytes.Buffer)

	n := New()
	n.Use(NewStatic("."))

	req, err := http.NewRequest("GET", "http://localhost:3000/negroni.go", nil)
	if err != nil {
		t.Error(err)
	}
	n.ServeHTTP(response, req)
	expect(t, response.Code, http.StatusOK)
	expect(t, response.Header().Get("Expires"), "")
	if response.Body.Len() == 0 {
		t.Errorf("Got empty body for GET request")
	}
}

func Test_Static_Head(t *testing.T) {
	response := httptest.NewRecorder()
	response.Body = new(bytes.Buffer)

	n := New()
	n.Use(NewStatic("."))

	req, err := http.NewRequest("HEAD", "http://localhost:3000/negroni.go", nil)
	if err != nil {
		t.Error(err)
	}

	n.ServeHTTP(response, req)
	expect(t, response.Code, http.StatusOK)
	if response.Body.Len() != 0 {
		t.Errorf("Got non-empty body for HEAD request")
	}
}

func Test_Static_As_Post(t *testing.T) {
	response := httptest.NewRecorder()

	n := New()
	n.Use(NewStatic("."))

	req, err := http.NewRequest("POST", "http://localhost:3000/negroni.go", nil)
	if err != nil {
		t.Error(err)
	}

	n.ServeHTTP(response, req)
	expect(t, response.Code, http.StatusNotFound)
}

func Test_Static_BadDir(t *testing.T) {
	response := httptest.NewRecorder()

	n := Classic()

	req, err := http.NewRequest("GET", "http://localhost:3000/negroni.go", nil)
	if err != nil {
		t.Error(err)
	}

	n.ServeHTTP(response, req)
	refute(t, response.Code, http.StatusOK)
}
