package main

import "errors"

type routeManager struct {
	routes []Route
}

func (m routeManager) routeForRemoteAddress(addr string) (Route, error) {
	for _, each := range m.routes {
		if each.TargetDomain == addr {
			return each, nil
		}
	}
	return Route{}, errors.New("no route for address:" + addr)
}
