package main

import (
	"net/http"
	"text/template"
)

var indexPage *template.Template

func init() {
	t, _ := Asset("dashboard/templates/index.html")
	indexPage = template.Must(template.New("index.html").Parse(string(t)))
}

func dashboardResourceIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	indexPage.Execute(w, linkMgr.APIGroups())
}
