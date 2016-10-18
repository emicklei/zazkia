package main

import (
	"io/ioutil"
	"os"
	"path"
	"testing"
)

var routesData = `
	[
	    {
	        "label": "oracle",
	        "target-domain": "some-host-name",
	        "target-port": 1521,
	        "listen-port": 49997,
	        "transport": {
				"throttle-service-response": 0,
				"delay-target-response": 0,
				"sending-to-client": true,
				"receiving-from-client": true,
				"sending-to-service": true,
				"receiving-from-service": true,
				"verbose": false
	        }
	    }
	]
`

func TestReadRoutes(t *testing.T) {
	dir := os.TempDir()
	loc := path.Join(dir, "routes.json")
	err := ioutil.WriteFile(loc, []byte(routesData), 0666)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir) // clean up
	rs, err := readRoutes(loc)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(rs), 1; got != want {
		t.Errorf("got %v want %v", got, want)
	}
}
