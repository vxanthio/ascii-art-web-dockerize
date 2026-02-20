// Package flagparser validates command-line arguments for the ascii-art-web program.
//
// Its responsibility is limited to validating input format and structure.
// It does NOT perform rendering or color validation logic.
//
// It ensures:
//   - the correct number of arguments is provided
//   - the --color flag (if present) appears in the correct position
//   - only one --color flag is used
//   - the --color flag contains a non-empty value
//
// Any invalid input results in a usage error.
package flagparser

import (
	"errors"
	"strings"
)

// Argument count boundaries according to the project specification.
const (
	minimumArgs = 2
	maximumArgs = 5
)

// errUsage is the single user-facing error returned for any invalid CLI input.
// This keeps command-line output consistent and predictable.
//
//nolint:staticcheck // ST1005: capitalized per project specification
var errUsage = errors.New("Usage: go run . [OPTION] [STRING]\n\nEX: go run . --color=<color> <substring to be colored> \"something\"")

// ParseArgs validates the provided command-line arguments.
//
// The function checks argument count boundaries, flag syntax, flag position,
// and ensures the --color flag (if present) contains a non-empty value.
//
// Parameters:
//   - args: The command-line arguments including the program name (os.Args).
//
// Returns:
//   - An error if the arguments are invalid, nil otherwise.
func ParseArgs(args []string) error {
	colorFlagCount := 0

	if len(args) < minimumArgs || len(args) > maximumArgs {
		return errUsage
	}

	if strings.HasPrefix(args[1], "-") && !strings.HasPrefix(args[1], "--color=") {
		return errUsage
	}

	for i := 1; i < len(args); i++ {
		if strings.HasPrefix(args[i], "--color=") {
			colorFlagCount++

			if colorFlagCount > 1 {
				return errUsage
			}

			if i != 1 {
				return errUsage
			}
		}
	}

	if strings.HasPrefix(args[1], "--color=") && len(args) < 3 {
		return errUsage
	}

	if strings.HasPrefix(args[1], "--color=") {
		_, color, found := strings.Cut(args[1], "=")
		if !found || color == "" {
			return errUsage
		}
	}

	return nil
}
