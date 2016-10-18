package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLinkMatcher(t *testing.T) {
	tokens := linksMatcher.FindStringSubmatch("/links/123/action")
	if got, want := tokens[1], "123"; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	if got, want := tokens[2], "action"; got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestVerbose(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/links/0/toggle-verbose", nil)
	m := linkResource{linkMgr}
	// add link 0
	link := new(link)
	link.resetTransport()
	linkMgr = newLinkManager()
	linkMgr.add(link)
	if got, want := link.transport.Verbose, false; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	// serve
	m.ServeHTTP(w, r)
	if got, want := w.Code, 303; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	if got, want := link.transport.Verbose, true; got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestSend(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/links/0/toggle-send", nil)
	m := linkResource{linkMgr}
	// add link 0
	link := new(link)
	link.resetTransport()
	linkMgr = newLinkManager()
	linkMgr.add(link)
	if got, want := link.transport.SendingToClient, true; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	if got, want := link.transport.SendingToService, true; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	// serve
	m.ServeHTTP(w, r)
	if got, want := w.Code, 303; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	if got, want := link.transport.SendingToClient, false; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	if got, want := link.transport.SendingToService, false; got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestReceive(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/links/0/toggle-receive", nil)
	m := linkResource{linkMgr}
	// add link 0
	link := new(link)
	link.resetTransport()
	linkMgr = newLinkManager()
	linkMgr.add(link)
	if got, want := link.transport.ReceivingFromClient, true; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	if got, want := link.transport.ReceivingFromService, true; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	// serve
	m.ServeHTTP(w, r)
	if got, want := w.Code, 303; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	if got, want := link.transport.ReceivingFromClient, false; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	if got, want := link.transport.ReceivingFromService, false; got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestDelayResponse(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/links/0/delay-response", nil)
	m := linkResource{linkMgr}
	// add link 0
	link := new(link)
	link.resetTransport()
	linkMgr = newLinkManager()
	linkMgr.add(link)
	if got, want := link.transport.DelayServiceResponse, 0; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	// serve
	m.ServeHTTP(w, r)
	if got, want := w.Code, 303; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	if got, want := link.transport.DelayServiceResponse, 10000; got != want {
		t.Errorf("got %v want %v", got, want)
	}
}
