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

var transportBufferSize = 32 * 1024

type receiver struct{}

func (rc receiver) Read(link *link, r io.Reader) (parcel, error) {
	buffer := make([]byte, transportBufferSize)
	read, err := r.Read(buffer)
	return parcel{buffer, read, 0}, err
}
