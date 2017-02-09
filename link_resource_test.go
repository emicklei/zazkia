package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	restful "github.com/emicklei/go-restful"
)

func TestVerbose(t *testing.T) {
	setup()
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/links/0/toggle-verbose", nil)
	// add link 0
	link := new(link)
	link.resetTransport()
	linkMgr.add(link)
	if got, want := link.transport.Verbose, false; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	restful.DefaultContainer.Dispatch(w, r)
	if got, want := w.Code, 200; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	if got, want := link.transport.Verbose, true; got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestWritesService(t *testing.T) {
	setup()
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/links/0/toggle-writes-service", nil)
	// add link 0
	link := new(link)
	link.resetTransport()
	linkMgr.add(link)
	if got, want := link.transport.SendingToService, true; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	// serve
	restful.DefaultContainer.Dispatch(w, r)
	if got, want := w.Code, 200; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	if got, want := link.transport.SendingToService, false; got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestWritesClient(t *testing.T) {
	setup()
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/links/0/toggle-writes-client", nil)
	// add link 0
	link := new(link)
	link.resetTransport()
	linkMgr.add(link)
	if got, want := link.transport.SendingToClient, true; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	// serve
	restful.DefaultContainer.Dispatch(w, r)
	if got, want := w.Code, 200; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	if got, want := link.transport.SendingToClient, false; got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestReceive(t *testing.T) {
	setup()
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/links/0/toggle-reads-client", nil)
	// add link 0
	link := new(link)
	link.resetTransport()
	linkMgr.add(link)
	if got, want := link.transport.ReceivingFromClient, true; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	// serve
	restful.DefaultContainer.Dispatch(w, r)
	if got, want := w.Code, 200; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	if got, want := link.transport.ReceivingFromClient, false; got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestDelayResponse(t *testing.T) {
	setup()
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/links/0/delay-response?ms=100", nil)
	// add link 0
	link := new(link)
	link.resetTransport()
	linkMgr.add(link)
	if got, want := link.transport.DelayServiceResponse, 0; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	// serve
	restful.DefaultContainer.Dispatch(w, r)
	if got, want := w.Code, 200; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	if got, want := link.transport.DelayServiceResponse, 100; got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

func setup() {
	restful.DefaultContainer = restful.NewContainer()
	linkMgr = newLinkManager()
	addRESTResources()
}
