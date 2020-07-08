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
		Version: "0.0.1",
		Commands: []*cli.Command{
			{
				Name:        "logout",
				Description: "Remove current session",
				Usage:       "logout",
				Action: func(context *cli.Context) error {
					homeDir, _ := os.UserHomeDir()
					sessionDir := path.Join(homeDir, ".gitlabcli", "session")

					if _, err := os.Stat(sessionDir); os.IsNotExist(err) {
						color.Red.Printf("Session does not exist\n")
						return nil
					}

					if err := os.Remove(sessionDir); err != nil {
						color.Error.Println(err.Error)
						os.Exit(1)
					}
					color.Success.Println("Logged out")
					return nil
				},
			},
			{
				Name:        "project",
				Description: "Project related commands",
				Subcommands: []*cli.Command{
					{
						Name:        "list",
						Description: "List projects",
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
				Description: "Merge Request related commands",
				Subcommands: []*cli.Command{
					{
						Name:        "create",
						Description: "Create new merge request",
						Usage:       "mergerequest create [path]",
						Action:      actions.CreateMergeRequest(&client),
					},
					{
						Name:        "assign",
						Description: "Assign user to existing merge request",
						Usage:       "mergerequest assign [path] <iid>",
						Action:      actions.AssignMergeRequest(&client),
					},
				},
			},
		},
	}

	app.EnableBashCompletion = true
	err := app.Run(os.Args)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	println()
}
