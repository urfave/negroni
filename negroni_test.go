package negroni

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

/* Test Helpers */
func expect(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("Expected %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func refute(t *testing.T, a interface{}, b interface{}) {
	if a == b {
		t.Errorf("Did not expect %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func TestNegroniRun(t *testing.T) {
	// just test that Run doesn't bomb
	go New().Run(":3000")
}

func TestNegroniServeHTTP(t *testing.T) {
	result := ""
	response := httptest.NewRecorder()

	n := New()
	n.Use(HandlerFunc(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		result += "foo"
		next(rw, r)
		result += "ban"
	}))
	n.Use(HandlerFunc(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		result += "bar"
		next(rw, r)
		result += "baz"
	}))
	n.Use(HandlerFunc(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		result += "bat"
		rw.WriteHeader(http.StatusBadRequest)
	}))

	n.ServeHTTP(response, (*http.Request)(nil))

	expect(t, result, "foobarbatbazban")
	expect(t, response.Code, http.StatusBadRequest)
}
