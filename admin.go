package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"text/template"
)

func commandHandler(w http.ResponseWriter, r *http.Request) {
	cmd := r.URL.String()
	switch {
	case cmd == "/links":
		w.Header().Set("Content-Type", "text/html")
		adminPage.Execute(w, linkMgr.sortedLinks())

	case strings.HasPrefix(cmd, "/close/"):
		token := cmd[len("/close/"):]
		id, err := strconv.Atoi(token)
		if err != nil {
			http.Redirect(w, r, "/links", 307)
			return
		}
		linkMgr.disconnectAndRemove(id)
		http.Redirect(w, r, "/links", 307)

	case strings.HasPrefix(cmd, "/toggle-send/"):
		token := cmd[len("/toggle-send/"):]
		id, err := strconv.Atoi(token)
		if err != nil {
			w.WriteHeader(400)
			return
		}
		link := linkMgr.get(id)
		if link != nil {
			link.sendingToClient = !link.sendingToClient
			link.sendingToService = !link.sendingToService
		} else {
			w.WriteHeader(404)
			return
		}
		http.Redirect(w, r, "/links", 307)

	case strings.HasPrefix(cmd, "/toggle-receive/"):
		token := cmd[len("/toggle-receive/"):]
		id, err := strconv.Atoi(token)
		if err != nil {
			w.WriteHeader(400)
			return
		}
		link := linkMgr.get(id)
		if link != nil {
			link.receivingFromClient = !link.receivingFromClient
			link.receivingFromService = !link.receivingFromService
		} else {
			w.WriteHeader(404)
			return
		}
		http.Redirect(w, r, "/links", 307)

	case cmd == "/favicon.ico":
	default:
		fmt.Fprintf(w, "sorry, unknown command")
	}
}

var adminPage = template.Must(template.New("admin").Parse(`<html>
<body>
<h1>zazkia - a tcp proxy for simulating errors</h1>
<h3>open links</h3>
<ul>
{{ range .}}
	<li>
		<h4>{{ . }}</h4>
		<ul>
			<li><a href="/close/{{ .ID }}">close</a></li>
			<li><a href="/toggle-send/{{ .ID }}">toggle-send</a></li>
			<li><a href="/toggle-receive/{{ .ID }}">toggle-receive</a></li>
		</ul>
	</li>
{{ end }}
</ul>
</body>
</html>
`))
