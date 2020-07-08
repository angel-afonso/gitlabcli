package actions

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/eiannone/keyboard"
	"github.com/urfave/cli/v2"
	"gitlab.com/angel-afonso/gitlabcli/api"
	"gitlab.com/angel-afonso/gitlabcli/utils"
	"gopkg.in/gookit/color.v1"
)

type projectList struct {
	ID          string `graphql-bind:"id"`
	Name        string
	Description string
	FullPath    string
	StarCount   int
}

func (project *projectList) Print() {
	reflVal := reflect.ValueOf(project).Elem()
	reflType := reflect.TypeOf(project).Elem()

	for i := 0; i < reflType.NumField(); i++ {
		color.Cyan.Printf("%s: ", reflType.Field(i).Name)
		color.Reset()
		fmt.Println(reflVal.Field(i))
	}
}

// Project struct
type Project struct {
	projectList `graphql:"inner"`
	ForksCount  int
	Visibility  string
	CreatedAt   string
}

// Print project data
func (project *Project) Print() {
	color.White.Printf("%s\n", color.Bold.Sprintf(project.FullPath))
	fmt.Printf("%s\n", project.Description)
}

// ProjectList send request to get user's project
// and print a table with projects
func ProjectList(client *api.Client) func(*cli.Context) error {
	return func(context *cli.Context) error {
		spinner := utils.ShowSpinner()

		var query struct {
			Projects struct {
				PageInfo struct {
					EndCursor string
				}
				Nodes []projectList
			} `graphql:"(membership: true, first: 10,after: $after)"`
		}

		variables := struct {
			after string
		}{
			after: "",
		}

		client.Query(&query, variables)

		for {
			spinner.Stop()
			if len(query.Projects.Nodes) == 0 {
				return nil
			}

			for _, project := range query.Projects.Nodes {
				project.Print()
				println()
			}

			if err := keyboard.Open(); err != nil {
				panic(err)
			}

			defer keyboard.Close()

			char, key, _ := keyboard.GetKey()

			if char == 'q' || key == keyboard.KeyCtrlC || key == keyboard.KeyEsc {
				println()
				return nil
			}

			println()

			variables.after = query.Projects.PageInfo.EndCursor
			spinner.Start()

			client.Query(&query, variables)
		}
	}
}

// ProjectView get and show data from a project by path
func ProjectView(client *api.Client) func(*cli.Context) error {
	return func(context *cli.Context) error {
		path := utils.GetPathParam(context)

		spinner := utils.ShowSpinner()
		var query struct {
			Project Project `graphql:"(fullPath:$path)"`
		}

		variables := struct {
			path string `graphql-type:"ID!"`
		}{
			path,
		}

		client.Query(&query, variables)

		spinner.Stop()

		query.Project.Print()

		return nil
	}
}

// ProjectMembers show project members
func ProjectMembers(client *api.Client) func(*cli.Context) error {
	return func(context *cli.Context) error {
		path := utils.GetPathParam(context)
		spinner := utils.ShowSpinner()

		users := getProjectMembers(client, path)

		spinner.Stop()
		for _, user := range users {
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
