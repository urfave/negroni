package negroni

import (
	"bytes"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"
)

// LoggerEntry is the structure passed to the template.
type LoggerEntry struct {
	StartTime string
	Status    int
	Duration  time.Duration
	Hostname  string
	Method    string
	Path      string
	Request   *http.Request
}

// LoggerDefaultFormat is the format logged used by the default Logger instance.
var LoggerDefaultFormat = "{{.StartTime}} | {{.Status}} | \t {{.Duration}} | {{.Hostname}} | {{.Method}} {{.Path}}"

// LoggerDefaultDateFormat is the format used for date by the default Logger instance.
var LoggerDefaultDateFormat = time.RFC3339

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
}

// NewLogger returns a new Logger instance
func NewLogger() *Logger {
	logger := &Logger{ALogger: log.New(os.Stdout, "[negroni] ", 0), dateFormat: LoggerDefaultDateFormat}
	logger.SetFormat(LoggerDefaultFormat)
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

	res, ok := rw.(ResponseWriter)
	if !ok {
		res = NewResponseWriter(rw)
		rw = res
	}

	next(rw, r)

	dateFormat := l.dateFormat
	if dateFormat == "" {
		dateFormat = LoggerDefaultDateFormat
	}

	path := ""
	host := ""
	method := ""
	if r != nil {
		host = r.Host
		method = r.Method
		if r.URL != nil {
			path = r.URL.Path
		}
	}

	log := LoggerEntry{
		StartTime: start.Format(dateFormat),
		Status:    res.Status(),
		Duration:  time.Since(start),
		Hostname:  host,
		Method:    method,
		Path:      path,
		Request:   r,
	}

	tmpl := l.template
	if tmpl == nil {
		tmpl = template.Must(template.New("negroni_parser").Parse(LoggerDefaultFormat))
	}

	buff := &bytes.Buffer{}
	tmpl.Execute(buff, log)
	if l.ALogger != nil {
		l.Println(buff.String())
	}
}
