package actions

import (
	"errors"
	"fmt"
	"strings"

	"github.com/eiannone/keyboard"
	"github.com/urfave/cli/v2"
	"gitlab.com/angel-afonso/gitlabcli/api"
	"gitlab.com/angel-afonso/gitlabcli/utils"
	"gopkg.in/gookit/color.v1"
)

type baseProject struct {
	Name              string
	Description       string
	NameWithNamespace string
	ForksCount        int
	StarCount         int
	Visibility        string
	FullPath          string
}

func (project *baseProject) Print() {
	fmt.Println(color.Bold.Sprintf(project.NameWithNamespace))
	fmt.Println(project.Name)
	fmt.Println(project.FullPath)

	if project.Description != "" {
		color.OpItalic.Println(project.Description)
	}

	color.OpItalic.Printf("Stars: %d Forks: %d\n", project.StarCount, project.ForksCount)
	color.OpItalic.Println(project.Visibility)
	println()
}

// Project struct
type Project struct {
	baseProject     `graphql:"inner"`
	CreatedAt       string
	OpenIssuesCount int
	SSHURLToRepo    string `graphql-bind:"sshUrlToRepo"`
	HTTPURLToRepo   string `graphql-bind:"httpUrlToRepo"`
	WebURL          string `graphql-bind:"webUrl"`
	Releases        struct {
		Nodes []struct {
			Name string
		}
	} `graphql:"(first: 1)"`
	Pipelines struct {
		Nodes []struct {
			DetailedStatus struct {
				Label string
			}
		}
	} `graphql:"(first: 1, ref: \\\"master\\\")"`
}

// Print project data
func (project *Project) Print() {
	color.White.Println(color.Bold.Sprintf(project.NameWithNamespace))
	fmt.Println(project.Name)
	color.OpItalic.Println(project.Description)
	color.OpItalic.Println(project.Visibility)
	println()
	color.OpUnderscore.Printf("Stars: %d Forks: %d\n", project.StarCount, project.ForksCount)

	if len(project.Pipelines.Nodes) > 0 {
		println()
		color.OpItalic.Printf("Pipeline status: %s \n", color.Bold.Sprint(project.Pipelines.Nodes[0].DetailedStatus.Label))
		println()
	}

	if len(project.Releases.Nodes) > 0 {
		color.OpItalic.Printf("Last Release: %s \n", project.Releases.Nodes[0].Name)
		println()
	}

	fmt.Printf("Open issues: %d\n", project.OpenIssuesCount)
	println()
	fmt.Printf("HTTP URL: %s\n", color.OpItalic.Sprint(project.HTTPURLToRepo))
	fmt.Printf("SSH URL: %s\n", color.OpItalic.Sprint(project.SSHURLToRepo))
	println()
	color.OpItalic.Println(color.Gray.Sprint(project.WebURL))
}

// ProjectList send request to get user's project
// and print a table with projects
func ProjectList(client *api.Client) func(*cli.Context) error {
	return func(context *cli.Context) error {
		spinner := utils.ShowSpinner()

		var query struct {
			Projects struct {
				PageInfo struct {
					EndCursor   string
					HasNextPage bool
				}
				Nodes []baseProject
			} `graphql:"(membership: true, first: 10,after: $after)"`
		}

		variables := struct {
			after string
		}{
			after: "",
		}

		if err := client.Query(&query, variables); err != nil {
			return err
		}

		for {
			spinner.Stop()
			for _, project := range query.Projects.Nodes {
				project.Print()
			}

			if !query.Projects.PageInfo.HasNextPage {
				return nil
			}

			if err := keyboard.Open(); err != nil {
				return err
			}

			defer keyboard.Close()

			for {
				char, key, _ := keyboard.GetKey()

				if key == keyboard.KeyEnter {
					break
				}

				if char == 'q' || key == keyboard.KeyCtrlC || key == keyboard.KeyEsc {
					println()
					return nil
				}
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
		path, err := utils.GetPathParam(context)

		if err != nil {
			return err
		}

		spinner := utils.ShowSpinner()

		var query struct {
			Project *Project `graphql:"(fullPath:$path)"`
		}

		variables := struct {
			path string `graphql-type:"ID!"`
		}{
			path,
		}

		if err := client.Query(&query, variables); err != nil {
			return err
		}

		spinner.Stop()

		if query.Project != nil {
			query.Project.Print()
			return nil
		}

		return errors.New("An error has occurred, check the repository path and permissions")
	}
}

// ProjectMembers show project members
func ProjectMembers(client *api.Client) func(*cli.Context) error {
	return func(context *cli.Context) error {
		path, err := utils.GetPathParam(context)

		if err != nil {
			return err
		}

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
