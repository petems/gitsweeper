package internal

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// AskForConfirmation asks the user for confirmation. A user must type in "yes" or "no" and
// then press enter. It has fuzzy matching, so "y", "Y", "yes", "YES", and "Yes" all count as
// confirmations. If the input is not recognized, it will ask again. The function does not return
// until it gets a valid response from the user.
//
//	AskForConfirmation("Do you want to restart?", os.Stdin)
func AskForConfirmation(s string, in io.Reader) (bool, error) {
	reader := bufio.NewReader(in)

	for {
		fmt.Printf("%s [y/n]: ", s)

		response, err := reader.ReadString('\n')
		if err != nil {
			return false, err
		}

		response = strings.ToLower(strings.TrimSpace(response))

		switch response {
		case "y", "yes":
			return true, nil
		case "n", "no":
			return false, nil
		}
	}
}
