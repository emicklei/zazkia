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
