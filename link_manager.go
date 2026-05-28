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
	"cmp"
	"errors"
	"log"
	"slices"
	"sync"
)

type linkManager struct {
	links map[int]*link
	mutex sync.RWMutex
}

func newLinkManager() *linkManager {
	return &linkManager{
		links: map[int]*link{},
	}
}

func (m *linkManager) add(l *link) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.links[l.ID] = l
}

func (m *linkManager) get(id int) *link {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.links[id]
}

func (m *linkManager) disconnectAndRemove(id int) error {
	m.mutex.Lock()
	link, ok := m.links[id]
	if !ok {
		m.mutex.Unlock()
		return errors.New("no link with id")
	}
	delete(m.links, link.ID)
	m.mutex.Unlock()
	log.Printf("[%s] disconnecting link %d\n", link.route.Label, link.ID)
	link.disconnect()
	return nil
}

// snapshotIDs returns the current set of link ids filtered by keep.
// If keep is nil, all ids are returned.
func (m *linkManager) snapshotIDs(keep func(*link) bool) []int {
	m.mutex.RLock()
	ids := make([]int, 0, len(m.links))
	for id, l := range m.links {
		if keep == nil || keep(l) {
			ids = append(ids, id)
		}
	}
	m.mutex.RUnlock()
	return ids
}

func (m *linkManager) close() {
	for _, id := range m.snapshotIDs(nil) {
		m.disconnectAndRemove(id)
	}
}

func (m *linkManager) closeAllWithError() {
	ids := m.snapshotIDs(func(l *link) bool {
		return len(l.clientErrorString()) != 0 || len(l.serviceErrorString()) > 0
	})
	for _, id := range ids {
		m.disconnectAndRemove(id)
	}
}

func (m *linkManager) closeAllForRoute(label string) {
	ids := m.snapshotIDs(func(l *link) bool {
		return l.route.Label == label
	})
	for _, id := range ids {
		m.disconnectAndRemove(id)
	}
}

func (m *linkManager) APIGroups() []*APILinkGroup {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	gm := map[string]*APILinkGroup{}
	// TEMP
	for _, each := range routeMgr.routes {
		_, ok := gm[each.Label]
		if !ok {
			g := new(APILinkGroup)
			g.Route = each
			gm[each.Label] = g
		}
	}
	for _, each := range m.links {
		g := gm[each.route.Label]
		g.Links = append(g.Links, NewAPILink(each))
	}
	all := []*APILinkGroup{}
	for _, each := range gm {
		slices.SortFunc(each.Links, func(a, b APILink) int { return a.ID - b.ID })
		all = append(all, each)
	}
	slices.SortFunc(all, func(a, b *APILinkGroup) int { return cmp.Compare(a.Route.Label, b.Route.Label) })
	return all
}

func (m *linkManager) APILinks(routeLabel string) (all []APILink) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	for _, each := range m.links {
		if each.route.Label == routeLabel {
			all = append(all, NewAPILink(each))
		}
	}
	slices.SortFunc(all, func(a, b APILink) int { return a.ID - b.ID })
	return
}
