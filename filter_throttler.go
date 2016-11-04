package main

import (
	"io"
	"time"
)

type throttler struct{}

func (t throttler) Write(link *link, w io.Writer, p parcel) (parcel, error) {
	if p.isEmpty() {
		return p, nil
	}
	if link.transport.ThrottleServiceResponse == 0 {
		return p, nil
	}
	throttler := time.NewTicker(time.Duration(1000/link.transport.ThrottleServiceResponse) * time.Millisecond) // ms
	defer throttler.Stop()
	for _, each := range p.data[:p.read] {
		<-throttler.C
		_, err := w.Write([]byte{each})
		if err != nil {
			return p, err
		}
	}
	return p, nil
}
