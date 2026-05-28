/*
Copyright 2017 Ernest Micklei

Licensed under the Apache License, Version 2.0 (the "License");
*/

package main

import "testing"

func TestAPIGroupsHelper(t *testing.T) {
	g := APIGroups([]APILink{{ID: 1}})
	if len(g) != 1 || len(g[0].Links) != 1 {
		t.Errorf("got %#v", g)
	}
}
