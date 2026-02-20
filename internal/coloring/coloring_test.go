package coloring_test

import (
	"strings"
	"testing"

	"ascii-art-web/internal/coloring"
)

func TestApplyColor_AdvancedCases(t *testing.T) {
	tests := []struct {
		name       string
		text       string
		substring  string
		charWidths []int
		asciiArt   []string
		wantCount  int
	}{
		{
			name:       "Match at beginning",
			text:       "hello",
			substring:  "he",
			charWidths: []int{1, 1, 1, 1, 1},
			asciiArt:   []string{"hello"},
			wantCount:  1,
		},
		{
			name:       "Variable character widths",
			text:       "ABC",
			substring:  "B",
			charWidths: []int{3, 6, 3},
			asciiArt:   []string{"AAABBBBBBCCC"},
			wantCount:  1,
		},
		{
			name:       "Multi-line ASCII art",
			text:       "hi",
			substring:  "h",
			charWidths: []int{2, 2},
			asciiArt:   []string{"HHII", "HHII"},
			wantCount:  2,
		},
		{
			name:       "Overlapping matches",
			text:       "banana",
			substring:  "ana",
			charWidths: []int{1, 1, 1, 1, 1, 1},
			asciiArt:   []string{"banana"},
			wantCount:  1,
		},
		{
			name:       "Art wider than widths",
			text:       "A",
			substring:  "A",
			charWidths: []int{1},
			asciiArt:   []string{"A (remainder)"},
			wantCount:  1,
		},
	}

	colorCode := "\033[31m"
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := coloring.ApplyColor(tt.asciiArt, tt.text, tt.substring, colorCode, tt.charWidths)
			count := 0
			for _, line := range got {
				count += strings.Count(line, colorCode)
				if tt.name == "Art wider than widths" && !strings.Contains(line, "(remainder)") {
					t.Errorf("Lost remainder content")
				}
			}
			if count != tt.wantCount {
				t.Errorf("got %d segments, want %d", count, tt.wantCount)
			}
		})
	}
}

func TestApplyColor_FinalCoverage(t *testing.T) {
	colorCode := "\033[31m"

	t.Run("Empty_Inputs", func(t *testing.T) {
		if coloring.ApplyColor(nil, "a", "a", colorCode, []int{1}) != nil {
			t.Error("Nil art fail")
		}
		art := []string{"a"}
		if len(coloring.ApplyColor(art, "a", "a", colorCode, []int{})) != 1 {
			t.Error("Empty widths fail")
		}
		if len(coloring.ApplyColor(art, "", "a", colorCode, []int{1})) != 1 {
			t.Error("Empty text fail")
		}
		if len(coloring.ApplyColor(art, "a", "", colorCode, []int{1})) != 1 {
			t.Error("Empty sub fail")
		}
	})

	t.Run("Non_Contiguous", func(t *testing.T) {
		res := coloring.ApplyColor([]string{"abcabc"}, "abcabc", "a", colorCode, []int{1, 1, 1, 1, 1, 1})
		if strings.Count(res[0], colorCode) != 2 {
			t.Error("Non-contiguous match fail")
		}
	})

	t.Run("Mismatched_Lengths", func(t *testing.T) {
		res := coloring.ApplyColor([]string{"A"}, "AB", "A", colorCode, []int{1})
		if !strings.Contains(res[0], "A") {
			t.Error("Mismatched length fail")
		}
	})

	t.Run("Varying_Line_Lengths", func(t *testing.T) {
		art := []string{"ABC", "A"}
		res := coloring.ApplyColor(art, "ABC", "A", colorCode, []int{1, 1, 1})
		if len(res) != 2 {
			t.Error("Varying length fail")
		}
	})

	t.Run("Edge_Transitions", func(t *testing.T) {
		res := coloring.ApplyColor([]string{"abc"}, "abc", "c", colorCode, []int{1, 1, 1})
		if !strings.HasSuffix(res[0], coloring.Reset) {
			t.Error("End transition fail")
		}
	})

	t.Run("Longer_Substring", func(t *testing.T) {
		res := coloring.ApplyColor([]string{"a"}, "a", "abc", colorCode, []int{1})
		if strings.Contains(res[0], colorCode) {
			t.Error("Long substring fail")
		}
	})

	t.Run("Short_Art_Truncation", func(t *testing.T) {
		res := coloring.ApplyColor([]string{"12"}, "A", "A", colorCode, []int{10})
		if !strings.Contains(res[0], "12") {
			t.Error("Truncation logic fail")
		}
	})
}
