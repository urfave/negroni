package negroni

import (
	"fmt"
	//	"io/ioutil"
	"bufio"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"
)

/* Test Helpers */
func expect(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("Expected %v (type %v) - Got [%v](type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func refute(t *testing.T, a interface{}, b interface{}) {
	if a == b {
		t.Errorf("Did not expect %v (type %v) - Got [%v] (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func expect2(t *testing.T, msg string, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("%s: Expected %v (type %v) - Got [%v] (type %v)", msg, b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func refute2(t *testing.T, msg string, a interface{}, b interface{}) {
	if a == b {
		t.Errorf("%s: Did not expect %v (type %v) - Got [%v] (type %v)", msg, b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func TestNegroniRun(t *testing.T) {
	// just test that Run doesn't bomb
	go New().Run(":3000")
}

type simpleServer struct{}

const helloWorld = "hello world"

// Test some options.
func TestNegroniRunWithOptions(t *testing.T) {
	n := New()
	n.Use(HandlerFunc(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		hello := []byte(helloWorld)
		rw.Write(hello)
	}))
	s := &http.Server{
		Addr:         ":3001",
		ReadTimeout:  time.Second,
		WriteTimeout: time.Second,
	}
	go n.RunWithOptions(s)
	// Just a bit of sleep to let it start.
	time.Sleep(time.Second)

	times := []int{900, 950, 1050, 1100, 1500}
	for i := range times {
		millis := times[i]
		timeout, _ := time.ParseDuration(fmt.Sprintf("%dms", millis))
		conn, err := net.Dial("tcp", "localhost:3001")
		if err != nil {
			t.Error(err)
		}
		time.Sleep(timeout)
		fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")
		rdr := bufio.NewReader(conn)
		status, err := rdr.ReadString('\n')
		status = strings.TrimSpace(status)
		if millis < 1000 {
			expect2(t, "Expected status 200", "HTTP/1.0 200 OK", status)
			expect2(t, fmt.Sprintf("Expected no error for timeout %s", timeout), err, nil)
			resp := make([]byte, 1024)
			_, err := rdr.Read(resp)
			expect2(t, fmt.Sprintf("Expected no error for timeout %s", timeout), err, nil)
			responseStr := strings.TrimSpace(string(resp))
			idx := strings.Index(responseStr, helloWorld)
//			log.Printf("::: %d\n", idx)
			msg := fmt.Sprintf("Expect response ending with [%s], received [%s]\n", idx, responseStr)
			refute2(t, msg, idx, -1)
			log.Printf("Got response\n%s\n", responseStr)
		} else {
			msg := fmt.Sprintf("Expected error for timeout %s", millis)
			refute2(t, msg, err, nil)
			expect2(t, msg, err, io.EOF)
		}
	}
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
