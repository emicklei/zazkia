package main

import (
	"bytes"
	"io"
	"log"
	"time"
)

const ReadsFromService = true

var TransportBufferSize = 32 * 1024

var tickerFunc = time.NewTicker

func transport(link *link, w io.Writer, r io.Reader, readsFromService bool) error {
	// io.Copy with simulated problems
	buffer := make([]byte, TransportBufferSize)
	for {
		var (
			err       error
			read      int
			throttler *time.Ticker
		)
		doRead := (readsFromService && link.transport.ReceivingFromService) ||
			!readsFromService && link.transport.ReceivingFromClient
		doWrite := (readsFromService && link.transport.SendingToClient) ||
			!readsFromService && link.transport.SendingToService

		if doRead {
			read, err = r.Read(buffer)
			if err != nil {
				return err
			}
		}

		if readsFromService && link.transport.DelayServiceResponse > 0 {
			if *verbose {
				log.Printf("[%s] delay %d ms of sending response with %d bytes",
					link.route.Label, link.transport.DelayServiceResponse, read)
			}
			time.Sleep(time.Duration(link.transport.DelayServiceResponse) * time.Millisecond)
		}

		doThrottle := link.transport.ThrottleServiceResponse > 0
		if doThrottle {
			throttler = time.NewTicker(time.Duration(1000/link.transport.ThrottleServiceResponse) * time.Millisecond) // ms
			defer throttler.Stop()
		}

		if doWrite {
			offset := 0
			towrite := read
			for towrite > 0 {
				var (
					err     error
					written int
				)
				subset := buffer[offset:read]
				if doThrottle {
					<-throttler.C
					written, err = w.Write(subset[:1])
					if err != nil {
						return err
					}
				} else {
					// unthrottled
					written, err = w.Write(subset)
					if err != nil {
						return err
					}
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
