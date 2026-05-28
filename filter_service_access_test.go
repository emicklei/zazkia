/*
Copyright 2017 Ernest Micklei

Licensed under the Apache License, Version 2.0 (the "License");
*/

package main

import "testing"

func Test_serviceAccess_blocked(t *testing.T) {
	r := newTestRoute("r1")
	l, cleanup := newTestLink(t, r)
	defer cleanup()
	l.transport.ReceivingFromService = false
	if p, _ := (serviceAccess{}).Read(l, nil, parcel{}); !p.isEmpty() {
		t.Error("expected empty parcel")
	}
	l.transport.SendingToService = false
	p := parcel{[]byte("x"), 1, 0}
	if rp, _ := (serviceAccess{}).Write(l, nil, p); !rp.isEmpty() {
		t.Error("expected empty when blocked")
	}
	(serviceAccess{}).Write(l, nil, parcel{})
}
