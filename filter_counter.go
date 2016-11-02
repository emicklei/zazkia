package main

import "io"

// counter will update the stats of the link.
type counter struct {
	accessService bool
}

func (c counter) Read(link *link, r io.Reader, p parcel) (parcel, error) {
	if c.accessService {
		link.stats.bytesReceivedFromService += int64(p.read)
	} else {
		link.stats.bytesReceivedFromClient += int64(p.read)
	}
	return p, nil

}

func (c counter) Write(link *link, w io.Writer, p parcel) (parcel, error) {
	if c.accessService {
		link.stats.bytesSentToService += int64(p.written)
	} else {
		link.stats.bytesSentToClient += int64(p.written)
	}
	return p, nil
}
