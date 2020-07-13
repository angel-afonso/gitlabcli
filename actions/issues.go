package actions

import (
	"errors"
	"fmt"

	"github.com/eiannone/keyboard"
	cli "github.com/urfave/cli/v2"
	"gitlab.com/angel-afonso/gitlabcli/api"
	"gitlab.com/angel-afonso/gitlabcli/utils"
	"gopkg.in/gookit/color.v1"
)

var (
	// Opened store flag --opened valie
	Opened bool
	// Closed store flag --closed valie
	Closed bool
)

type issuesList struct {
	Iid       string
	Title     string
	State     string
	Assignees struct {
		Nodes []User
	}
}

func (il *issuesList) Print() {
	color.Bold.Println(il.Title)
	fmt.Println(il.State)
	color.OpItalic.Printf("#%s\n", il.Iid)
	fmt.Println("Assignees:")

	for _, user := range il.Assignees.Nodes {
		fmt.Print("  - ")
		user.Print()
	}

	println()
}

// Issue struct representation
type Issue struct {
	issuesList `graphql:"inner"`
	Author     User
}

// Print issue
func (i *Issue) Print() {
	color.Bold.Println(i.Title)
	fmt.Println(i.State)
	println()
	color.OpItalic.Printf("#%s\n", i.Iid)
	println()
	fmt.Print("Author: ")
	i.Author.Print()
	println()
	fmt.Println("Assignees:")
	println()

	for _, user := range i.Assignees.Nodes {
		fmt.Print("  - ")
		user.Print()
	}
}

// IssuesList display a project's issues list
func IssuesList(client *api.Client) func(*cli.Context) error {
	return func(context *cli.Context) error {
		path, err := utils.GetPathParam(context)

		if err != nil {
			return err
		}

		spinner := utils.ShowSpinner()

		var query struct {
			Project *struct {
				Issues struct {
					PageInfo struct {
						EndCursor   string
						HasNextPage bool
					}
					Nodes []issuesList
				} `graphql:"(first: 10, after: $after)"`
			} `graphql:"(fullPath:$path)"`
		}

		variables := struct {
			path  string `graphql-type:"ID!"`
			after string
		}{
			path:  path,
			after: "",
		}

		if err := client.Query(&query, variables); err != nil {
			return err
		}

		for {
			spinner.Stop()

			if query.Project == nil {
				return errors.New("An error has occurred, check the repository path and permissions")
			}

			for _, issue := range query.Project.Issues.Nodes {
				issue.Print()
			}

			if !query.Project.Issues.PageInfo.HasNextPage {
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

			variables.after = query.Project.Issues.PageInfo.EndCursor
			spinner.Start()

			client.Query(&query, variables)
		}
	}
}

// ShowIssue search and display a issue by iid
func ShowIssue(client *api.Client) func(*cli.Context) error {
	return func(context *cli.Context) error {
		path, err := utils.GetPathParam(context)

		if err != nil {
			return err
		}

		var iid string

		if context.Args().Len() > 1 {
			iid = context.Args().Get(1)
		} else {
			iid = context.Args().Get(0)
		}

		if iid == "" {
			return errors.New("iid is required")
		}

		spinner := utils.ShowSpinner()

		var query struct {
			Project struct {
				Issue *Issue `graphql:"(iid:$iid)"`
			} `graphql:"(fullPath:$path)"`
		}

		variables := struct {
			path string `graphql-type:"ID!"`
			iid  string
		}{
			path,
			iid,
		}

		if err := client.Query(&query, variables); err != nil {
			return err
		}

		spinner.Stop()

		if query.Project.Issue != nil {
			query.Project.Issue.Print()
			return nil
		}

		return errors.New("An error has occurred, check the repository path and permissions")
	}
}
