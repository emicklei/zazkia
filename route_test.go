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
	        "service-hostname": "some-host-name",
	        "service-port": 1521,
	        "listen-port": 49997,
	        "transport": {
				"throttle-service-response": 1,
				"delay-service-response": 1,
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
	for _, each := range rs {
		if len(each.Label) == 0 {
			t.Error("missing Label")
		}
		if len(each.ServiceHostname) == 0 {
			t.Error("missing ServiceHostname")
		}
		if each.ServicePort == 0 {
			t.Error("missing ServicePort")
		}
		if each.ListenPort == 0 {
			t.Error("missing ListenPort")
		}
		if each.Transport.DelayServiceResponse == 0 {
			t.Error("missing DelayServiceResponse")
		}
		if each.Transport.ThrottleServiceResponse == 0 {
			t.Error("missing ThrottleServiceResponse")
		}
	}
}
