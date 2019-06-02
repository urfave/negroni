package negroni

import (
	"net/http"
	"testing"
)

type voidHandler struct{}

func (vh *voidHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	next(rw, r)
}

func BenchmarkNegroni(b *testing.B) {
	h1 := &voidHandler{}
	h2 := &voidHandler{}
	h3 := &voidHandler{}
	h4 := &voidHandler{}
	h5 := &voidHandler{}
	h6 := &voidHandler{}
	h7 := &voidHandler{}
	h8 := &voidHandler{}
	h9 := &voidHandler{}
	h10 := &voidHandler{}

	n := New(h1, h2, h3, h4, h5, h6, h7, h8, h9, h10)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		n.ServeHTTP(nil, nil)
	}
}
