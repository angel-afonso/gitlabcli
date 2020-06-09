package main

import (
	"log"
	"os"

	cli "github.com/urfave/cli/v2"
	"gitlab.com/angel-afonso/gitlabcli/actions"
	"gitlab.com/angel-afonso/gitlabcli/api"
	"gitlab.com/angel-afonso/gitlabcli/auth"
)

func main() {
	client := api.NewClient(auth.OpenSession())

	app := &cli.App{
		Version: "0.0.1",
		Commands: []*cli.Command{
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
						Usage:       "project view <path>",
						Action:      actions.ProjectView(&client),
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
						Usage:       "mergerequest create [project path]",
						Action:      actions.CreateMergeRequest(&client),
					},
				},
			},
		},
	}

	app.EnableBashCompletion = true

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
