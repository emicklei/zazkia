package main

import (
	"io"
	"time"
)

type linkReader interface {
	Read(link *link, r io.Reader) (parcel, error)
}

type linkWriter interface {
	Write(link *link, w io.Writer, p parcel) (parcel, error)
}

type parcel struct {
	data    []byte
	read    int
	written int
}

var emptyData = []byte{}

var emptyParcel = parcel{emptyData, 0, 0}

func (p parcel) isEmpty() bool { return p.read == p.written }

type serviceAccess struct{}

func (s serviceAccess) Read(link *link, r io.Reader) (parcel, error) {
	if !link.transport.ReceivingFromService {
		return emptyParcel, nil
	}
	return receiver{}.Read(link, r)
}

type clientAccess struct{}

func (c clientAccess) Read(link *link, r io.Reader) (parcel, error) {
	if !link.transport.ReceivingFromClient {
		return emptyParcel, nil
	}
	return receiver{}.Read(link, r)
}

func (c clientAccess) Write(link *link, w io.Writer, p parcel) (parcel, error) {
	if p.isEmpty() {
		return p, nil
	}
	if !link.transport.SendingToClient {
		return emptyParcel, nil
	}
	return p, nil
}

type receiver struct{}

func (c receiver) Read(link *link, r io.Reader) (parcel, error) {
	buffer := make([]byte, TransportBufferSize)
	read, err := r.Read(buffer)
	return parcel{buffer, read, 0}, err
}

func (s serviceAccess) Write(link *link, w io.Writer, p parcel) (parcel, error) {
	if p.isEmpty() {
		return p, nil
	}
	if !link.transport.SendingToService {
		return emptyParcel, nil
	}
	return p, nil
}

type delayer struct{}

func (d delayer) Write(link *link, w io.Writer, p parcel) (parcel, error) {
	if p.isEmpty() {
		return p, nil
	}
	if link.transport.DelayServiceResponse == 0 {
		return p, nil
	}
	time.Sleep(time.Duration(link.transport.DelayServiceResponse) * time.Millisecond)
	return p, nil
}

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

type sender struct{}

func (s sender) Write(link *link, w io.Writer, p parcel) (parcel, error) {
	if p.isEmpty() {
		return p, nil
	}
	offset := 0
	towrite := p.read
	for towrite > 0 {
		var (
			err     error
			written int
		)
		subset := p.data[offset:p.read]
		written, err = w.Write(subset)
		if err != nil {
			return p, err
		}
		offset += written
		towrite -= written
	}
	return emptyParcel, nil
}
