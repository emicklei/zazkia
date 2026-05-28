/*
Copyright 2017 Ernest Micklei

Licensed under the Apache License, Version 2.0 (the "License");
*/

package main

import (
	"errors"
	"testing"
)

func TestLinkManager_DisconnectAndRemove(t *testing.T) {
	m := newLinkManager()
	r := newTestRoute("r1")
	l, cleanup := newTestLink(t, r)
	defer cleanup()
	m.add(l)
	if err := m.disconnectAndRemove(l.ID); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if err := m.disconnectAndRemove(99999); err == nil {
		t.Error("expected error for missing id")
	}
}

func TestLinkManager_Close(t *testing.T) {
	m := newLinkManager()
	r := newTestRoute("r1")
	l, cleanup := newTestLink(t, r)
	defer cleanup()
	m.add(l)
	done := make(chan struct{})
	go func() {
		m.close()
		close(done)
	}()
	select {
	case <-done:
	case <-mustTimer(t):
		t.Fatal("close timed out")
	}
}

func TestLinkManager_CloseAllForRoute(t *testing.T) {
	m := newLinkManager()
	r1 := newTestRoute("a")
	r2 := newTestRoute("b")
	l1, c1 := newTestLink(t, r1)
	l2, c2 := newTestLink(t, r2)
	defer c1()
	defer c2()
	m.add(l1)
	m.add(l2)

	done := make(chan struct{})
	go func() {
		m.closeAllForRoute("a")
		close(done)
	}()
	select {
	case <-done:
	case <-mustTimer(t):
		t.Fatal("closeAllForRoute timed out")
	}
	if m.get(l2.ID) == nil {
		t.Error("route b link should remain")
	}
}

func TestLinkManager_CloseAllWithError(t *testing.T) {
	m := newLinkManager()
	r := newTestRoute("r1")
	lOk, c1 := newTestLink(t, r)
	lErr, c2 := newTestLink(t, r)
	defer c1()
	defer c2()
	lErr.clientError = errors.New("boom")
	m.add(lOk)
	m.add(lErr)

	done := make(chan struct{})
	go func() {
		m.closeAllWithError()
		close(done)
	}()
	select {
	case <-done:
	case <-mustTimer(t):
		t.Fatal("closeAllWithError timed out")
	}
	if m.get(lOk.ID) == nil {
		t.Error("ok link should remain")
	}
}

func TestLinkManager_APIGroupsAndLinks(t *testing.T) {
	r := newTestRoute("r1")
	routeMgr = routeManager{routes: []*Route{r}}
	m := newLinkManager()
	l, cleanup := newTestLink(t, r)
	defer cleanup()
	m.add(l)

	groups := m.APIGroups()
	if len(groups) != 1 || len(groups[0].Links) != 1 {
		t.Errorf("got %#v", groups)
	}
	links := m.APILinks("r1")
	if len(links) != 1 {
		t.Errorf("got %d links", len(links))
	}
	if other := m.APILinks("other"); len(other) != 0 {
		t.Errorf("expected no links, got %d", len(other))
	}
}
