package actions

import (
	"log"
	"strings"

	"github.com/angel-afonso/gitlabcli/graphql"
	"github.com/angel-afonso/gitlabcli/utils"
	"github.com/gookit/color"
	"github.com/urfave/cli/v2"
)

// CreateMergeRequest send a request to create merge request
// by given project path
func CreateMergeRequest(client *graphql.Client) func(*cli.Context) error {
	return func(context *cli.Context) error {

		if context.Args().Len() < 1 {
			log.Fatal("Expected project path")
		}

		path := context.Args().First()
		color.Cyan.Print("Merge request title: ")
		title := utils.ReadLine()
		color.Cyan.Print("Source Branch: ")
		source := utils.ReadLine()
		color.Cyan.Print("Target Branch: ")
		target := utils.ReadLine()
		color.Cyan.Print("Description: ")
		description := utils.ReadLine()

		var mutation struct {
			MergeRequestCreate struct {
				MergeRequest struct {
					Iid string
				}
				Errors []string
			} `graphql:"(input:{title:$title,projectPath:$path,sourceBranch:$source,targetBranch:$target,description:$description})"`
		}

		variables := struct {
			path        string `graphql-type:"ID!"`
			title       string `graphql-type:"String!"`
			source      string `graphql-type:"String!"`
			target      string `graphql-type:"String!"`
			description string
		}{
			path,
			title,
			source,
			target,
			description,
		}

		client.Mutation(&mutation, variables)

		if len(mutation.MergeRequestCreate.Errors) > 0 {
			color.Red.Println(strings.Join(mutation.MergeRequestCreate.Errors, "\n"))
			return nil
		}

		color.Green.Print("Created merge request ")
		color.BgCyan.Printf("!%s\n", mutation.MergeRequestCreate.MergeRequest.Iid)

		return nil
	}
}
