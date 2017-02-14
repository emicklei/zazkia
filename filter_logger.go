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
	"io"
	"log"
)

type logger struct {
	accessService bool
}

func (l logger) Read(link *link, r io.Reader, p parcel) (parcel, error) {
	if p.read == 0 {
		return p, nil
	}
	if link.transport.Verbose {
		src := "client"
		if l.accessService {
			src = "service"
		}
		log.Printf("[%s.%d] read %d bytes from %s", link.route.Label, link.ID, p.read, src)
		if !l.accessService {
			// only log what is read from the client
			log.Println(printable(p.data[:p.read]))
		}
	}
	return p, nil
}

func (l logger) Write(link *link, w io.Writer, p parcel) (parcel, error) {
	if p.written == 0 {
		return p, nil
	}
	if link.transport.Verbose {
		target := "client"
		if l.accessService {
			target = "service"
		}
		log.Printf("[%s.%d] written %d bytes to %s", link.route.Label, link.ID, p.written, target)
		if !l.accessService {
			// only log what is written to the client
			log.Println(printable(p.data[:p.written]))
		}
	}
	return p, nil
}

// printable returns a string from bytes with characters that can be displayed on a console.
// linefeeds are preserved
func printable(data []byte) string {
	b := new(bytes.Buffer)
	for _, each := range string(data) {
		if each == '\u000A' { // LF
			b.WriteRune(each)
			continue
		}
		if each >= ' ' && each <= '~' {
			b.WriteRune(each)
			continue
		}
		b.WriteByte(46) // dot
	}
	return b.String()
}
