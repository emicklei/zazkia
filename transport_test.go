package main

import (
	"bytes"
	"strings"
	"testing"
)

func Test_transport(t *testing.T) {
	l := new(link)
	l.resetTransport()
	r := strings.NewReader("abcdefghij")
	w := new(bytes.Buffer)
	transport(l, w, r, ReadsFromService)
	if got, want := w.Len(), 10; got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

func Test_transport_no_send_service(t *testing.T) {
	l := new(link)
	l.resetTransport()
	l.transport.SendingToClient = false
	r := strings.NewReader("abcdefghij")
	w := new(bytes.Buffer)
	transport(l, w, r, ReadsFromService)
	if got, want := w.Len(), 0; got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

func Test_transport_small_buffer(t *testing.T) {
	TransportBufferSize = 4
	defer func() {
		TransportBufferSize = 32 * 1024
	}()
	l := new(link)
	l.resetTransport()
	r := strings.NewReader("abcdefghij")
	w := new(bytes.Buffer)
	transport(l, w, r, ReadsFromService)
	if got, want := w.Len(), 10; got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

func Test_transport_throttled(t *testing.T) {
	TransportBufferSize = 4
	defer func() {
		TransportBufferSize = 32 * 1024
	}()
	l := new(link)
	l.resetTransport()
	l.transport.ThrottleServiceResponse = 1
	r := strings.NewReader("abcdefghij")
	w := new(bytes.Buffer)
	transport(l, w, r, ReadsFromService)
	if got, want := w.Len(), 10; got != want {
		t.Errorf("got %v want %v", got, want)
	}
}
