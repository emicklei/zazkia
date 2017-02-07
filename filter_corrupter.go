package main

import (
	"io"
	"math/rand"
)

// corrupter mangles the byte sequence for writing
type corrupter struct{}

func (c corrupter) Write(link *link, w io.Writer, p parcel) (parcel, error) {
	m := link.transport.ServiceResponseCorruptMethod
	if len(m) == 0 {
		return p, nil
	}
	if p.read == 0 {
		return p, nil
	}
	switch m {
	case "randomize":
		for i := 0; i < p.read; i++ {
			j := rand.Intn(i + 1)
			p.data[i], p.data[j] = p.data[j], p.data[i]
		}
	}
	return p, nil
}
