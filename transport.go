package main

import (
	"bytes"
	"io"
	"log"
)

const ReadsFromService = true

var TransportBufferSize = 32 * 1024

func transport(link *link, w io.Writer, r io.Reader, readsFromService bool) error {
	// io.Copy with simulated problems
	buffer := make([]byte, TransportBufferSize)
	for {
		var (
			err  error
			read int
		)
		doRead := (readsFromService && link.receivingFromService) ||
			!readsFromService && link.receivingFromClient
		doWrite := (readsFromService && link.sendingToClient) ||
			!readsFromService && link.sendingToService

		if doRead {
			read, err = r.Read(buffer)
			if err != nil {
				return err
			}
		}
		if doWrite {
			offset := 0
			towrite := read
			for towrite > 0 {
				subset := buffer[offset:read]
				written, err := w.Write(subset)
				if err != nil {
					return err
				}
				if *verbose {
					log.Printf("[%s] written %d from %d", link.route.Label, written, read)
					log.Println(printable(subset))
				}
				offset += written
				towrite -= written
			}
		} else {
			if *verbose {
				log.Printf("[%s] flushing %d bytes", link.route.Label, read)
			}
		}
	}
}

func printable(data []byte) string {
	b := new(bytes.Buffer)
	for _, each := range data {
		if each == 10 || each == 13 { // CR,LF
			b.WriteByte(each)
			continue
		}
		if each >= 32 && each <= 126 {
			b.WriteByte(each)
		} else {
			b.WriteByte(46) // dot
		}
	}
	return b.String()
}
