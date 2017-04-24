//+build go1.8

package negroni

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type pusherRecorder struct {
	*httptest.ResponseRecorder
	pushed bool
}

func newPusherRecorder() *pusherRecorder {
	return &pusherRecorder{ResponseRecorder: httptest.NewRecorder()}
}

func (c *pusherRecorder) Push(target string, opts *http.PushOptions) error {
	c.pushed = true
	return nil
}

func TestResponseWriterPush(t *testing.T) {
	pushable := newPusherRecorder()
	rw := NewResponseWriter(pushable)
	pusher, ok := rw.(http.Pusher)
	expect(t, ok, true)
	err := pusher.Push("", nil)
	if err != nil {
		t.Error(err)
	}
	expect(t, pushable.pushed, true)
}
