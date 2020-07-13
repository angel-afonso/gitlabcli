package actions

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/gookit/color"
	cli "github.com/urfave/cli/v2"
	"gitlab.com/angel-afonso/gitlabcli/api"
	"gitlab.com/angel-afonso/gitlabcli/utils"
)

// CreateMergeRequest send a request to create merge request
// by given project path
func CreateMergeRequest(client *api.Client) func(*cli.Context) error {
	return func(context *cli.Context) error {
		path, err := utils.GetPathParam(context)
		if err != nil {
			return err
		}

		var head *plumbing.Reference
		var commit *object.Commit

		if context.Args().Len() == 0 {
			head = utils.RepoHead()
			commit = utils.RepoLastCommit()
		}

		color.Bold.Print("Merge request title ")
		color.Reset()

		if commit != nil {
			color.LightBlue.Printf(" (Default: %s)", strings.TrimSpace(commit.Message))
		}

		color.White.Print(": ")
		title := utils.ReadLineOptional(
			utils.Ternary(commit != nil, strings.TrimSpace(commit.Message), "").(string),
		)

		color.Bold.Print("Source Branch ")
		color.Reset()

		if head != nil {
			color.LightBlue.Printf(" (Default: %s)", head.Name().Short())
		}

		color.White.Print(": ")
		source := utils.ReadLineOptional(
			utils.Ternary(head != nil, head.Name().Short(), "").(string),
		)

		color.Bold.Print("Target Branch ")
		color.Reset()
		color.LightBlue.Print("(Default: master)")
		color.White.Print(": ")

		target := utils.ReadLineOptional("master")

		color.Bold.Print("Description: ")
		color.Reset()
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

		color.LightGreen.Printf("Created merge request !%s\n", mutation.MergeRequestCreate.MergeRequest.Iid)

		fmt.Print("Assign merge request? ")
		color.Blue.Print("y/n ")
		color.Gray.Print("default: n ")
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
