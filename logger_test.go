package negroni

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_Logger(t *testing.T) {
	var buff bytes.Buffer
	recorder := httptest.NewRecorder()

	l := NewLogger()
	l.ALogger = log.New(&buff, "[negroni] ", 0)

	n := New()
	// replace log for testing
	n.Use(l)
	n.UseHandler(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusNotFound)
	}))

	req, err := http.NewRequest("GET", "http://localhost:3000/foobar", nil)
	if err != nil {
		t.Error(err)
	}

	n.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusNotFound)
	refute(t, len(buff.String()), 0)
}

func Test_LoggerCustomFormat(t *testing.T) {
	var buff bytes.Buffer
	recorder := httptest.NewRecorder()

	l := NewLogger()
	l.ALogger = log.New(&buff, "[negroni] ", 0)
	l.SetFormat("{{.Request.URL.Query.Get \"foo\"}} {{.Request.UserAgent}} - {{.Status}}")

	n := New()
	n.Use(l)
	n.UseHandler(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte("OK"))
	}))

	userAgent := "Negroni-Test"
	req, err := http.NewRequest("GET", "http://localhost:3000/foobar?foo=bar", nil)
	if err != nil {
		t.Error(err)
	}
	req.Header.Set("User-Agent", userAgent)

	n.ServeHTTP(recorder, req)
	expect(t, strings.TrimSpace(buff.String()), "[negroni] bar "+userAgent+" - 200")
}
