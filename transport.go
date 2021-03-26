/*
Copyright 2017 Ernest Micklei

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
			breaker{},
			corrupter{},
			delayer{},
			throttler{},
			sender{},
			logger{!AccessesService},
			counter{!AccessesService})
	} else {
		// reads from client
		readers = append(readers, clientAccess{}, counter{!AccessesService}, logger{!AccessesService})
		// writes to service
		// note: no corrupter,delayer or throttler here
		writers = append(writers,
			serviceAccess{},
			sender{},
			logger{AccessesService},
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
