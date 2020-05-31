package main

import (
	"log"
	"os"

	"github.com/angel-afonso/gitlabcli/actions"
	"github.com/angel-afonso/gitlabcli/auth"
	"github.com/angel-afonso/gitlabcli/graphql"
	cli "github.com/urfave/cli/v2"
)

func main() {
	client := graphql.NewClient(auth.OpenSession())

	app := &cli.App{
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
		},
	}

	app.EnableBashCompletion = true

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
