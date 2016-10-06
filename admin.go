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

	// poor mans routing

	switch {
	case cmd == "/links":
		w.Header().Set("Content-Type", "text/html")
		adminPage.Execute(w, linkMgr.sortedLinks())

	case strings.HasPrefix(cmd, "/delay-response/"):
		token := cmd[len("/delay-response/"):]
		id, err := strconv.Atoi(token)
		if err != nil {
			http.Redirect(w, r, "/links", 307)
			return
		}
		link := linkMgr.get(id)
		if link != nil {
			// toggle
			if link.transport.DelayServiceResponse > 0 {
				link.transport.DelayServiceResponse = 0
			} else {
				link.transport.DelayServiceResponse = 10000 // 10 sec
			}
		} else {
			w.WriteHeader(404)
			return
		}
		http.Redirect(w, r, "/links", 307)

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
			link.transport.SendingToClient = !link.transport.SendingToClient
			link.transport.SendingToService = !link.transport.SendingToService
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
			link.transport.ReceivingFromClient = !link.transport.ReceivingFromClient
			link.transport.ReceivingFromService = !link.transport.ReceivingFromService
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
			<li><a href="/delay-response/{{ .ID }}">delay-response</a></li>
		</ul>
	</li>
{{ end }}
</ul>
</body>
</html>
`))
