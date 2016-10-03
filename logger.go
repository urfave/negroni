package negroni

import (
	"log"
	"net/http"
	"time"
)

// ALogger interface
type aLogger interface {
	Println(v ...interface{})
	Printf(format string, v ...interface{})
}

// Logger is a middleware handler that logs the request as it goes in and the response as it goes out.
type Logger struct {
	// ALogger implements just enough log.Logger interface to be compatible with other implementations
	aLogger
}

// NewLogger returns a new Logger instance
func NewLogger(l *log.Logger) *Logger {
	return &Logger{l}
}

func (l *Logger) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	start := time.Now()
	l.Printf("Started %s %s", r.Method, r.URL.Path)

	next(rw, r)

	res := rw.(ResponseWriter)
	l.Printf("Completed %v %s in %v", res.Status(), http.StatusText(res.Status()), time.Since(start))
}
