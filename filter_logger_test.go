/*
Copyright 2017 Ernest Micklei

Licensed under the Apache License, Version 2.0 (the "License");
*/

package main

import "testing"

func TestPrintable(t *testing.T) {
	got := printable([]byte("ab\ncd\x01ef"))
	if got != "ab\ncd.ef" {
		t.Errorf("got %q", got)
	}
}

func Test_logger_verbose(t *testing.T) {
	r := newTestRoute("r1")
	l, cleanup := newTestLink(t, r)
	defer cleanup()
	l.transport.Verbose = true

	p := parcel{[]byte("hi"), 2, 2}
	logger{accessService: false}.Read(l, nil, p)
	logger{accessService: false}.Write(l, nil, p)
	logger{accessService: true}.Read(l, nil, p)
	logger{accessService: true}.Write(l, nil, p)
	logger{}.Read(l, nil, parcel{})
	logger{}.Write(l, nil, parcel{})
}
