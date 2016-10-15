package main

import "io"

var transportBufferSize = 32 * 1024

type receiver struct{}

func (rc receiver) Read(link *link, r io.Reader) (parcel, error) {
	buffer := make([]byte, transportBufferSize)
	read, err := r.Read(buffer)
	return parcel{buffer, read, 0}, err
}
