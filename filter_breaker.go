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

import (
	"io"
	"log"
	"math/rand"
)

// breaker controls whether the writing of a byte sequence must be broken by the link
type breaker struct{}

func (d breaker) Write(link *link, w io.Writer, p parcel) (parcel, error) {
	if p.isEmpty() {
		return p, nil
	}
	if link.transport.BreakServiceResponse == 0 {
		return p, nil
	}
	pct := rand.Intn(99)+1
	if link.transport.Verbose {
		log.Printf("[%s] broken %d pct of writing %d bytes",
			link.route.Label, link.transport.BreakServiceResponse, p.read)
	}
	if( pct<=link.transport.BreakServiceResponse ) {
		err := errBreak
		return emptyParcel, err
	} else {
		return p, nil
	}
	return p, nil
}
