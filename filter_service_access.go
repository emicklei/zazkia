package main

import "io"

type serviceAccess struct{}

func (s serviceAccess) Read(link *link, r io.Reader, p parcel) (parcel, error) {
	if !link.transport.ReceivingFromService {
		return emptyParcel, nil
	}
	return receiver{}.Read(link, r)
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
