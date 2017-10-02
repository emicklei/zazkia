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
	"encoding/json"
	"fmt"
	"os"
)

// Route represents the specification of a listener and client for connections (link).
type Route struct {
	Label           string          `json:"label"`
	ServiceHostname string          `json:"service-hostname"`
	ServicePort     int             `json:"service-port"`
	ListenPort      int             `json:"listen-port"`
	Transport       *TransportState `json:"transport"`
}

func readRoutes(location string) (routes []*Route, err error) {
	f, err := os.Open(location)
	if err != nil {
		return
	}
	defer f.Close()
	err = json.NewDecoder(f).Decode(&routes)
	// make sure transport is initialized
	for _, each := range routes {
		if each.Transport == nil {
			each.Transport = newDefaultTransportState()
		}
	}
	return
}

func (r *Route) tcp() string {
	return fmt.Sprintf("%s:%d", r.ServiceHostname, r.ServicePort)
}

func (r *Route) String() string {
	return fmt.Sprintf("[%s] :%d <-> %s:%d {%v}", r.Label, r.ListenPort, r.ServiceHostname, r.ServicePort, r.Transport)
}
