package actions

import (
	"fmt"
	"log"

	"github.com/urfave/cli/v2"
	"gitlab.com/angel-afonso/gitlabcli/graphql"
	"gitlab.com/angel-afonso/gitlabcli/utils"
)

// ProjectList send request to get user's project
// and print a table with projects
func ProjectList(client *api.Client) func(*cli.Context) error {
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
func ProjectView(client *api.Client) func(*cli.Context) error {
	return func(context *cli.Context) error {
		var path string

		if utils.IsGitRepository() && context.Args().Len() == 0 {
			remotes := utils.GetRemote()
			if len(remotes) > 1 {
				path = utils.GetRemotePath(utils.AskRemote(remotes))
			}
			path = utils.GetRemotePath(remotes[0])
		} else if context.Args().Len() > 0 {
			path = context.Args().First()
		} else {
			log.Fatal("Expected project path")
		}

		var query struct {
			Project struct {
				ID          string `graphql-bind:"id"`
				Description string
				Name        string
			} `graphql:"(fullPath:$path)"`
		}

		variables := struct {
			path string `graphql-type:"ID!"`
		}{
			path,
		}

		client.Query(&query, variables)

		fmt.Printf("ID: %s\nName: %s\nDescription: %s\n\n", query.Project.ID, query.Project.Name, query.Project.Description)

		return nil
	}
}
