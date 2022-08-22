package negroni

import (
	"bufio"
	"net"
	"net/http"
)

const (
	flusher = 1 << iota
	hijacker
	closeNotifier
)

type (
	flusherFeature       struct{ *responseWriter }
	hijackerFeature      struct{ *responseWriter }
	closeNotifierFeature struct{ *responseWriter }
)

func (f flusherFeature) Flush() {
	if !f.Written() {
		// The status will be StatusOK if WriteHeader has not been called yet
		f.WriteHeader(http.StatusOK)
	}
	f.ResponseWriter.(http.Flusher).Flush()
}

func (f hijackerFeature) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return f.ResponseWriter.(http.Hijacker).Hijack()
}

func (f closeNotifierFeature) CloseNotify() <-chan bool {
	return f.ResponseWriter.(http.CloseNotifier).CloseNotify()
}

var featurePicker = make([]func(writer *responseWriter) ResponseWriter, 8)

func initFeaturePicker() {
	featurePicker[0] = func(w *responseWriter) ResponseWriter {
		return w
	}
	featurePicker[flusher] = func(w *responseWriter) ResponseWriter {
		return struct {
			*responseWriter
			http.Flusher
		}{w, flusherFeature{w}}
	}
	featurePicker[hijacker] = func(w *responseWriter) ResponseWriter {
		return struct {
			*responseWriter
			http.Hijacker
		}{w, hijackerFeature{w}}
	}
	featurePicker[closeNotifier] = func(w *responseWriter) ResponseWriter {
		return struct {
			*responseWriter
			http.Flusher
		}{w, flusherFeature{w}}
	}
	featurePicker[flusher|hijacker] = func(w *responseWriter) ResponseWriter {
		return struct {
			*responseWriter
			http.Flusher
			http.Hijacker
		}{w, flusherFeature{w}, hijackerFeature{w}}
	}
	featurePicker[flusher|closeNotifier] = func(w *responseWriter) ResponseWriter {
		return struct {
			*responseWriter
			http.Flusher
			http.CloseNotifier
		}{w, flusherFeature{w}, closeNotifierFeature{w}}
	}
	featurePicker[hijacker|closeNotifier] = func(w *responseWriter) ResponseWriter {
		return struct {
			*responseWriter
			http.Hijacker
			http.CloseNotifier
		}{w, hijackerFeature{w}, closeNotifierFeature{w}}
	}
	featurePicker[flusher|hijacker|closeNotifier] = func(w *responseWriter) ResponseWriter {
		return struct {
			*responseWriter
			http.Flusher
			http.Hijacker
			http.CloseNotifier
		}{w, flusherFeature{w}, hijackerFeature{w}, closeNotifierFeature{w}}
	}
}

func wrapFeature(w *responseWriter) ResponseWriter {
	rw := w.ResponseWriter

	feature := 0
	if _, ok := rw.(http.Flusher); ok {
		feature |= flusher
	}
	if _, ok := rw.(http.Hijacker); ok {
		feature |= hijacker
	}
	if _, ok := rw.(http.CloseNotifier); ok {
		feature |= closeNotifier
	}

	return featurePicker[feature](w)
}
