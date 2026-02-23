package sanitize

import "html"

func HTML(input string) string {
	return html.EscapeString(input)
}
