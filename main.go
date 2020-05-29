package main

import (
	"github.com/angel-afonso/gitlabcli/auth"
	"github.com/angel-afonso/gitlabcli/graphql"
)

func main() {
	client := graphql.NewClient(auth.OpenSession())

	var query struct {
		Projects struct {
			Nodes []struct {
				Name string
			}
		} `graphql:"(membership: true)"`
	}

	client.Query(&query)
}
