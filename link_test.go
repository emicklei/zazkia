/*
Copyright 2017 Ernest Micklei

Licensed under the Apache License, Version 2.0 (the "License");
*/

package main

import (
	"errors"
	"net"
	"strings"
	"testing"
)

func TestLink_NewAPILink_AndErrors(t *testing.T) {
	r := newTestRoute("r1")
	l, cleanup := newTestLink(t, r)
	defer cleanup()
	l.clientError = errors.New("c-err")
	l.serviceError = errors.New("s-err")

	if got := l.clientErrorString(); got != "c-err" {
		t.Errorf("got %q", got)
	}
	if got := l.serviceErrorString(); got != "s-err" {
		t.Errorf("got %q", got)
	}
	api := NewAPILink(l)
	if api.ID != l.ID || api.ClientReceiveError != "c-err" || api.ServiceReceiveError != "s-err" {
		t.Errorf("unexpected APILink: %#v", api)
	}

	l.clientError = nil
	l.serviceError = nil
	if l.clientErrorString() != "" || l.serviceErrorString() != "" {
		t.Error("expected empty error strings")
	}
}

func TestLink_String(t *testing.T) {
	r := newTestRoute("r1")
	l, cleanup := newTestLink(t, r)
	defer cleanup()
	if s := l.String(); !strings.Contains(s, "r1") {
		t.Errorf("missing label in %q", s)
	}
}

func TestTransportStats_String(t *testing.T) {
	ts := TransportStats{BytesSentToClient: 1, BytesSentToService: 2, BytesReceivedFromClient: 3, BytesReceivedFromService: 4}
	if s := ts.String(); !strings.Contains(s, "client sent:3") {
		t.Errorf("got %q", s)
	}
}

func TestNewLink_UsesRouteTransport(t *testing.T) {
	r := newTestRoute("r1")
	r.Transport.Verbose = true
	cConn, cOther := net.Pipe()
	sConn, sOther := net.Pipe()
	defer cConn.Close()
	defer cOther.Close()
	defer sConn.Close()
	defer sOther.Close()

	l := newLink(r, cConn, sConn)
	if !l.transport.Verbose {
		t.Error("expected route transport copied")
	}

	r2 := newTestRoute("r2")
	r2.Transport = nil
	l2 := newLink(r2, cConn, sConn)
	if !l2.transport.SendingToClient {
		t.Error("expected default transport")
	}
}

func TestLink_Disconnect(t *testing.T) {
	r := newTestRoute("r1")
	l, cleanup := newTestLink(t, r)
	defer cleanup()
	if err := l.disconnect(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestLink_DisconnectBranches(t *testing.T) {
	// both nil
	l := &link{}
	if err := l.disconnect(); err != nil {
		t.Errorf("unexpected: %v", err)
	}
	// only serviceConn
	_, sOther := net.Pipe()
	defer sOther.Close()
	sConn, sOther2 := net.Pipe()
	defer sOther2.Close()
	l = &link{serviceConn: sConn}
	if err := l.disconnect(); err != nil {
		t.Errorf("unexpected: %v", err)
	}
}
