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
	"math/rand"
)

// corrupter mangles the byte sequence for writing
type corrupter struct{}

func (c corrupter) Write(link *link, w io.Writer, p parcel) (parcel, error) {
	m := link.transport.ServiceResponseCorruptMethod
	if len(m) == 0 {
		return p, nil
	}
	if p.read == 0 {
		return p, nil
	}
	switch m {
	case "randomize":
		for i := 0; i < p.read; i++ {
			j := rand.Intn(i + 1)
			p.data[i], p.data[j] = p.data[j], p.data[i]
		}
	}
	return p, nil
}
