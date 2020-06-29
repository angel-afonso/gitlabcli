package actions

import (
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"
	"gitlab.com/angel-afonso/gitlabcli/api"
	"gitlab.com/angel-afonso/gitlabcli/utils"
	"gopkg.in/gookit/color.v1"
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
		path := utils.GetPathParam(context)

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

// ProjectMembers show project members
func ProjectMembers(client *api.Client) func(*cli.Context) error {
	return func(context *cli.Context) error {
		path := utils.GetPathParam(context)
		for _, user := range getProjectMembers(client, path) {
			color.FgGreen.Print(user.Name)
			color.FgGreen.Printf(" (%s)\n", user.Username)
			color.Reset()
		}
		color.Reset()
		return nil
	}
}

func getProjectMembers(client *api.Client, path string) []User {
	var users []User

	client.Get(fmt.Sprintf("projects/%s/users", strings.ReplaceAll(path, "/", "%2F")), &users)
	return users
}
