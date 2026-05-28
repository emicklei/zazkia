/*
Copyright 2017 Ernest Micklei

Licensed under the Apache License, Version 2.0 (the "License");
*/

package main

import "testing"

func Test_counter_branches(t *testing.T) {
	r := newTestRoute("r1")
	l, cleanup := newTestLink(t, r)
	defer cleanup()
	p := parcel{read: 5, written: 5}
	counter{accessService: true}.Read(l, nil, p)
	counter{accessService: false}.Read(l, nil, p)
	counter{accessService: true}.Write(l, nil, p)
	counter{accessService: false}.Write(l, nil, p)
	if l.stats.BytesReceivedFromService != 5 || l.stats.BytesReceivedFromClient != 5 ||
		l.stats.BytesSentToService != 5 || l.stats.BytesSentToClient != 5 {
		t.Errorf("unexpected stats: %#v", l.stats)
	}
}
