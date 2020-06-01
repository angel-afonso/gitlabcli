package utils

import (
	"bufio"
	"fmt"
	"os"
)

// ReadLine read text from stdin until break line
func ReadLine(text string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(text)
	readed, _ := reader.ReadString('\n')
	return readed[:len(readed)-1]
}
