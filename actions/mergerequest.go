package actions

import (
	"fmt"
	"log"
	"strings"

	"github.com/angel-afonso/gitlabcli/graphql"
	"github.com/angel-afonso/gitlabcli/utils"
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
		title := utils.ReadLine("Merge request title: ")
		source := utils.ReadLine("Source Branch: ")
		target := utils.ReadLine("Target Branch: ")
		description := utils.ReadLine("Description: ")

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
			fmt.Println(strings.Join(mutation.MergeRequestCreate.Errors, "\n"))
			return nil
		}

		fmt.Printf("Created merge request !%s\n", mutation.MergeRequestCreate.MergeRequest.Iid)

		return nil
	}
}
