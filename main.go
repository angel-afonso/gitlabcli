package main

import (
	"log"
	"os"

	"github.com/angel-afonso/gitlabcli/auth"
	"github.com/angel-afonso/gitlabcli/graphql"
	cli "github.com/urfave/cli/v2"
)

func main() {
	graphql.NewClient(auth.OpenSession())

	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:        "Project",
				Description: "asdads",
			},
		},
	}

	app.EnableBashCompletion = true

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
