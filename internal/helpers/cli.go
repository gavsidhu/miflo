package helpers

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func PromptForConfirmation(prompt string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s [y/n]: ", prompt)

	response, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		return false
	}

	response = strings.TrimSpace(response)

	return strings.ToLower(response) == "y" || strings.ToLower(response) == "yes"
}
