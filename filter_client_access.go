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

// clientAccess is a filter to control whether reading or writing a client is allowed for the link.
type clientAccess struct{}

func (c clientAccess) Read(link *link, r io.Reader, p parcel) (parcel, error) {
	if !link.transport.ReceivingFromClient {
		return emptyParcel, nil
	}
	return receiver{}.Read(link, r)
}

func (c clientAccess) Write(link *link, w io.Writer, p parcel) (parcel, error) {
	if p.isEmpty() {
		return p, nil
	}
	if !link.transport.SendingToClient {
		return emptyParcel, nil
	}
	return p, nil
}
