package utils

import (
	"bufio"
	"os"
)

// ReadLine read text from stdin until break line
func ReadLine() string {
	reader := bufio.NewReader(os.Stdin)
	readed, _ := reader.ReadString('\n')
	return readed[:len(readed)-1]
}
