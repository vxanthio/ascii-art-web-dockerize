// Package main provides the ASCII art generator CLI application.
//
// The application orchestrates the parser, renderer, color, flagparser, and coloring
// packages to convert text input into graphical ASCII art representations, optionally
// with ANSI color support for full text or specific substrings.
//
// Usage:
//
//	go run . "text" [banner]
//	go run . --color=<color> "text" [banner]
//	go run . --color=<color> <substring> "text" [banner]
//
// Responsibilities of this package:
//   - Parse and validate command-line arguments
//   - Route between normal mode and color mode
//   - Validate and resolve banner file paths
//   - Coordinate between parser, renderer, and coloring
//   - Handle errors with appropriate exit codes
//
// Any invalid input, missing files, or rendering errors are reported to stderr.
package main

import (
	"fmt"
	"os"

	"ascii-art-web/internal/parser"
	"ascii-art-web/internal/renderer"
)

const (
	// Exit codes for different error scenarios.
	exitCodeUsageError  = 1
	exitCodeBannerError = 2
	exitCodeRenderError = 3
	exitCodeColorError  = 4

	// Default banner style.
	defaultBanner = "standard"
)

// main is the entry point of the ascii-art application.
//
// It determines whether to run in normal mode or color mode based on
// the presence of the --color flag, then orchestrates the appropriate
// packages to render ASCII art with optional ANSI color codes.
func main() {
	if hasColorFlag(os.Args) {
		runColorMode(os.Args)
		return
	}

	text, banner, err := ParseArgs(os.Args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(exitCodeUsageError)
	}

	bannerPath, err := GetBannerPath(banner)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(exitCodeUsageError)
	}

	charMap, err := parser.LoadBanner(GetBannerFS(), bannerPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading banner file: %v\n", err)
		os.Exit(exitCodeBannerError)
	}

	result, err := renderer.ASCII(text, charMap)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error rendering text: %v\n", err)
		os.Exit(exitCodeRenderError)
	}

	fmt.Print(result)
}
