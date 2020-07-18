package actions

import (
	"errors"
	"fmt"
	"strings"

	"github.com/eiannone/keyboard"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/gookit/color"
	cli "github.com/urfave/cli/v2"
	"gitlab.com/angel-afonso/gitlabcli/api"
	"gitlab.com/angel-afonso/gitlabcli/utils"
)

type baseMergeRequest struct {
	Iid    string `graphql-type:"string!"`
	Title  string
	State  string
	Author User
}

func (mr *baseMergeRequest) Print() {
	color.OpItalic.Printf("!%s\n", mr.Iid)
	color.Bold.Println(mr.Title)
	fmt.Printf("State: %s\n", mr.State)
	mr.Author.Print()
	println()
}

// MergeRequest graphql struct
type MergeRequest struct {
	baseMergeRequest `graphql:"inner"`
	Description      string
	SourceBranch     string
	TargetBranch     string
	Assignees        struct {
		Nodes []User
	}
	Author User
}

// Print merge request
func (mr *MergeRequest) Print() {
	color.OpItalic.Printf("!%s\n", mr.Iid)
	color.Bold.Println(mr.Title)
	if mr.Description != "" {
		fmt.Println(mr.Description)
	}
	println()
	color.OpItalic.Printf("Source branch: %s\n", mr.SourceBranch)
	color.OpItalic.Printf("Target branch: %s\n", mr.TargetBranch)
	println()

	fmt.Println("Assignees:")

	for _, user := range mr.Assignees.Nodes {
		fmt.Print("  - ")
		user.Print()
	}

	println()
	fmt.Printf("State: %s\n", color.Bold.Sprint(mr.State))
	println()
	fmt.Print("Author: ")
	mr.Author.Print()

	println()
}

func mergeRequestState() string {
	if Opened {
		return "opened"
	}

	if Closed {
		return "closed"
	}

	if Merged {
		return "merged"
	}

	return "null"
}

// MergeRequestList display a paginated merge request for a given project by path
func MergeRequestList(client *api.Client) func(*cli.Context) error {
	return func(context *cli.Context) error {
		path, err := utils.GetPathParam(context)

		if err != nil {
			return err
		}

		spinner := utils.ShowSpinner()

		var query struct {
			Project *struct {
				MergeRequests struct {
					PageInfo struct {
						EndCursor   string
						HasNextPage bool
					}
					Nodes []baseMergeRequest
				} `graphql:"(first: 10, after: $after, state: $state)"`
			} `graphql:"(fullPath:$path)"`
		}

		variables := struct {
			path  string `graphql-type:"ID!"`
			after string
			state string `graphql-type:"MergeRequestState"`
		}{
			path:  path,
			after: "",
			state: mergeRequestState(),
		}

		if err := client.Query(&query, variables); err != nil {
			return err
		}

		for {
			spinner.Stop()

			if query.Project == nil {
				return errors.New("An error has occurred, check the repository path and permissions")
			}

			for _, issue := range query.Project.MergeRequests.Nodes {
				issue.Print()
			}

			if !query.Project.MergeRequests.PageInfo.HasNextPage {
				return nil
			}

			if err := keyboard.Open(); err != nil {
				return err
			}

			defer keyboard.Close()

			for {
				char, key, _ := keyboard.GetKey()

				if key == keyboard.KeyEnter {
					break
				}

				if char == 'q' || key == keyboard.KeyCtrlC || key == keyboard.KeyEsc {
					println()
					return nil
				}
			}

			println()

			variables.after = query.Project.MergeRequests.PageInfo.EndCursor
			spinner.Start()

			client.Query(&query, variables)
		}
	}
}

// ShowMergeRequest search and display a merge request by iid
func ShowMergeRequest(client *api.Client) func(*cli.Context) error {
	return func(context *cli.Context) error {
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

		spinner := utils.ShowSpinner()

		var query struct {
			Project struct {
				MergeRequest *MergeRequest `graphql:"(iid:$iid)"`
			} `graphql:"(fullPath:$path)"`
		}

		variables := struct {
			path string `graphql-type:"ID!"`
			iid  string `graphql-type:"String!"`
		}{
			path,
			iid,
		}

		if err := client.Query(&query, variables); err != nil {
			return err
		}

		spinner.Stop()

		if query.Project.MergeRequest != nil {
			query.Project.MergeRequest.Print()
			return nil
		}

		return errors.New("An error has occurred, check the repository path and permissions")
	}
}

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
				color.Reset()
				fmt.Printf("%s (%s)\n", color.Bold.Sprint(user.Name), color.OpItalic.Sprint(user.Username))
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
			color.Reset()
			fmt.Printf("%s (%s)\n", color.Bold.Sprint(user.Name), color.OpItalic.Sprint(user.Username))
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

	color.LightGreen.Printf("%s assigned to merge request !%s\n", usernames[0], iid)
}
