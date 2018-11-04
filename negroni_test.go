package negroni

import (
	"net/http"
	"net/http/httptest"
	"os"
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

func TestNegroniWith(t *testing.T) {
	result := ""
	response := httptest.NewRecorder()

	n1 := New()
	n1.Use(HandlerFunc(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		result = "one"
		next(rw, r)
	}))
	n1.Use(HandlerFunc(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		result += "two"
		next(rw, r)
	}))

	n1.ServeHTTP(response, (*http.Request)(nil))
	expect(t, 2, len(n1.Handlers()))
	expect(t, result, "onetwo")

	n2 := n1.With(HandlerFunc(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		result += "three"
		next(rw, r)
	}))

	// Verify that n1 was left intact and not modified.
	n1.ServeHTTP(response, (*http.Request)(nil))
	expect(t, 2, len(n1.Handlers()))
	expect(t, result, "onetwo")

	n2.ServeHTTP(response, (*http.Request)(nil))
	expect(t, 3, len(n2.Handlers()))
	expect(t, result, "onetwothree")
}

func TestNegroniWith_doNotModifyOriginal(t *testing.T) {
	result := ""
	response := httptest.NewRecorder()

	n1 := New()
	n1.handlers = make([]Handler, 0, 10) // enforce initial capacity
	n1.Use(HandlerFunc(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		result = "one"
		next(rw, r)
	}))

	n1.ServeHTTP(response, (*http.Request)(nil))
	expect(t, 1, len(n1.Handlers()))

	n2 := n1.With(HandlerFunc(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		result += "two"
		next(rw, r)
	}))
	n3 := n1.With(HandlerFunc(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		result += "three"
		next(rw, r)
	}))

	// rebuilds middleware
	n2.UseHandlerFunc(func(rw http.ResponseWriter, r *http.Request) {})
	n3.UseHandlerFunc(func(rw http.ResponseWriter, r *http.Request) {})

	n1.ServeHTTP(response, (*http.Request)(nil))
	expect(t, 1, len(n1.Handlers()))
	expect(t, result, "one")

	n2.ServeHTTP(response, (*http.Request)(nil))
	expect(t, 3, len(n2.Handlers()))
	expect(t, result, "onetwo")

	n3.ServeHTTP(response, (*http.Request)(nil))
	expect(t, 3, len(n3.Handlers()))
	expect(t, result, "onethree")
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

// Ensures that a Negroni middleware chain
// can correctly return all of its handlers.
func TestHandlers(t *testing.T) {
	response := httptest.NewRecorder()
	n := New()
	handlers := n.Handlers()
	expect(t, 0, len(handlers))

	n.Use(HandlerFunc(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		rw.WriteHeader(http.StatusOK)
	}))

	// Expects the length of handlers to be exactly 1
	// after adding exactly one handler to the middleware chain
	handlers = n.Handlers()
	expect(t, 1, len(handlers))

	// Ensures that the first handler that is in sequence behaves
	// exactly the same as the one that was registered earlier
	handlers[0].ServeHTTP(response, (*http.Request)(nil), nil)
	expect(t, response.Code, http.StatusOK)
}

func TestNegroni_Use_Nil(t *testing.T) {
	defer func() {
		err := recover()
		if err == nil {
			t.Errorf("Expected negroni.Use(nil) to panic, but it did not")
		}
	}()

	n := New()
	n.Use(nil)
}

func TestDetectAddress(t *testing.T) {
	if detectAddress() != DefaultAddress {
		t.Error("Expected the DefaultAddress")
	}

	if detectAddress(":6060") != ":6060" {
		t.Error("Expected the provided address")
	}

	os.Setenv("PORT", "8080")
	if detectAddress() != ":8080" {
		t.Error("Expected the PORT env var with a prefixed colon")
	}
}

func voidHTTPHandlerFunc(rw http.ResponseWriter, r *http.Request) {
	// Do nothing
}

// Test for function Wrap
func TestWrap(t *testing.T) {
	response := httptest.NewRecorder()

	handler := Wrap(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(response, (*http.Request)(nil), voidHTTPHandlerFunc)

	expect(t, response.Code, http.StatusOK)
}

// Test for function WrapFunc
func TestWrapFunc(t *testing.T) {
	response := httptest.NewRecorder()

	// WrapFunc(f) equals Wrap(http.HandlerFunc(f)), it's simpler and usefull.
	handler := WrapFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusOK)
	})

	handler.ServeHTTP(response, (*http.Request)(nil), voidHTTPHandlerFunc)

	expect(t, response.Code, http.StatusOK)
}
