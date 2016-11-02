package main

import "io"

// AccessesService is a parameter name
const AccessesService = true

func transport(link *link, w io.Writer, r io.Reader, readsFromService bool) error {
	readers := []linkReader{}
	writers := []linkWriter{}
	if readsFromService {
		readers = append(readers, serviceAccess{}, counter{AccessesService}, logger{AccessesService})
		// writes to client
		writers = append(writers,
			clientAccess{},
			logger{!AccessesService},
			corrupter{},
			delayer{},
			throttler{},
			sender{},
			counter{!AccessesService})
	} else {
		// reads from client
		readers = append(readers, clientAccess{}, counter{!AccessesService}, logger{!AccessesService})
		// writes to service
		// note: no corrupter,delayer or throttler here
		writers = append(writers,
			serviceAccess{},
			logger{AccessesService},
			sender{},
			counter{AccessesService})
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
}
