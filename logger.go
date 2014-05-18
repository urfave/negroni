package negroni

import (
	"log"
	"net/http"
	"os"
	"time"
)

// Logger is a middleware handler that logs the request as it goes in and the response as it goes out.
type Logger struct {
	Logger *log.Logger
}

func NewLogger() *Logger {
	return &Logger{
		Logger: log.New(os.Stdout, "[negroni] ", 0),
	}
}

func (l *Logger) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	start := time.Now()
	log.Printf("Started %s %s", r.Method, r.URL.Path)

	next(rw, r)

	res := rw.(ResponseWriter)
	log.Printf("Completed %v %s in %v\n", res.Status(), http.StatusText(res.Status()), time.Since(start))
}
