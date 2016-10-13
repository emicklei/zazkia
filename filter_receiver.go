package main

import "io"

type receiver struct{}

func (c receiver) Read(link *link, r io.Reader) (parcel, error) {
	buffer := make([]byte, TransportBufferSize)
	read, err := r.Read(buffer)
	return parcel{buffer, read, 0}, err
}
