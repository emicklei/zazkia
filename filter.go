package main

import "io"

type linkReader interface {
	Read(link *link, r io.Reader, p parcel) (parcel, error)
}

type linkWriter interface {
	Write(link *link, w io.Writer, p parcel) (parcel, error)
}

type parcel struct {
	data    []byte
	read    int
	written int
}

var emptyParcel = parcel{[]byte{}, 0, 0}

func (p parcel) isEmpty() bool { return p.read == p.written }
