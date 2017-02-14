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
)

func Test_transport(t *testing.T) {
	l := new(link)
	l.resetTransport()
	r := strings.NewReader("abcdefghij")
	w := new(bytes.Buffer)
	transport(l, w, r, AccessesService)
	if got, want := w.Len(), 10; got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

func Test_transport_no_send_service(t *testing.T) {
	l := new(link)
	l.resetTransport()
	l.transport.SendingToClient = false
	r := strings.NewReader("abcdefghij")
	w := new(bytes.Buffer)
	transport(l, w, r, AccessesService)
	if got, want := w.Len(), 0; got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

func Test_transport_small_buffer(t *testing.T) {
	transportBufferSize = 4
	defer func() {
		transportBufferSize = 32 * 1024
	}()
	l := new(link)
	l.resetTransport()
	r := strings.NewReader("abcdefghij")
	w := new(bytes.Buffer)
	transport(l, w, r, AccessesService)
	if got, want := w.Len(), 10; got != want {
		t.Errorf("got %v want %v", got, want)
	}
}
