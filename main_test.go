/*
Copyright 2017 Ernest Micklei

Licensed under the Apache License, Version 2.0 (the "License");
*/

package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	restful "github.com/emicklei/go-restful/v3"
	"github.com/go-openapi/spec"
)

func TestMIMESetter_ServeHTTP(t *testing.T) {
	h := MIMESetter{Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})}
	for _, p := range []string{"/x.css", "/x.js", "/x.png"} {
		w := httptest.NewRecorder()
		r, _ := http.NewRequest("GET", p, nil)
		h.ServeHTTP(w, r)
		if w.Code != 200 {
			t.Errorf("%s: got %d", p, w.Code)
		}
	}
}

func TestExtendSwaggerObject(t *testing.T) {
	s := &spec.Swagger{}
	extendSwaggerObject(s)
	if s.Info == nil || s.Info.Title == "" {
		t.Error("swagger info not set")
	}
	if len(s.Tags) == 0 {
		t.Error("swagger tags not set")
	}
}

func TestAddSwagger(t *testing.T) {
	restful.DefaultContainer = restful.NewContainer()
	addRESTResources()
	addSwagger()
}
