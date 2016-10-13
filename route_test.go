package main

import "testing"

func TestReadRoutes(t *testing.T) {
	rs, err := readRoutes("routes.json")
	t.Log(rs, err)
}
