package main

import "io"

const ReadsFromService = true

func  transport(link *link, w io.Writer, r io.Reader, readsFromService bool) error {
	readers := []linkReader{}
	writers := []linkWriter{}
	if readsFromService {
		readers = append(readers, serviceAccess{}, logger{true})
		writers = append(writers, clientAccess{}, logger{false}, delayer{}, throttler{}, sender{})
	} else {
		// readsFromClient
		readers = append(readers, clientAccess{}, logger{false})
		writers = append(writers, serviceAccess{}, logger{true}, delayer{}, throttler{}, sender{})
	}
	for {
		var p parcel
		for _, each := range readers {
			pp, err := each.Read(link, r, p)
			if err != nil {
				return err
			}
			p = pp
		}
		for _, each := range writers {
			pp, err := each.Write(link, w, p)
			if err != nil {
				return err
			}
			p = pp
		}
	}
	return nil
}
