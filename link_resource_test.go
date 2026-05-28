/*
Copyright 2017 Ernest Micklei

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	restful "github.com/emicklei/go-restful/v3"
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

func TestLinkResource_Close(t *testing.T) {
	setupWithRoute("r1")
	r := newTestRoute("r1")
	l, cleanup := newTestLink(t, r)
	defer cleanup()
	linkMgr.add(l)

	w := doRequest("POST", "/links/"+itoa(l.ID)+"/close")
	if w.Code != 200 {
		t.Errorf("got %d", w.Code)
	}
}

func TestLinkResource_Close_NotFound(t *testing.T) {
	setup()
	if w := doRequest("POST", "/links/9999/close"); w.Code != http.StatusNotFound {
		t.Errorf("got %d", w.Code)
	}
	if w := doRequest("POST", "/links/abc/close"); w.Code != http.StatusBadRequest {
		t.Errorf("got %d", w.Code)
	}
}

func TestLinkResource_BreakResponse(t *testing.T) {
	setup()
	r := newTestRoute("r1")
	l, cleanup := newTestLink(t, r)
	defer cleanup()
	linkMgr.add(l)
	path := "/links/" + itoa(l.ID) + "/break-response"

	if w := doRequest("POST", path+"?pct=50"); w.Code != 200 {
		t.Errorf("got %d", w.Code)
	}
	if l.transport.BreakServiceResponse != 50 {
		t.Errorf("got %d", l.transport.BreakServiceResponse)
	}
	if w := doRequest("POST", path+"?pct=50"); w.Code != 200 {
		t.Errorf("got %d", w.Code)
	}
	if l.transport.BreakServiceResponse != 0 {
		t.Error("expected reset to 0")
	}
	if w := doRequest("POST", path+"?pct=200"); w.Code != http.StatusBadRequest {
		t.Errorf("got %d", w.Code)
	}
	if w := doRequest("POST", path+"?pct=abc"); w.Code != http.StatusBadRequest {
		t.Errorf("got %d", w.Code)
	}
	if w := doRequest("POST", "/links/9999/break-response?pct=10"); w.Code != http.StatusNotFound {
		t.Errorf("got %d", w.Code)
	}
}

func TestLinkResource_DelayResponse_Invalid(t *testing.T) {
	setup()
	r := newTestRoute("r1")
	l, cleanup := newTestLink(t, r)
	defer cleanup()
	linkMgr.add(l)
	path := "/links/" + itoa(l.ID) + "/delay-response"

	if w := doRequest("POST", path+"?ms=abc"); w.Code != http.StatusBadRequest {
		t.Errorf("got %d", w.Code)
	}
	if w := doRequest("POST", "/links/9999/delay-response?ms=10"); w.Code != http.StatusNotFound {
		t.Errorf("got %d", w.Code)
	}
}

func TestLinkResource_ToggleReadsService(t *testing.T) {
	setup()
	r := newTestRoute("r1")
	l, cleanup := newTestLink(t, r)
	defer cleanup()
	linkMgr.add(l)
	if w := doRequest("POST", "/links/"+itoa(l.ID)+"/toggle-reads-service"); w.Code != 200 {
		t.Errorf("got %d", w.Code)
	}
	if l.transport.ReceivingFromService {
		t.Error("expected toggled off")
	}
	if w := doRequest("POST", "/links/9999/toggle-reads-service"); w.Code != http.StatusNotFound {
		t.Errorf("got %d", w.Code)
	}
}

func TestLinkResource_Stats(t *testing.T) {
	setup()
	r := newTestRoute("r1")
	l, cleanup := newTestLink(t, r)
	defer cleanup()
	linkMgr.add(l)
	if w := doRequest("GET", "/links/"+itoa(l.ID)+"/stats"); w.Code != 200 {
		t.Errorf("got %d", w.Code)
	}
	if w := doRequest("GET", "/links/9999/stats"); w.Code != http.StatusNotFound {
		t.Errorf("got %d", w.Code)
	}
}

func TestLinkResource_GetLinksAndCloseAllWithError(t *testing.T) {
	setupWithRoute("r1")
	r := routeMgr.routes[0]
	l, cleanup := newTestLink(t, r)
	defer cleanup()
	linkMgr.add(l)

	if w := doRequest("GET", "/links/"); w.Code != 200 {
		t.Errorf("got %d", w.Code)
	}
	if w := doRequest("POST", "/links/closeAllWithError"); w.Code != 200 {
		t.Errorf("got %d", w.Code)
	}
}

func TestLinkResource_ToggleNotFound(t *testing.T) {
	setup()
	for _, p := range []string{
		"/links/9999/toggle-verbose",
		"/links/9999/toggle-writes-service",
		"/links/9999/toggle-writes-client",
		"/links/9999/toggle-reads-client",
	} {
		if w := doRequest("POST", p); w.Code != http.StatusNotFound {
			t.Errorf("%s: got %d", p, w.Code)
		}
	}
}
