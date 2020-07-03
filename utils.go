package waka

import (
	"fmt"
	"os"

	wzlib_utils "github.com/infra-whizz/wzlib/utils"
)

// ExitOnError with the error code 1, allowing a custom message.
func ExitOnErrorPreamble(err error, preamble string) {
	if preamble != "" {
		preamble += " "
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s%s\n", preamble, err.Error())
		os.Exit(wzlib_utils.EX_GENERIC)
	}
}

// ExitOnError with the error code 1.
func ExitOnError(err error) {
	ExitOnErrorPreamble(err, "")
}
