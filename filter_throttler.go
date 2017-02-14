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
	"time"
)

type throttler struct{}

func (t throttler) Write(link *link, w io.Writer, p parcel) (parcel, error) {
	if p.isEmpty() {
		return p, nil
	}
	if link.transport.ThrottleServiceResponse == 0 {
		return p, nil
	}
	bytesPerSecond := link.transport.ThrottleServiceResponse
	delayMicrosecondsPerByte := 1000 * 1000 / bytesPerSecond
	if delayMicrosecondsPerByte == 0 {
		return sender{}.Write(link, w, p)
	}
	throttler := time.NewTicker(time.Duration(delayMicrosecondsPerByte) * time.Microsecond)
	defer throttler.Stop()
	for _, each := range p.data[:p.read] {
		<-throttler.C
		_, err := w.Write([]byte{each})
		if err != nil {
			return p, err
		}
	}
	return p, nil
}
