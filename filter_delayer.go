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
	"time"
)

// delayer controls whether the writing of a byte sequence must be delayed by the link
type delayer struct{}

func (d delayer) Write(link *link, w io.Writer, p parcel) (parcel, error) {
	if p.isEmpty() {
		return p, nil
	}
	if link.transport.DelayServiceResponse == 0 {
		return p, nil
	}
	if link.transport.Verbose {
		log.Printf("[%s] delay %d ms of writing %d bytes",
			link.route.Label, link.transport.DelayServiceResponse, p.read)
	}
	time.Sleep(time.Duration(link.transport.DelayServiceResponse) * time.Millisecond)
	return p, nil
}
