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

	// Write in chunks paced by a ticker. Aim for ~1ms granularity so we don't
	// take a syscall per byte; fall back to per-byte pacing only when the
	// configured rate is below 1000 B/s.
	const targetTicksPerSecond = 1000
	chunk := bytesPerSecond / targetTicksPerSecond
	if chunk < 1 {
		chunk = 1
	}
	tickInterval := time.Second / time.Duration(bytesPerSecond/chunk)
	ticker := time.NewTicker(tickInterval)
	defer ticker.Stop()

	data := p.data[:p.read]
	for offset := 0; offset < len(data); {
		<-ticker.C
		end := offset + chunk
		if end > len(data) {
			end = len(data)
		}
		n, err := w.Write(data[offset:end])
		offset += n
		if err != nil {
			return p, err
		}
	}
	return p, nil
}
