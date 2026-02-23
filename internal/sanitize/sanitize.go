// Package sanitize provides HTML sanitization to prevent XSS attacks.
// It escapes dangerous HTML characters in user input.
package sanitize

import "html"

func HTML(input string) string {
	return html.EscapeString(input)
}
