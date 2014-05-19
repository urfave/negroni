package negroni

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_Recovery(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()

	rec := NewRecovery()
	rec.Logger = log.New(buff, "[negroni] ", 0)

	m := New()
	// replace log for testing
	m.Use(rec)
	m.UseHandler(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		panic("here is a panic!")
	}))
	m.ServeHTTP(recorder, (*http.Request)(nil))
	expect(t, recorder.Code, http.StatusInternalServerError)
	refute(t, recorder.Body.Len(), 0)
	refute(t, len(buff.String()), 0)
}
