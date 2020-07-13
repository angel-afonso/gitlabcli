package main

import (
	"fmt"
	"os"
	"path"

	"github.com/gookit/color"
	cli "github.com/urfave/cli/v2"
	"gitlab.com/angel-afonso/gitlabcli/actions"
	"gitlab.com/angel-afonso/gitlabcli/api"
	"gitlab.com/angel-afonso/gitlabcli/auth"
)

func main() {
	client := api.NewClient(auth.OpenSession())

	fmt.Println()

	app := &cli.App{
		Name:        "gitlabcli",
		Usage:       "Gitlab CLI",
		Version:     "0.1.2",
		Description: "Command line interface to interact with the gitlab API",
		Commands: []*cli.Command{
			{
				Name:        "logout",
				Description: "Remove current session",
				Usage:       "Remove session",
				UsageText:   "gitlabcli logout",
				Action: func(context *cli.Context) error {
					homeDir, _ := os.UserHomeDir()
					sessionDir := path.Join(homeDir, ".gitlabcli", "session")

					if _, err := os.Stat(sessionDir); os.IsNotExist(err) {
						color.Red.Printf("Session does not exist\n")
						return nil
					}

					if err := os.Remove(sessionDir); err != nil {
						return err
					}
					color.Success.Println("Logged out")
					return nil
				},
			},
			{
				Name:        "project",
				Usage:       "Handle Gitlab project",
				Description: "Project related commands",
				Subcommands: []*cli.Command{
					{
						Name:        "list",
						UsageText:   "gitlabcli project list",
						Usage:       "List projects",
						Description: "Display a list with user's project",
						Action:      actions.ProjectList(&client),
					},
					{
						Name:        "view",
						Description: "View project",
						Usage:       "project view [path]",
						Action:      actions.ProjectView(&client),
					},
					{
						Name:        "members",
						Description: "View project members",
						Usage:       "project members [path]",
						Action:      actions.ProjectMembers(&client),
					},
				},
			},
			{
				Name:        "mergerequest",
				Usage:       "Handle merge request",
				Description: "Merge Request related commands",
				Subcommands: []*cli.Command{
					{
						Name:        "create",
						Usage:       "Create new merge request",
						Description: "Create new merge request. Path is optional if the current directory is a git repository with remote in gitlab",
						UsageText:   "gitlabcli mergerequest create [path]",
						Action:      actions.CreateMergeRequest(&client),
					},
					{
						Name:        "assign",
						Description: "Assign user to existing merge request. Path is optional if the current directory is a git repository with remote in gitlab",
						Usage:       "Assign user to a merge request",
						UsageText:   "gitlabcli mergerequest assign [path] <iid>",
						Action:      actions.AssignMergeRequest(&client),
					},
				},
			},
			{
				Name:        "issue",
				Usage:       "Handle project issues",
				Description: "Issues related commands",
				Subcommands: []*cli.Command{
					{
						Name:        "list",
						Usage:       "List project issues",
						Description: "Display a issue list. Path is optional if the current directory is a git repository with remote in gitlab",
						UsageText:   "gitlabcli issue list [path]",
						Action:      actions.IssuesList(&client),
						Flags: []cli.Flag{
							&cli.BoolFlag{
								Name:        "opened",
								Value:       false,
								Usage:       "Display only opened issues",
								Aliases:     []string{"o"},
								Destination: &actions.Opened,
							},
							&cli.BoolFlag{
								Name:        "cloed",
								Value:       false,
								Usage:       "Display only closed issues",
								Aliases:     []string{"c"},
								Destination: &actions.Closed,
							},
						},
					},
					{
						Name:        "view",
						Usage:       "View project issue",
						Description: "Display a specific issue. Path is optional if the current directory is a git repository with remote in gitlab",
						UsageText:   "gitlabcli issue view [path] <iid>",
						Action:      actions.ShowIssue(&client),
					},
				},
			},
		},
	}

	app.EnableBashCompletion = true
	err := app.Run(os.Args)

	if err != nil {
		color.Red.Println(err.Error())
	}

	println()
}
