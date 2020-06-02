package utils

import (
	"bufio"
	"os"
	"os/exec"
	"strings"
)

// ReadLine read text from stdin until break line
func ReadLine() string {
	reader := bufio.NewReader(os.Stdin)
	readed, _ := reader.ReadString('\n')
	return readed[:len(readed)-1]
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

	remotes := strings.Split(string(output), "\n")
	return remotes
}

// GetRemotePath returns a string with the repo path for the given remote repository
func GetRemotePath(remote string) string {
	cmd := exec.Command("git", "remote", "get-url", remote)
	output, _ := cmd.Output()
	return strings.NewReplacer("https://gitlab.com/", "", "git@gitlab.com:", "", ".git", "").Replace(string(output))
}
