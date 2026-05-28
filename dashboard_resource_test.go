/*
Copyright 2017 Ernest Micklei

Licensed under the Apache License, Version 2.0 (the "License");
*/

package main

import "testing"

func TestDashboard_GetIndexAndHelp(t *testing.T) {
	setupWithRoute("r1")
	if w := doRequest("GET", "/"); w.Code != 200 {
		t.Errorf("got %d", w.Code)
	}
	if w := doRequest("GET", "/index.html"); w.Code != 200 {
		t.Errorf("got %d", w.Code)
	}
	if w := doRequest("GET", "/help.html"); w.Code != 200 {
		t.Errorf("got %d", w.Code)
	}
}
