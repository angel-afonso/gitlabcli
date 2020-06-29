package actions

import (
	"fmt"
	"log"
	"strings"

	"github.com/gookit/color"
	"github.com/urfave/cli/v2"
	"gitlab.com/angel-afonso/gitlabcli/api"
	"gitlab.com/angel-afonso/gitlabcli/utils"
)

// CreateMergeRequest send a request to create merge request
// by given project path
func CreateMergeRequest(client *api.Client) func(*cli.Context) error {
	return func(context *cli.Context) error {

		var path string

		if utils.IsGitRepository() {
			remotes := utils.GetRemote()
			if len(remotes) > 1 {
				path = utils.GetRemotePath(utils.AskRemote(remotes))
			}
			path = utils.GetRemotePath(remotes[0])

		} else if context.Args().Len() > 0 {
			path = context.Args().First()
		} else {
			log.Fatal("Expected project path")
		}

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

		color.Green.Printf("Created merge request !%s\n", mutation.MergeRequestCreate.MergeRequest.Iid)
		color.Reset()

		fmt.Print("Assign merge request? ")
		color.Blue.Print("y/n")
		color.FgGray.Print("default: n")
		color.Reset()

		if choice := utils.ReadLine(); choice == "y" || choice == "yes" {
			users := getProjectMembers(client, path)
			for index, user := range users {
				color.Blue.Printf("%d ", index+1)
				color.Green.Printf("%s (%s)\n", user.Name, user.Username)
			}

			index := utils.ReadInt()

			var assignMutation struct {
				MergeRequestSetAssignees struct {
					Errors []string
				} `graphql:"(input:{projectPath:$path,iid:$iid,assigneeUsernames:$usernames})"`
			}

			assignVariables := struct {
				path      string   `graphql-type:"ID!"`
				iid       string   `graphql-type:"String!"`
				usernames []string `graphql-type:"[String!]!"`
			}{
				path:      path,
				iid:       mutation.MergeRequestCreate.MergeRequest.Iid,
				usernames: []string{users[index].Username},
			}

			client.Mutation(assignMutation, assignVariables)
		}

		return nil
	}
}
