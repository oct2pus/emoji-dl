package emoji

import "fmt"

// hasError is a helper function to print errors.
func hasError(err error) bool {
	if err != nil {
		if Verbose.Errors {
			fmt.Printf("error: %v\n", err)
		}
		return true
	}
	return false
}
