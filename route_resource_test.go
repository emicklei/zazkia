/*
Copyright 2017 Ernest Micklei

Licensed under the Apache License, Version 2.0 (the "License");
*/

package main

import (
	"net/http"
	"testing"
)

func TestRouteResource_GetRoutes(t *testing.T) {
	setupWithRoute("r1")
	if w := doRequest("GET", "/routes/"); w.Code != 200 {
		t.Errorf("got %d", w.Code)
	}
}

func TestRouteResource_GetLinksForRoute(t *testing.T) {
	setupWithRoute("r1")
	if w := doRequest("GET", "/routes/r1/links"); w.Code != 200 {
		t.Errorf("got %d", w.Code)
	}
}

func TestRouteResource_CloseAllLinksForRoute(t *testing.T) {
	setupWithRoute("r1")
	if w := doRequest("POST", "/routes/r1/links/close"); w.Code != 200 {
		t.Errorf("got %d", w.Code)
	}
}

func TestRouteResource_ToggleAcceptConnections(t *testing.T) {
	setupWithRoute("r1")
	if w := doRequest("POST", "/routes/r1/toggle-accept"); w.Code != 200 {
		t.Errorf("got %d", w.Code)
	}
	if routeMgr.routes[0].Transport.isAcceptConnections() {
		t.Error("expected toggled off")
	}
	if w := doRequest("POST", "/routes/r1/toggle-accept"); w.Code != 200 {
		t.Errorf("got %d", w.Code)
	}
	if w := doRequest("POST", "/routes/unknown/toggle-accept"); w.Code != http.StatusNotFound {
		t.Errorf("got %d", w.Code)
	}
}
