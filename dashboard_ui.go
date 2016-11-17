package main

type APILinkGroup struct {
	Route *Route    `json:"route"`
	Links []APILink `json:"links"`
}

func APIGroups(links []APILink) []APILinkGroup {
	return []APILinkGroup{
		APILinkGroup{
			Links: links,
		},
	}
}
