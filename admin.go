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
			<li><a href="/v1/links/{{ .ID }}/close">close</a></li>
			<li><a href="/v1/links/{{ .ID }}/delay-response">delay-response</a></li>
			<li><a href="/v1/links/{{ .ID }}/toggle-send">toggle-send</a></li>
			<li><a href="/v1/links/{{ .ID }}/toggle-receive">toggle-receive</a></li>			
			<li><a href="/v1/links/{{ .ID }}/toggle-verbose">toggle-verbose</a></li>
		</ul>
	</li>
{{ end }}
</ul>
</body>
</html>
`))
