package main

import "testing"

func TestReadRoutes(t *testing.T) {
	rs, err := readRoutes()
	t.Log(rs, err)
}
