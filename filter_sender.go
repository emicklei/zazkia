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

// sender will write all the data of a parcel to a writer.
type sender struct{}

func (s sender) Write(link *link, w io.Writer, p parcel) (parcel, error) {
	if p.isEmpty() {
		return p, nil
	}
	offset := 0
	towrite := p.read
	for towrite > 0 {
		var (
			err     error
			written int
		)
		subset := p.data[offset:p.read]
		written, err = w.Write(subset)
		if err != nil {
			return p, err
		}
		offset += written
		towrite -= written
	}
	return parcel{p.data, p.read, p.read}, nil
}
