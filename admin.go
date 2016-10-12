package main

import "text/template"

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
