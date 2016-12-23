package negroni

import (
	"bytes"

	"log"
	"net/http"
	"os"
	"sync"
	"text/template"
	"time"
)

// Log is the structure
// passed to the template.
type Log struct {
	StartTime string
	Status    int
	Duration  time.Duration
	Hostname  string
	Method    string
	Path      string
}

// DefaultFormat is the format
// logged used by the default Logger instance.
var DefaultFormat = "{{.StartTime}} | {{.Status}} | \t {{.Duration}} | {{.Hostname}} | {{.Method}} {{.Path}} \n"

// DefaultDateFormat is the
// format used for date by the
// default Logger instance.
var DefaultDateFormat = time.RFC3339

// ALogger interface
type ALogger interface {
	Println(v ...interface{})
	Printf(format string, v ...interface{})
}

// Logger is a middleware handler that logs the request as it goes in and the response as it goes out.
type Logger struct {
	// ALogger implements just enough log.Logger interface to be compatible with other implementations
	ALogger
	dateFormat string
	template   *template.Template
	buffer     *bytes.Buffer
	mu         sync.Mutex
}

// NewLogger returns a new Logger instance
func NewLogger() *Logger {
	logger := &Logger{ALogger: log.New(os.Stdout, "[negroni] ", 0), dateFormat: DefaultDateFormat, buffer: bytes.NewBufferString("")}
	logger.SetFormat(DefaultFormat)
	return logger
}

func (l *Logger) SetFormat(format string) {
	l.template = template.Must(template.New("negroni_parser").Parse(format))
}

func (l *Logger) SetDateFormat(format string) {
	l.dateFormat = format
}

func (l *Logger) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	start := time.Now()

	next(rw, r)

	res := rw.(ResponseWriter)
	log := Log{start.Format(l.dateFormat), res.Status(), time.Since(start), r.Host, r.Method, r.URL.Path}

	l.mu.Lock()
	defer l.mu.Unlock()
	l.template.Execute(l.buffer, log)
	l.Printf(l.buffer.String())
	l.buffer.Reset()
}
