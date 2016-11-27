package negroni

import (
	"log"
	"net/http"
	"os"
	"time"
)

// ALogger interface
type ALogger interface {
	Println(v ...interface{})
	Printf(format string, v ...interface{})
}

// Logger is a middleware handler that logs the request as it goes in and the response as it goes out.
type Logger struct {
	// ALogger implements just enough log.Logger interface to be compatible with other implementations
	ALogger
}

// NewLogger returns a new Logger instance
func NewLogger() *Logger {
	return &Logger{log.New(os.Stdout, "[negroni] ", 0)}
}

func (l *Logger) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	start := time.Now()

	next(rw, r)

	res := rw.(ResponseWriter)
	l.Printf("%v | %v | \t %v | %v | %v %v", start.Format(time.RFC3339), res.Status(), time.Since(start), r.Host, r.Method, r.URL.Path)
}
