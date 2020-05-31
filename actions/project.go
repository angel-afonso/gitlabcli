package actions

import (
	"fmt"

	"github.com/angel-afonso/gitlabcli/graphql"
	"github.com/urfave/cli/v2"
)

// ProjectList send request to get user's project
// and print a table with projects
func ProjectList(client *graphql.Client) func(*cli.Context) error {
	return func(context *cli.Context) error {
		var query struct {
			Projects struct {
				Nodes []struct {
					ID          string `graphql-bind:"id"`
					Name        string
					Description string
					FullPath    string
				}
			} `graphql:"(membership: true)"`
		}

		client.Query(&query, nil)

		for _, project := range query.Projects.Nodes {
			fmt.Printf("ID: %s\nName: %s\nDescription: %s\nPath: %s\n\n", project.ID, project.Name, project.Description, project.FullPath)
		}

		return nil
	}
}

// ProjectView get and show data from a project by path
func ProjectView(client *graphql.Client) func(*cli.Context) error {
	return func(context *cli.Context) error {
		var query struct {
			Project struct {
				ID          string `graphql-bind:"id"`
				Description string
				Name        string
			} `graphql:"(fullPath:$path)"`
		}

		variables := struct {
			Path string `graphql-type:"ID!"`
		}{
			Path: context.Args().Get(0),
		}

		client.Query(&query, variables)

		fmt.Printf("ID: %s\nName: %s\nDescription: %s\n\n", query.Project.ID, query.Project.Name, query.Project.Description)

		return nil
	}
}
