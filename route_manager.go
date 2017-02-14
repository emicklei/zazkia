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

import "errors"

type routeManager struct {
	routes []*Route
}

func (m routeManager) routeForRemoteAddress(addr string) (*Route, error) {
	for _, each := range m.routes {
		if each.ServiceHostname == addr {
			return each, nil
		}
	}
	return new(Route), errors.New("no route for address:" + addr)
}

func (m routeManager) routeForLabel(label string) *Route {
	for _, each := range m.routes {
		if each.Label == label {
			return each
		}
	}
	return nil
}
