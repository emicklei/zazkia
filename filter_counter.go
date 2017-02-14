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

// counter will update the stats of the link.
type counter struct {
	accessService bool
}

func (c counter) Read(link *link, r io.Reader, p parcel) (parcel, error) {
	if c.accessService {
		link.stats.BytesReceivedFromService += int64(p.read)
	} else {
		link.stats.BytesReceivedFromClient += int64(p.read)
	}
	return p, nil

}

func (c counter) Write(link *link, w io.Writer, p parcel) (parcel, error) {
	if c.accessService {
		link.stats.BytesSentToService += int64(p.written)
	} else {
		link.stats.BytesSentToClient += int64(p.written)
	}
	return p, nil
}
