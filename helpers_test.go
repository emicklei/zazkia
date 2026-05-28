/*
Copyright 2017 Ernest Micklei

Licensed under the Apache License, Version 2.0 (the "License");
*/

package main

import (
	"net"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	restful "github.com/emicklei/go-restful/v3"
)

// newTestRoute returns a Route with sensible defaults for tests.
func newTestRoute(label string) *Route {
	return &Route{
		Label:           label,
		ServiceHostname: "localhost",
		ServicePort:     1234,
		ListenPort:      4321,
		Transport:       newDefaultTransportState(),
	}
}

// newTestLink returns a link with piped client/service connections.
// The cleanup func closes both ends of both pipes.
func newTestLink(t *testing.T, r *Route) (*link, func()) {
	t.Helper()
	cConn, cOther := net.Pipe()
	sConn, sOther := net.Pipe()
	l := &link{
		ID:          <-idGen,
		route:       r,
		clientConn:  cConn,
		serviceConn: sConn,
	}
	l.resetTransport()
	cleanup := func() {
		cConn.Close()
		cOther.Close()
		sConn.Close()
		sOther.Close()
	}
	return l, cleanup
}

// setupWithRoute installs a single named route in routeMgr and re-registers
// REST resources so handlers see the new route manager.
func setupWithRoute(label string) *Route {
	setup()
	r := newTestRoute(label)
	routeMgr = routeManager{routes: []*Route{r}}
	restful.DefaultContainer = restful.NewContainer()
	addRESTResources()
	return r
}

// doRequest dispatches an HTTP request through the default container.
func doRequest(method, path string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(method, path, nil)
	restful.DefaultContainer.Dispatch(w, r)
	return w
}

func itoa(i int) string { return strconv.Itoa(i) }

// mustTimer returns a channel that fires after a short period; used to avoid
// indefinite hangs in tests that exercise locking edge cases.
func mustTimer(t *testing.T) <-chan struct{} {
	t.Helper()
	ch := make(chan struct{})
	go func() {
		<-time.After(2 * time.Second)
		close(ch)
	}()
	return ch
}
