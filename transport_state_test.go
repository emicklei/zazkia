/*
Copyright 2017 Ernest Micklei

Licensed under the Apache License, Version 2.0 (the "License");
*/

package main

import (
	"strings"
	"testing"
)

func TestTransportState(t *testing.T) {
	s := newDefaultTransportState()
	if !s.isAcceptConnections() {
		t.Error("expected accepting")
	}
	if v := s.toggleAcceptConnections(); v {
		t.Error("expected false after toggle")
	}
	if s.isAcceptConnections() {
		t.Error("expected not accepting")
	}
	if str := s.String(); !strings.Contains(str, "accept=false") {
		t.Errorf("unexpected string: %q", str)
	}
}

func TestTransportState_postReadInitialize(t *testing.T) {
	s := &TransportState{}
	s.postReadInitialize()
	if s.mutex == nil {
		t.Error("expected mutex to be set")
	}
}
