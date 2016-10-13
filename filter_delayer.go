package main

import (
	"io"
	"log"
	"time"
)

type delayer struct{}

func (d delayer) Write(link *link, w io.Writer, p parcel) (parcel, error) {
	if p.isEmpty() {
		return p, nil
	}
	if link.transport.DelayServiceResponse == 0 {
		return p, nil
	}
	if *verbose {
		log.Printf("[%s] delay %d ms of writing %d bytes",
			link.route.Label, link.transport.DelayServiceResponse, p.read)
	}
	time.Sleep(time.Duration(link.transport.DelayServiceResponse) * time.Millisecond)
	return p, nil
}
