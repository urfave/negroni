package negroni

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"runtime"
)

// Recovery is a Negroni middleware that recovers from any panics and writes a 500 if there was one.
type Recovery struct {
	Logger     *log.Logger
	PrintStack bool
	StackAll   bool
	StackSize  int
}

// NewRecovery returns a new instance of Recovery
func NewRecovery() *Recovery {
	return &Recovery{
		Logger:     log.New(os.Stdout, "[negroni] ", 0),
		PrintStack: true,
		StackAll:   false,
		StackSize:  1024 * 8,
	}
}

func (rec *Recovery) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	nr := httptest.NewRecorder()
	nrw := NewResponseWriter(nr)

	defer func() {
		if err := recover(); err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			stack := make([]byte, rec.StackSize)
			stack = stack[:runtime.Stack(stack, rec.StackAll)]

			f := "PANIC: %s\n%s"
			rec.Logger.Printf(f, err, stack)

			if rec.PrintStack {
				fmt.Fprintf(rw, f, err, stack)
			}
		} else {
			for k, v := range nrw.Header() {
				rw.Header()[k] = v
			}
			rw.WriteHeader(nr.Code)
			rw.Write(nr.Body.Bytes())
		}
	}()

	next(nrw, r)
}
