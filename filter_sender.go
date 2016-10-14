package main

import "io"

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
