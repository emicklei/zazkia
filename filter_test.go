package main

import (
	"bytes"
	"strings"
	"testing"
)

func Test_serviceAccess(t *testing.T) {
	TransportBufferSize = 4
	defer func() { TransportBufferSize = 1024 }()
	l := new(link)
	l.resetTransport()
	r := strings.NewReader("abcdefghij")

	f := serviceAccess{}
	p, err := f.Read(l, r)
	if err != nil {
		t.Error(err.Error())
	}
	t.Logf("%#v", p)
}

func Test_client_to_service(t *testing.T) {
	//	readers := []linkReader{
	//		clientAccess{},
	//	}
	writers := []linkWriter{
		serviceAccess{},
		delayer{},
		throttler{},
		sender{},
	}
	l := new(link)
	l.resetTransport()
	// READ
	r := strings.NewReader("abcdefghij")
	p, err := clientAccess{}.Read(l, r)
	if err != nil {
		t.Error(err.Error())
	}
	t.Logf("%#v", p)
	w := new(bytes.Buffer)
	// WRITE
	for _, each := range writers {
		np, err := each.Write(l, w, p)
		if err != nil {
			t.Error(err.Error())
			return
		}
		p = np
	}
}
