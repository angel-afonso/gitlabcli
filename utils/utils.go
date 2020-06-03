package utils

import (
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strings"

	"gopkg.in/gookit/color.v1"
)

// ReadLine read text from stdin until break line
func ReadLine() string {
	var input string
	fmt.Scanf("%s", &input)
	return input
}

// ReadInt get a int value from user input
func ReadInt() int {
	var input int
	fmt.Scanf("%d", &input)
	return input
}

// IsGitRepository returns true if the current working directory is a valid git repository
func IsGitRepository() bool {
	err := exec.Command("git", "-C", ".", "rev-parse").Run()
	return err == nil
}

// GetRemote returns an array with the remote repositories names
func GetRemote() []string {
	cmd := exec.Command("git", "remote")
	output, _ := cmd.Output()

	remotes := strings.Split(string(output), " \n")
	return remotes
}

// AskRemote ask to the user which remote repository yse
func AskRemote(remotes []string) string {
	color.Cyan.Println("Select Remote")
	for i := 0; i < len(remotes); i++ {
		fmt.Printf("(%d): %s\n", i, remotes[i])
	}

	index := ReadInt()

	if index > len(remotes)-1 {
		log.Fatal("Invalid index")
	}

	return remotes[index]
}

// GetRemotePath returns a string with the repo path for the given remote repository
func GetRemotePath(remote string) string {
	cmd := exec.Command("git", "remote", "get-url", strings.TrimSpace(remote))
	output, _ := cmd.Output()

	return strings.TrimSpace(regexp.MustCompile(`https://.+[^.].com/|git@.+[^:]:|\.git`).ReplaceAllString(string(output), ""))
}
