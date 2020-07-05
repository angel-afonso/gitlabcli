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
		path := utils.GetPathParam(context)

		color.Cyan.Print("Merge request title: ")
		title := utils.ReadLine()
		color.Cyan.Print("Source Branch: ")
		source := utils.ReadLine()
		color.Cyan.Print("Target Branch: ")
		target := utils.ReadLine()
		color.Cyan.Print("Description: ")
		description := utils.ReadLine()

		spinner := utils.ShowSpinner()

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

		spinner.Stop()

		color.Green.Printf("Created merge request !%s\n", mutation.MergeRequestCreate.MergeRequest.Iid)
		color.Reset()

		fmt.Print("Assign merge request? ")
		color.Blue.Print("y/n ")
		color.FgGray.Print("default: n")
		color.Reset()

		if choice := utils.ReadLine(); choice == "y" || choice == "yes" {
			spinner.Start()

			users := getProjectMembers(client, path)

			spinner.Stop()

			for index, user := range users {
				color.Blue.Printf("%d ", index+1)
				color.Green.Printf("%s (%s)\n", user.Name, user.Username)
			}

			index := utils.ReadInt()

			assignUserForMergeRequest(client,
				mutation.MergeRequestCreate.MergeRequest.Iid,
				path,
				[]string{`"` + users[index-1].Username + `"`},
			)
		}

		return nil
	}
}

// AssignMergeRequest interact with the graphql api to assign user to merge request
func AssignMergeRequest(client *api.Client) func(*cli.Context) error {
	return func(context *cli.Context) error {
		spinner := utils.ShowSpinner()

		path := utils.GetPathParam(context)

		args := context.Args()

		var iid string

		if args.Len() > 1 {
			iid = args.Get(1)
		} else {
			iid = args.Get(0)
		}

		if iid == "" {
			log.Fatal("iid is required")
		}

		users := getProjectMembers(client, path)

		spinner.Stop()

		for index, user := range users {
			color.Blue.Printf("%d ", index+1)
			color.Green.Printf("%s (%s)\n", user.Name, user.Username)
		}

		index := utils.ReadInt()

		assignUserForMergeRequest(client,
			iid,
			path,
			[]string{`"` + users[index-1].Username + `"`},
		)
		return nil
	}
}

func assignUserForMergeRequest(client *api.Client, iid string, path string, usernames []string) {
	spinner := utils.ShowSpinner()

	var assignMutation struct {
		MergeRequestSetAssignees struct {
			MergeRequest struct {
				iid string
			}
		} `graphql:"(input:{projectPath:$path,iid:$iid,assigneeUsernames:$usernames})"`
	}

	assignVariables := struct {
		path      string   `graphql-type:"ID!"`
		iid       string   `graphql-type:"String!"`
		usernames []string `graphql-type:"[String!]!"`
	}{
		path:      path,
		iid:       iid,
		usernames: usernames,
	}

	client.Mutation(&assignMutation, assignVariables)

	spinner.Stop()

	color.Green.Printf("%s assigned to merge request !%s\n", usernames[0], iid)
}
