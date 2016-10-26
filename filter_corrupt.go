package main

import (
	"io"
	"math/rand"
)

// corrupt mangles the byte sequence for writing
type corrupt struct{}

func (c corrupt) Write(link *link, w io.Writer, p parcel) (parcel, error) {
	m := link.transport.ServiceResponseCorruptMethod
	if len(m) == 0 {
		return p, nil
	}
	if p.read == 0 {
		return p, nil
	}
	switch m {
	case "firstbyte":
		return parcel{p.data[:1], 1, 0}, nil
	case "randomize":
		for i := 0; i < p.read; i++ {
			j := rand.Intn(i + 1)
			p.data[i], p.data[j] = p.data[j], p.data[i]
		}
	}
	return p, nil
}
