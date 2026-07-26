package negroni

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRecovery(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	panicHandlerCalled := false
	handlerCalled := false

	rec := NewRecovery()
	rec.Logger = log.New(buff, "[negroni] ", 0)
	rec.ErrorHandlerFunc = func(i interface{}) {
		handlerCalled = true
	}
	rec.PanicHandlerFunc = func(i *PanicInformation) {
		panicHandlerCalled = (i != nil)
	}

	n := New()
	// replace log for testing
	n.Use(rec)
	n.UseHandler(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		panic("here is a panic!")
	}))
	n.ServeHTTP(recorder, (*http.Request)(nil))
	expect(t, recorder.Header().Get("Content-Type"), "text/plain; charset=utf-8")
	expect(t, recorder.Code, http.StatusInternalServerError)
	expect(t, panicHandlerCalled, true)
	expect(t, handlerCalled, true)
	refute(t, recorder.Body.Len(), 0)
	refute(t, len(buff.String()), 0)
}

func TestRecovery_noContentTypeOverwrite(t *testing.T) {
	recorder := httptest.NewRecorder()

	rec := NewRecovery()
	rec.Logger = log.New(bytes.NewBuffer([]byte{}), "[negroni] ", 0)

	n := New()
	n.Use(rec)
	n.UseHandler(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "application/javascript; charset=utf-8")
		panic("here is a panic!")
	}))
	n.ServeHTTP(recorder, (*http.Request)(nil))
	expect(t, recorder.Header().Get("Content-Type"), "application/javascript; charset=utf-8")
}

func TestRecovery_callbackPanic(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()

	rec := NewRecovery()
	rec.Logger = log.New(buff, "[negroni] ", 0)
	rec.ErrorHandlerFunc = func(i interface{}) {
		panic("callback panic")
	}

	n := New()
	n.Use(rec)
	n.UseHandler(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		panic("here is a panic!")
	}))
	n.ServeHTTP(recorder, (*http.Request)(nil))

	expect(t, strings.Contains(buff.String(), "callback panic"), true)
}

func TestRecovery_handlerPanic(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()

	rec := NewRecovery()
	rec.Logger = log.New(buff, "[negroni] ", 0)
	rec.PanicHandlerFunc = func(i *PanicInformation) {
		panic("panic handler panic")
	}

	n := New()
	n.Use(rec)
	n.UseHandler(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		panic("here is a panic!")
	}))
	n.ServeHTTP(recorder, (*http.Request)(nil))

	expect(t, strings.Contains(buff.String(), "panic handler panic"), true)
}

type testOutput struct {
	*bytes.Buffer
}

func newTestOutput() *testOutput {
	buf := bytes.NewBufferString("")
	return &testOutput{buf}
}

func (t *testOutput) FormatPanicError(rw http.ResponseWriter, r *http.Request, infos *PanicInformation) {
	fmt.Fprint(t, formatInfos(infos))
}

func formatInfos(infos *PanicInformation) string {
	return fmt.Sprintf("%s %s", infos.RequestDescription(), infos.RecoveredPanic)
}
func TestRecovery_formatter(t *testing.T) {
	recorder := httptest.NewRecorder()
	formatter := newTestOutput()

	req, _ := http.NewRequest("GET", "http://localhost:3003/somePath?element=true", nil)
	var element interface{} = "here is a panic!"
	expectedInfos := &PanicInformation{RecoveredPanic: element, Request: req}

	rec := NewRecovery()
	rec.Formatter = formatter
	n := New()
	n.Use(rec)
	n.UseHandler(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		panic(element)
	}))

	n.ServeHTTP(recorder, req)

	expect(t, formatInfos(expectedInfos), formatter.String())
}

func TestRecovery_PanicInformation(t *testing.T) {
	// Request with query
	req, _ := http.NewRequest("GET", "http://localhost:3003/somePath?element=true", nil)
	var element interface{} = "here is a panic!"
	expectedInfos := &PanicInformation{RecoveredPanic: element, Request: req}

	expect(t, expectedInfos.RequestDescription(), "GET /somePath?element=true")

	// Request without Query
	req, _ = http.NewRequest("POST", "http://localhost:3003/somePath", nil)
	element = "here is a panic!"
	expectedInfos = &PanicInformation{RecoveredPanic: element, Request: req}

	expect(t, expectedInfos.RequestDescription(), "POST /somePath")

	// Nil request
	expectedInfos = &PanicInformation{RecoveredPanic: element, Request: nil}
	expect(t, expectedInfos.RequestDescription(), nilRequestMessage)

	// Stack
	stackValue := "Some Stack element"
	expectedInfos = &PanicInformation{RecoveredPanic: element, Request: req, Stack: []byte(stackValue)}
	expect(t, expectedInfos.StackAsString(), stackValue)
}

func TestRecovery_HTMLFormatter(t *testing.T) {
	recorder := httptest.NewRecorder()
	rec := NewRecovery()
	rec.Formatter = &HTMLPanicFormatter{}
	n := New()
	n.Use(rec)
	n.UseHandler(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		panic("some panic")
	}))

	n.ServeHTTP(recorder, (*http.Request)(nil))
	expect(t, recorder.Header().Get("Content-Type"), "text/html; charset=utf-8")
	refute(t, recorder.Body.Len(), 0)
}

type jsonPanicFormatter struct{}

func (j *jsonPanicFormatter) FormatPanicError(rw http.ResponseWriter, r *http.Request, infos *PanicInformation) {
	rw.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(rw, `{"error":"%v"}`, infos.RecoveredPanic)
}

func TestRecovery_customFormatterHeaders(t *testing.T) {
	// Regression for #241: custom formatters must be able to set Content-Type
	// before the 500 status is committed.
	recorder := httptest.NewRecorder()
	rec := NewRecovery()
	rec.Logger = log.New(bytes.NewBuffer([]byte{}), "[negroni] ", 0)
	rec.Formatter = &jsonPanicFormatter{}

	n := New()
	n.Use(rec)
	n.UseHandler(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		panic("boom")
	}))
	n.ServeHTTP(recorder, httptest.NewRequest("GET", "/", nil))

	expect(t, recorder.Code, http.StatusInternalServerError)
	expect(t, recorder.Header().Get("Content-Type"), "application/json")
	expect(t, strings.Contains(recorder.Body.String(), "boom"), true)
}

type negroniWriterCheckFormatter struct {
	sawNegroniWriter bool
}

func (f *negroniWriterCheckFormatter) FormatPanicError(rw http.ResponseWriter, r *http.Request, infos *PanicInformation) {
	_, f.sawNegroniWriter = rw.(ResponseWriter)
}

func TestRecovery_formatterReceivesNegroniResponseWriter(t *testing.T) {
	recorder := httptest.NewRecorder()
	formatter := &negroniWriterCheckFormatter{}

	rec := NewRecovery()
	rec.Logger = log.New(bytes.NewBuffer([]byte{}), "[negroni] ", 0)
	rec.Formatter = formatter

	n := New()
	n.Use(rec)
	n.UseHandler(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		panic("boom")
	}))
	n.ServeHTTP(recorder, httptest.NewRequest("GET", "/", nil))

	expect(t, formatter.sawNegroniWriter, true)
}
