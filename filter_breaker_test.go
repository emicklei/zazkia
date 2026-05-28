/*
Copyright 2017 Ernest Micklei

Licensed under the Apache License, Version 2.0 (the "License");
*/

package main

import (
	"bytes"
	"testing"
)

func Test_breaker_triggers(t *testing.T) {
	r := newTestRoute("r1")
	l, cleanup := newTestLink(t, r)
	defer cleanup()
	w := new(bytes.Buffer)
	p := parcel{[]byte{1, 2, 3}, 3, 0}

	l.transport.BreakServiceResponse = 100
	l.transport.Verbose = true
	if _, err := (breaker{}).Write(l, w, p); err == nil {
		t.Error("expected break error")
	}
	l.transport.BreakServiceResponse = 0
	if _, err := (breaker{}).Write(l, w, p); err != nil {
		t.Errorf("unexpected: %v", err)
	}
	(breaker{}).Write(l, w, parcel{})
}
