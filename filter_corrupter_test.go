/*
Copyright 2017 Ernest Micklei

Licensed under the Apache License, Version 2.0 (the "License");
*/

package main

import "testing"

func Test_corrupter_empty(t *testing.T) {
	r := newTestRoute("r1")
	l, cleanup := newTestLink(t, r)
	defer cleanup()
	if _, err := (corrupter{}).Write(l, nil, parcel{[]byte("x"), 1, 0}); err != nil {
		t.Error(err)
	}
	l.transport.ServiceResponseCorruptMethod = "randomize"
	if _, err := (corrupter{}).Write(l, nil, parcel{}); err != nil {
		t.Error(err)
	}
}
