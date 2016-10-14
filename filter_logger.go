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
	if link.transport.Verbose {
		target := "client"
		if l.accessService {
			target = "service"
		}
		log.Printf("[%s.%d] written %d bytes to %s", link.route.Label, link.ID, p.read, target)
		if !l.accessService {
			// only log what is written to the client
			log.Println(printable(p.data[:p.read]))
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
