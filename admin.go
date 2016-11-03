package main

import "text/template"

var adminPage = template.Must(template.New("admin").Parse(`<html>
<body>
<h1>zazkia - a tcp proxy for simulating errors</h1>
<h3>open links</h3>
<ul>
<!-- apigroups -->
{{ range .}}
	<li>
		<!-- apigroup -->
		<h4>Route {{ .Route }}</h4>
		<ol>
		{{ range .Links }} 
			<h5>{{ . }}</h5>		
			<ul>
				<li><a href="/links/{{ .ID }}/close">close</a></li>
				<li><a href="/links/{{ .ID }}/delay-response">delay-response</a></li>
				<li><a href="/links/{{ .ID }}/toggle-send">toggle-send</a></li>
				<li><a href="/links/{{ .ID }}/toggle-receive">toggle-receive</a></li>			
				<li><a href="/links/{{ .ID }}/toggle-verbose">toggle-verbose</a></li>
			</ul>
		{{ end }}
		</ol>
	</li>
{{ end }}
</ul>
</body>
</html>
`))

type APILinkGroup struct {
	Route Route     `json:"route"`
	Links []APILink `json:"links"`
}

func APIGroups(links []APILink) []APILinkGroup {
	return []APILinkGroup{
		APILinkGroup{
			Links: links,
		},
	}
}
