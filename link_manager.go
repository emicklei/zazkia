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
	"errors"
	"log"
	"sort"
	"sync"
)

type linkManager struct {
	links map[int]*link
	mutex *sync.Mutex
}

func newLinkManager() *linkManager {
	return &linkManager{
		links: map[int]*link{},
		mutex: new(sync.Mutex),
	}
}

func (m *linkManager) add(l *link) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.links[l.ID] = l
}

func (m *linkManager) remove(l *link) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	delete(m.links, l.ID)
}

func (m *linkManager) get(id int) *link {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	l, ok := m.links[id]
	if !ok {
		return nil
	}
	return l
}

func (m *linkManager) disconnectAndRemove(id int) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	link, ok := m.links[id]
	if !ok {
		return errors.New("no link with id")
	}
	log.Printf("[%s] disconnecting link %d\n", link.route.Label, link.ID)
	link.disconnect()
	delete(m.links, link.ID)
	return nil
}

func (m *linkManager) close() {
	for _, each := range m.links {
		m.disconnectAndRemove(each.ID)
	}
}

func (m *linkManager) closeAllWithError() {
	for _, each := range m.links {
		if len(each.clientErrorString()) != 0 || len(each.serviceErrorString()) > 0 {
			m.disconnectAndRemove(each.ID)
		}
	}
}

func (m *linkManager) APIGroups() []*APILinkGroup {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	gm := map[string]*APILinkGroup{}
	// TEMP
	for _, each := range routeMgr.routes {
		g, ok := gm[each.Label]
		if !ok {
			g = new(APILinkGroup)
			g.Route = each
			gm[each.Label] = g
		}
	}
	for _, each := range m.links {
		g, _ := gm[each.route.Label]
		g.Links = append(g.Links, NewAPILink(each))
	}
	all := []*APILinkGroup{}
	for _, each := range gm {
		sort.Sort(APILinkSorter(each.Links))
		all = append(all, each)
	}
	sort.Sort(APILinkGroupSorter(all))
	return all
}

func (m *linkManager) APILinks(routeLabel string) (all []APILink) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	for _, each := range m.links {
		if each.route.Label == routeLabel {
			all = append(all, NewAPILink(each))
		}
	}
	sort.Sort(APILinkSorter(all))
	return
}

type APILinkGroupSorter []*APILinkGroup

func (s APILinkGroupSorter) Len() int {
	return len(s)
}
func (s APILinkGroupSorter) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s APILinkGroupSorter) Less(i, j int) bool {
	a, b := s[i], s[j]
	return a.Route.Label < b.Route.Label
}

type APILinkSorter []APILink

func (s APILinkSorter) Len() int {
	return len(s)
}
func (s APILinkSorter) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s APILinkSorter) Less(i, j int) bool {
	a, b := s[i], s[j]
	return a.ID < b.ID
}
