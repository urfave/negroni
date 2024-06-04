package negroni

import (
	"bufio"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type closeNotifyingRecorder struct {
	*httptest.ResponseRecorder
	closed chan bool
}

func newCloseNotifyingRecorder() *closeNotifyingRecorder {
	return &closeNotifyingRecorder{
		httptest.NewRecorder(),
		make(chan bool, 1),
	}
}

func (c *closeNotifyingRecorder) close() {
	c.closed <- true
}

func (c *closeNotifyingRecorder) CloseNotify() <-chan bool {
	return c.closed
}

type hijackableResponse struct {
	Hijacked bool
}

func newHijackableResponse() *hijackableResponse {
	return &hijackableResponse{}
}

func (h *hijackableResponse) Header() http.Header           { return nil }
func (h *hijackableResponse) Write(buf []byte) (int, error) { return 0, nil }
func (h *hijackableResponse) WriteHeader(code int)          {}
func (h *hijackableResponse) Flush()                        {}
func (h *hijackableResponse) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h.Hijacked = true
	return nil, nil, nil
}

func TestResponseWriterBeforeWrite(t *testing.T) {
	rec := httptest.NewRecorder()
	rw := NewResponseWriter(rec)

	expect(t, rw.Status(), 0)
	expect(t, rw.Written(), false)
}

func TestResponseWriterBeforeFuncHasAccessToStatus(t *testing.T) {
	var status int

	rec := httptest.NewRecorder()
	rw := NewResponseWriter(rec)

	rw.Before(func(w ResponseWriter) {
		status = w.Status()
	})
	rw.WriteHeader(http.StatusCreated)

	expect(t, status, http.StatusCreated)
}

func TestResponseWriterBeforeFuncCanChangeStatus(t *testing.T) {
	rec := httptest.NewRecorder()
	rw := NewResponseWriter(rec)

	// Always respond with 200.
	rw.Before(func(w ResponseWriter) {
		w.WriteHeader(http.StatusOK)
	})

	rw.WriteHeader(http.StatusBadRequest)
	expect(t, rec.Code, http.StatusOK)
}

func TestResponseWriterBeforeFuncChangesStatusMultipleTimes(t *testing.T) {
	rec := httptest.NewRecorder()
	rw := NewResponseWriter(rec)

	rw.Before(func(w ResponseWriter) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	rw.Before(func(w ResponseWriter) {
		w.WriteHeader(http.StatusNotFound)
	})

	rw.WriteHeader(http.StatusOK)
	expect(t, rec.Code, http.StatusNotFound)
}

func TestResponseWriterWritingString(t *testing.T) {
	rec := httptest.NewRecorder()
	rw := NewResponseWriter(rec)

	rw.Write([]byte("Hello world"))

	expect(t, rec.Code, rw.Status())
	expect(t, rec.Body.String(), "Hello world")
	expect(t, rw.Status(), http.StatusOK)
	expect(t, rw.Size(), 11)
	expect(t, rw.Written(), true)
}

func TestResponseWriterWrittenStatusCode(t *testing.T) {
	rec := httptest.NewRecorder()
	rw := NewResponseWriter(rec)

	expect(t, rw.Written(), false)
	for status := http.StatusContinue; status <= http.StatusEarlyHints; status++ {
		if status == http.StatusSwitchingProtocols {
			continue
		}
		rw.WriteHeader(status)
		expected := false
		expect(t, rw.Written(), expected)
	}
	rw.WriteHeader(http.StatusCreated)
	expect(t, rw.Written(), true)

	rw2 := NewResponseWriter(rec)
	expect(t, rw2.Written(), false)
	rw2.WriteHeader(http.StatusSwitchingProtocols)
	expect(t, rw2.Written(), true)

}

func TestResponseWriterWritingStrings(t *testing.T) {
	rec := httptest.NewRecorder()
	rw := NewResponseWriter(rec)

	rw.Write([]byte("Hello world"))
	rw.Write([]byte("foo bar bat baz"))

	expect(t, rec.Code, rw.Status())
	expect(t, rec.Body.String(), "Hello worldfoo bar bat baz")
	expect(t, rw.Status(), http.StatusOK)
	expect(t, rw.Size(), 26)
}

func TestResponseWriterWritingHeader(t *testing.T) {
	rec := httptest.NewRecorder()
	rw := NewResponseWriter(rec)

	rw.WriteHeader(http.StatusNotFound)

	expect(t, rec.Code, rw.Status())
	expect(t, rec.Body.String(), "")
	expect(t, rw.Status(), http.StatusNotFound)
	expect(t, rw.Size(), 0)
}

func TestResponseWriterWritingHeaderTwice(t *testing.T) {
	rec := httptest.NewRecorder()
	rw := NewResponseWriter(rec)

	rw.WriteHeader(http.StatusNotFound)
	rw.WriteHeader(http.StatusInternalServerError)

	expect(t, rec.Code, rw.Status())
	expect(t, rec.Body.String(), "")
	expect(t, rw.Status(), http.StatusNotFound)
	expect(t, rw.Size(), 0)
}

func TestResponseWriterBefore(t *testing.T) {
	rec := httptest.NewRecorder()
	rw := NewResponseWriter(rec)
	result := ""

	rw.Before(func(ResponseWriter) {
		result += "foo"
	})
	rw.Before(func(ResponseWriter) {
		result += "bar"
	})

	rw.WriteHeader(http.StatusNotFound)

	expect(t, rec.Code, rw.Status())
	expect(t, rec.Body.String(), "")
	expect(t, rw.Status(), http.StatusNotFound)
	expect(t, rw.Size(), 0)
	expect(t, result, "barfoo")
}

func TestResponseWriterUnwrap(t *testing.T) {
	rec := httptest.NewRecorder()
	rw := NewResponseWriter(rec)
	switch v := rw.(type) {
	case interface{ Unwrap() http.ResponseWriter }:
		expect(t, v.Unwrap(), rec)
	default:
		t.Error("Does not implement Unwrap()")
	}
}

func TestResponseWriterHijack(t *testing.T) {
	hijackable := newHijackableResponse()
	rw := NewResponseWriter(hijackable)
	hijacker, ok := rw.(http.Hijacker)
	expect(t, ok, true)
	_, _, err := hijacker.Hijack()
	if err != nil {
		t.Error(err)
	}
	expect(t, hijackable.Hijacked, true)
}

func TestResponseWriteHijackNotOK(t *testing.T) {
	hijackable := new(http.ResponseWriter)
	rw := NewResponseWriter(*hijackable)
	_, ok := rw.(http.Hijacker)
	expect(t, ok, false)
}

func TestResponseWriterCloseNotify(t *testing.T) {
	rec := newCloseNotifyingRecorder()
	rw := NewResponseWriter(rec)
	closed := false
	notifier := rw.(http.CloseNotifier).CloseNotify()
	rec.close()
	select {
	case <-notifier:
		closed = true
	case <-time.After(time.Second):
	}
	expect(t, closed, true)
}

func TestResponseWriterNonCloseNotify(t *testing.T) {
	rw := NewResponseWriter(httptest.NewRecorder())
	_, ok := rw.(http.CloseNotifier)
	expect(t, ok, false)
}

func TestResponseWriterFlusher(t *testing.T) {
	rec := httptest.NewRecorder()
	rw := NewResponseWriter(rec)

	_, ok := rw.(http.Flusher)
	expect(t, ok, true)
}

func TestResponseWriter_Flush_marksWritten(t *testing.T) {
	rec := httptest.NewRecorder()
	rw := NewResponseWriter(rec)

	rw.(http.Flusher).Flush()
	expect(t, rw.Status(), http.StatusOK)
	expect(t, rw.Written(), true)
}

// mockReader only implements io.Reader without other methods like WriterTo
type mockReader struct {
	readStr string
	eof     bool
}

func (r *mockReader) Read(p []byte) (n int, err error) {
	if r.eof {
		return 0, io.EOF
	}
	copy(p, []byte(r.readStr))
	r.eof = true
	return len(r.readStr), nil
}

func TestResponseWriterWithoutReadFrom(t *testing.T) {
	writeString := "Hello world"

	rec := httptest.NewRecorder()
	rw := NewResponseWriter(rec)

	n, err := io.Copy(rw, &mockReader{readStr: writeString})
	expect(t, err, nil)
	expect(t, rw.Status(), http.StatusOK)
	expect(t, rw.Written(), true)
	expect(t, rw.Size(), len(writeString))
	expect(t, int(n), len(writeString))
	expect(t, rec.Body.String(), writeString)
}

type mockResponseWriterWithReadFrom struct {
	*httptest.ResponseRecorder
	writtenStr string
}

func (rw *mockResponseWriterWithReadFrom) ReadFrom(r io.Reader) (n int64, err error) {
	bytes, err := io.ReadAll(r)
	if err != nil {
		return 0, err
	}
	rw.writtenStr = string(bytes)
	rw.ResponseRecorder.Write(bytes)
	return int64(len(bytes)), nil
}

func TestResponseWriterWithReadFrom(t *testing.T) {
	writeString := "Hello world"
	mrw := &mockResponseWriterWithReadFrom{ResponseRecorder: httptest.NewRecorder()}
	rw := NewResponseWriter(mrw)
	n, err := io.Copy(rw, &mockReader{readStr: writeString})
	expect(t, err, nil)
	expect(t, rw.Status(), http.StatusOK)
	expect(t, rw.Written(), true)
	expect(t, rw.Size(), len(writeString))
	expect(t, int(n), len(writeString))
	expect(t, mrw.Body.String(), writeString)
	expect(t, mrw.writtenStr, writeString)
}
