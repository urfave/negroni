package negroni

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime/debug"
)

type ErrorHandlerFunc interface {
	Handle(err interface{}, report bool)
}

// Recovery is a Negroni middleware that recovers from any panics and writes a 500 if there was one.
type Recovery struct {
	Logger           *log.Logger
	PrintStack       bool
	ErrorHandlerFunc ErrorHandlerFunc
}

// NewRecovery returns a new instance of Recovery
func NewRecovery() *Recovery {
	return &Recovery{
		Logger:     log.New(os.Stdout, "[negroni] ", 0),
		PrintStack: true,
	}
}

// NewRecoveryWithFunc returns a nwe instance of Recovery with an
// attached errorHandlerFunction
func NewRecoveryWithFunc(errorHandlerFunc ErrorHandlerFunc) *Recovery {
	return &Recovery{
		Logger:           log.New(os.Stdout, "[negroni] ", 0),
		PrintStack:       true,
		ErrorHandlerFunc: errorHandlerFunc,
	}
}

func (rec *Recovery) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	defer func() {
		if err := recover(); err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			stack := debug.Stack()
			f := "PANIC: %s\n%s"
			rec.Logger.Printf(f, err, stack)

			if rec.PrintStack {
				fmt.Fprintf(rw, f, err, stack)
			}

			if rec.ErrorHandlerFunc != nil {
				rec.ErrorHandlerFunc.Handle(err, true)
			}
		}
	}()

	next(rw, r)
}
