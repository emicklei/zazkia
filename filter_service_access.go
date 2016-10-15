package main

import "io"

// serviceAccess is a filter to control whether reading or writing a service is allowed for the link.
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
