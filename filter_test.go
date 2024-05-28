/*
Copyright 2017 Ernest Micklei

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func Test_serviceAccess(t *testing.T) {
	transportBufferSize = 4
	defer func() { transportBufferSize = 1024 }()
	l := new(link)
	l.resetTransport()
	r := strings.NewReader("abcdefghij")

	f := serviceAccess{}
	p, err := f.Read(l, r, parcel{})
	if err != nil {
		t.Error(err.Error())
	}
	t.Logf("%#v", p)
}

func Test_client_to_service(t *testing.T) {
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
	p, err := clientAccess{}.Read(l, r, parcel{})
	if err != nil {
		t.Error(err.Error())
	}
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

func Test_flusher(t *testing.T) {
	transportBufferSize = 4
	defer func() { transportBufferSize = 1024 }()
	l := new(link)
	l.resetTransport()
	w := new(bytes.Buffer)
	p := parcel{[]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, 10, 0}

	f := sender{}
	_, err := f.Write(l, w, p)
	if err != nil {
		t.Error(err.Error())
	}
	if got, want := w.Len(), 10; got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

func Test_throttler(t *testing.T) {
	l := new(link)
	l.resetTransport()
	w := new(bytes.Buffer)
	p := parcel{[]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, 10, 0}

	f := throttler{}
	l.transport.ThrottleServiceResponse = 2
	_, err := f.Write(l, w, p)
	if err != nil {
		t.Error(err.Error())
	}
	if got, want := w.Len(), 10; got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

func Test_delayer(t *testing.T) {
	l := new(link)
	l.resetTransport()
	w := new(bytes.Buffer)
	p := parcel{[]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, 10, 0}

	d := delayer{}
	l.transport.DelayServiceResponse = 100
	now := time.Now()
	_, err := d.Write(l, w, p)
	if err != nil {
		t.Error(err.Error())
	}
	if got, want := w.Len(), 0; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	if got, want := time.Now(), now.Add(100*time.Millisecond); got.Before(want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func Test_receiver(t *testing.T) {
	l := new(link)
	l.resetTransport()
	r := strings.NewReader("abcdefghij")
	f := receiver{}
	p, err := f.Read(l, r)
	if err != nil {
		t.Error(err.Error())
	}
	if got, want := p.read, 10; got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

func Test_corrupt_randomize(t *testing.T) {
	l := new(link)
	l.resetTransport()
	w := new(bytes.Buffer)
	p := parcel{[]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}, 10, 0}

	c := corrupter{}
	l.transport.ServiceResponseCorruptMethod = "randomize"
	np, _ := c.Write(l, w, p)
	if got, want := np.read, 10; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	if got, want := p.data[10], byte(11); got != want {
		t.Errorf("got %v want %v", got, want)
	}
}
