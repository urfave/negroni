package negroni

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Regression test for https://github.com/urfave/negroni/issues/269:
// when a handler panics after the response has already been committed,
// the recovery middleware cannot send a 500 anymore. Appending the panic
// details then only corrupts the in-flight response and leaks internals
// into what still looks like a successful answer to the client.
func TestRecovery_committedResponse(t *testing.T) {
	recorder := httptest.NewRecorder()

	rec := NewRecovery()
	rec.Logger = log.New(bytes.NewBuffer([]byte{}), "[negroni] ", 0)

	n := New()
	n.Use(rec)
	n.UseHandler(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if _, err := res.Write([]byte("partial")); err != nil {
			t.Fatalf("write err: %v", err)
		}
		panic("here is a panic!")
	}))
	n.ServeHTTP(recorder, (*http.Request)(nil))

	if recorder.Body.String() != "partial" {
		t.Errorf("expected body %q, got %q", "partial", recorder.Body.String())
	}
}

func TestRecovery_committedResponseNoPrintStack(t *testing.T) {
	recorder := httptest.NewRecorder()

	rec := NewRecovery()
	rec.Logger = log.New(bytes.NewBuffer([]byte{}), "[negroni] ", 0)
	rec.PrintStack = false

	n := New()
	n.Use(rec)
	n.UseHandler(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if _, err := res.Write([]byte("partial")); err != nil {
			t.Fatalf("write err: %v", err)
		}
		panic("here is a panic!")
	}))
	n.ServeHTTP(recorder, (*http.Request)(nil))

	if recorder.Body.String() != "partial" {
		t.Errorf("expected body %q, got %q", "partial", recorder.Body.String())
	}
}
