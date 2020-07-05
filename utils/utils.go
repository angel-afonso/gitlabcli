package utils

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/urfave/cli/v2"
	"gopkg.in/gookit/color.v1"
)

// ShowSpinner display and return spinner instance
func ShowSpinner() *spinner.Spinner {
	spinner := spinner.New(spinner.CharSets[32], 100*time.Millisecond)
	spinner.Start()
	return spinner
}

// Ternary simulate ternary operator
func Ternary(condition bool, value1 interface{}, value2 interface{}) interface{} {
	if condition {
		return value1
	}
	return value2
}

// ReadLine read text from stdin until break line
func ReadLine() string {
	reader := bufio.NewReader(os.Stdin)
	readed, _ := reader.ReadString('\n')

	return readed[:len(readed)-1]
}

// ReadLineOptional read text from stdin until break line, if the input is empty, return a default value
func ReadLineOptional(defaultValue string) string {
	input := ReadLine()

	if input == "" {
		return defaultValue
	}

	return input
}

// ReadInt get a int value from user input
func ReadInt() int {
	var input string
	for {
		fmt.Scanln(&input)
		regexp.MustCompile(`\d+`).FindString(input)
		number, err := strconv.Atoi(regexp.MustCompile(`\d+`).FindString(input))

		if err == nil {
			return number
		}

		color.Red.Println("Invalid input.")
		color.Reset()
	}

}

// IsGitRepository returns true if the current working directory is a valid git repository
func IsGitRepository() bool {
	return Repository() != nil
}

// GetRemotes returns an array with references to remote repositories
func GetRemotes() []*git.Remote {
	remotes, _ := Repository().Remotes()
	return remotes
}

// AskRemote ask to the user which remote repository yse
func AskRemote(remotes []*git.Remote) string {
	color.Cyan.Println("Select Remote")
	for i := 0; i < len(remotes); i++ {
		color.Cyan.Printf("(%d): ", i)
		color.Green.Printf("%s\n", remotes[i].Config().Name)
	}

	index := ReadInt()

	if index > len(remotes)-1 {
		log.Fatal("Invalid index")
	}

	return remotes[index].Config().Name
}

// GetRemotePath returns a string with the repo path for the given remote repository
func GetRemotePath(remote string) string {
	cmd := exec.Command("git", "remote", "get-url", strings.TrimSpace(remote))
	output, _ := cmd.Output()

	return strings.TrimSpace(regexp.MustCompile(`https://.+[^.].com/|git@.+[^:]:|\.git`).ReplaceAllString(string(output), ""))
}

// Repository return current repository reference
func Repository() *git.Repository {
	dir, _ := os.Getwd()
	repository, _ := git.PlainOpen(dir)
	return repository
}

// RepoHead return repository head reference
func RepoHead() *plumbing.Reference {
	if repository := Repository(); repository != nil {
		head, _ := repository.Head()
		return head
	}
	return nil
}

// RepoCommits return repository commits iterator
func RepoCommits() object.CommitIter {
	if repository := Repository(); repository != nil {
		commits, _ := repository.CommitObjects()
		return commits
	}
	return nil
}

// RepoLastCommit return last commit in git repository
func RepoLastCommit() *object.Commit {
	if repository := RepoHead(); repository != nil {
		commit, _ := Repository().CommitObject(repository.Hash())
		return commit
	}
	return nil
}

// GetPathParam find path repository in command line arg or in the directory
func GetPathParam(context *cli.Context) string {
	var path string

	if IsGitRepository() && context.Args().Len() == 0 {
		if remotes := GetRemotes(); len(remotes) > 1 {
			path = GetRemotePath(AskRemote(remotes))
		} else if len(remotes) == 1 {
			path = GetRemotePath(remotes[0].Config().Name)
		} else {
			color.Red.Println("No repo path provided")
			color.Reset()
			os.Exit(1)
		}

	} else if context.Args().Len() > 0 {
		path = context.Args().First()
	} else {
		log.Fatal("Expected project path")
	}

	return path
}
