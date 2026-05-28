/*
Copyright 2017 Ernest Micklei

Licensed under the Apache License, Version 2.0 (the "License");
*/

package main

import "testing"

func TestStartListeners_NoRoutes(t *testing.T) {
	startListeners(nil)
}
