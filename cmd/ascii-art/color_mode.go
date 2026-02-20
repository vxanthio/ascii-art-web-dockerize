package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"ascii-art-web/internal/color"
	"ascii-art-web/internal/coloring"
	"ascii-art-web/internal/flagparser"
	"ascii-art-web/internal/parser"
	"ascii-art-web/internal/renderer"
)

// runColorMode handles execution when the --color flag is detected.
//
// The function validates color mode arguments, parses the color specification,
// loads the banner, and renders ASCII art with ANSI color codes applied.
// It exits with appropriate error codes if validation or rendering fails.
//
// Parameters:
//   - args: Command-line arguments including os.Args[0].
func runColorMode(args []string) {
	if err := flagparser.ParseArgs(args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(exitCodeUsageError)
	}

	colorSpec, substring, text, bannerName, err := extractColorArgs(args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(exitCodeUsageError)
	}

	rgb, err := color.Parse(colorSpec)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(exitCodeColorError)
	}

	bannerPath, err := GetBannerPath(bannerName)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(exitCodeUsageError)
	}

	charMap, err := parser.LoadBanner(GetBannerFS(), bannerPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading banner file: %v\n", err)
		os.Exit(exitCodeBannerError)
	}

	colorCode := color.ANSI(rgb)
	lines := strings.Split(text, "\n")

	for i, line := range lines {
		if line == "" {
			if i < len(lines)-1 {
				fmt.Println()
			}
			continue
		}

		art, err := renderer.ASCII(line, charMap)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error rendering text: %v\n", err)
			os.Exit(exitCodeRenderError)
		}

		artLines := strings.Split(strings.TrimSuffix(art, "\n"), "\n")
		widths := parser.CharWidths(line, charMap)
		colored := coloring.ApplyColor(artLines, line, substring, colorCode, widths)

		for _, cl := range colored {
			fmt.Println(cl)
		}
	}
}

// hasColorFlag checks whether the first user argument uses the --color flag.
//
// Parameters:
//   - args: Command-line arguments slice including os.Args[0].
//
// Returns:
//   - true if args[1] starts with "--color", false otherwise.
func hasColorFlag(args []string) bool {
	return len(args) > 1 && strings.HasPrefix(args[1], "--color")
}

// extractColorArgs extracts color spec, substring, text, and banner from color-mode arguments.
//
// The function expects args[1] to be the --color=<value> flag. The remaining arguments
// are interpreted as follows:
//   - 3 args: prog --color=X text (no substring, default banner)
//   - 4 args: prog --color=X text banner (if last arg is valid banner name)
//   - 4 args: prog --color=X substring text (otherwise, default banner)
//   - 5 args: prog --color=X substring text banner
//
// Parameters:
//   - args: Command-line arguments including program name.
//
// Returns:
//   - colorSpec: The color value from the --color= flag.
//   - substring: The substring to color (empty if not provided).
//   - text: The text to render (with escape sequences interpreted).
//   - banner: The banner name to use.
//   - err: An error if extraction fails.
func extractColorArgs(args []string) (colorSpec, substring, text, banner string, err error) {
	_, colorSpec, _ = strings.Cut(args[1], "=")

	remaining := args[2:]

	switch len(remaining) {
	case 0:
		return "", "", "", "", errors.New("missing text argument")
	case 1:
		text = remaining[0]
		banner = defaultBanner
	case 2:
		if isValidBanner(remaining[1]) {
			text = remaining[0]
			banner = remaining[1]
		} else {
			substring = remaining[0]
			text = remaining[1]
			banner = defaultBanner
		}
	case 3:
		substring = remaining[0]
		text = remaining[1]
		banner = remaining[2]
	default:
		return "", "", "", "", errors.New("too many arguments")
	}

	text = strings.ReplaceAll(text, "\\n", "\n")

	return colorSpec, substring, text, banner, nil
}
