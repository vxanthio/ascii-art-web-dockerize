package renderer_test

import (
	"strings"
	"testing"

	"ascii-art-web/internal/renderer"
)

func TestEmptyInput(t *testing.T) {
	input := ""
	banner := map[rune][]string{}
	output, err := renderer.ASCII(input, banner)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if input != output {
		t.Errorf("expected:\n%q\ngot:\n%q", input, output)
	}
}

func TestSingleCharacter(t *testing.T) {
	input := "A"
	expected := `A1
A2
A3
A4
A5
A6
A7
A8
`
	banner := map[rune][]string{
		'A': {"A1", "A2", "A3", "A4", "A5", "A6", "A7", "A8"},
	}
	output, err := renderer.ASCII(input, banner)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if expected != output {
		t.Errorf("expected:\n%q\ngot:\n%q", expected, output)
	}
}

func TestMultipleCharacters(t *testing.T) {
	input := "AB"
	expected := `A1B1
A2B2
A3B3
A4B4
A5B5
A6B6
A7B7
A8B8
`
	banner := map[rune][]string{
		'A': {"A1", "A2", "A3", "A4", "A5", "A6", "A7", "A8"},
		'B': {"B1", "B2", "B3", "B4", "B5", "B6", "B7", "B8"},
	}
	output, err := renderer.ASCII(input, banner)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if expected != output {
		t.Errorf("expected:\n%q\ngot:\n%q", expected, output)
	}
}

func TestSpaceBetweenCharacters(t *testing.T) {
	input := "A A"
	expected := `A1  A1
A2  A2
A3  A3
A4  A4
A5  A5
A6  A6
A7  A7
A8  A8
`
	banner := map[rune][]string{
		'A': {"A1", "A2", "A3", "A4", "A5", "A6", "A7", "A8"},
		' ': {"  ", "  ", "  ", "  ", "  ", "  ", "  ", "  "},
	}
	output, err := renderer.ASCII(input, banner)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if expected != output {
		t.Errorf("expected:\n%q\ngot:\n%q", expected, output)
	}
}

func TestNumbersBetweenCharacters(t *testing.T) {
	input := "A1A"
	expected := `A11A1
A21A2
A31A3
A41A4
A51A5
A61A6
A71A7
A81A8
`
	banner := map[rune][]string{
		'A': {"A1", "A2", "A3", "A4", "A5", "A6", "A7", "A8"},
		'1': {"1", "1", "1", "1", "1", "1", "1", "1"},
	}
	output, err := renderer.ASCII(input, banner)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if expected != output {
		t.Errorf("expected:\n%q\ngot:\n%q", expected, output)
	}
}

func TestAllSpecialCharacters(t *testing.T) {
	specials := `!"#$%&'()*+,-./:;<=>?@[\]^_{|}~`

	banner := make(map[rune][]string)
	for _, ch := range specials {
		banner[ch] = []string{
			string(ch), string(ch), string(ch), string(ch),
			string(ch), string(ch), string(ch), string(ch),
		}
	}

	output, err := renderer.ASCII(specials, banner)
	if err != nil {
		t.Fatalf("ASCII failed for special characters: %v", err)
	}

	lines := strings.Split(output, "\n")
	if len(lines) != 9 {
		t.Fatalf("expected 9 lines (8 content + trailing newline), got %d", len(lines))
	}

	for i, line := range lines[:8] {
		if len(line) != len(specials) {
			t.Errorf("line %d: expected %d chars, got %d", i, len(specials), len(line))
		}
	}
}

func TestNewlineBetweenCharacters(t *testing.T) {
	input := "A\nB"
	expected := `A1
A2
A3
A4
A5
A6
A7
A8
B1
B2
B3
B4
B5
B6
B7
B8
`
	banner := map[rune][]string{
		'A': {"A1", "A2", "A3", "A4", "A5", "A6", "A7", "A8"},
		'B': {"B1", "B2", "B3", "B4", "B5", "B6", "B7", "B8"},
	}
	output, err := renderer.ASCII(input, banner)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if expected != output {
		t.Errorf("expected:\n%q\ngot:\n%q", expected, output)
	}
}

func TestTrailingNewline(t *testing.T) {
	input := "A\n"
	expected := `A1
A2
A3
A4
A5
A6
A7
A8
`
	banner := map[rune][]string{
		'A': {"A1", "A2", "A3", "A4", "A5", "A6", "A7", "A8"},
	}
	output, err := renderer.ASCII(input, banner)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if expected != output {
		t.Errorf("expected:\n%q\ngot:\n%q", expected, output)
	}
}

func TestConsecutiveNewlines(t *testing.T) {
	input := "A\n\nB"
	expected := `A1
A2
A3
A4
A5
A6
A7
A8

B1
B2
B3
B4
B5
B6
B7
B8
`
	banner := map[rune][]string{
		'A': {"A1", "A2", "A3", "A4", "A5", "A6", "A7", "A8"},
		'B': {"B1", "B2", "B3", "B4", "B5", "B6", "B7", "B8"},
	}
	output, err := renderer.ASCII(input, banner)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if expected != output {
		t.Errorf("expected:\n%q\ngot:\n%q", expected, output)
	}
}

func TestMissingCharacter(t *testing.T) {
	input := "AB"
	banner := map[rune][]string{
		'A': {"A1", "A2", "A3", "A4", "A5", "A6", "A7", "A8"},
	}
	output, err := renderer.ASCII(input, banner)
	if err == nil {
		t.Error("expected error for missing character 'B', got nil")
	}
	expectedMsg := "character B"
	if !strings.Contains(err.Error(), expectedMsg) {
		t.Errorf("error message should mention character B, got: %v", err)
	}
	if output != "" {
		t.Errorf("expected empty output on error, got %q", output)
	}
}

func TestCorruptedBanner(t *testing.T) {
	input := "A"
	banner := map[rune][]string{
		'A': {"A1", "A2", "A3", "A4", "A6", "A7", "A8"},
	}
	output, err := renderer.ASCII(input, banner)
	if err == nil {
		t.Error("expected error for corrupted banner, got nil")
	}
	if output != "" {
		t.Errorf("expected empty output on error, got %q", output)
	}
}

func TestInvalidCharacters(t *testing.T) {
	input := "A	B"
	banner := map[rune][]string{
		'A': {"A1", "A2", "A3", "A4", "A5", "A6", "A7", "A8"},
		'B': {"B1", "B2", "B3", "B4", "B5", "B6", "B7", "B8"},
	}
	output, err := renderer.ASCII(input, banner)
	if err == nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if output != "" {
		t.Errorf("expected empty output on error, got %q", output)
	}
}

func TestCompleteASCIIRange(t *testing.T) {
	banner := make(map[rune][]string)

	for ch := rune(32); ch <= 126; ch++ {
		banner[ch] = []string{
			"*", "*", "*", "*", "*", "*", "*", "*",
		}
	}

	var input strings.Builder
	for ch := rune(32); ch <= 126; ch++ {
		input.WriteRune(ch)
	}

	output, err := renderer.ASCII(input.String(), banner)
	if err != nil {
		t.Fatalf("ASCII failed for ASCII range: %v", err)
	}

	lines := strings.Split(output, "\n")
	if len(lines) != 9 {
		t.Fatalf("expected 9 lines (8 content + trailing newline), got %d", len(lines))
	}

	for i, line := range lines[:8] {
		if len(line) != 95 {
			t.Errorf("line %d: expected 95 chars, got %d", i, len(line))
		}
	}
}
