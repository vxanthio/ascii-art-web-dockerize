package color_test

import (
	"testing"

	"ascii-art-web/internal/color"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name      string
		colorSpec string
		want      color.RGB
		wantErr   bool
	}{
		// Named colors
		{"named_red", "red", color.RGB{255, 0, 0}, false},
		{"named_green", "green", color.RGB{0, 255, 0}, false},
		{"named_blue", "blue", color.RGB{0, 0, 255}, false},
		{"named_case_insensitive", "RED", color.RGB{255, 0, 0}, false},
		{"named_unknown", "blurple", color.RGB{}, true},

		// Extended names
		{"named_orange", "orange", color.RGB{255, 165, 0}, false},
		{"named_purple", "purple", color.RGB{128, 0, 128}, false},
		{"named_pink", "pink", color.RGB{255, 192, 203}, false},
		{"named_brown", "brown", color.RGB{165, 42, 42}, false},
		{"named_gray", "gray", color.RGB{128, 128, 128}, false},

		// Hex
		{"hex_red", "#ff0000", color.RGB{255, 0, 0}, false},
		{"hex_invalid_length_short", "#ff0", color.RGB{}, true},
		{"hex_invalid_length_long", "#ff000000", color.RGB{}, true},
		{"hex_invalid_chars", "#gg0000", color.RGB{}, true},
		{"hex_invalid_green", "#ffgg00", color.RGB{}, true},
		{"hex_invalid_blue", "#ffffzz", color.RGB{}, true},

		// RGB
		{"rgb_red", "rgb(255, 0, 0)", color.RGB{255, 0, 0}, false},
		{"rgb_spaces", "rgb( 255 , 0 , 0 )", color.RGB{255, 0, 0}, false},
		{"rgb_invalid_count", "rgb(255)", color.RGB{}, true},
		{"rgb_out_of_range", "rgb(300, 0, 0)", color.RGB{}, true},
		{"rgb_non_number", "rgb(a, 0, 0)", color.RGB{}, true},
		{"rgb_uppercase", "RGB(255,0,0)", color.RGB{255, 0, 0}, false},
		{"rgb_no_space", "rgb(255,0,0)", color.RGB{255, 0, 0}, false},
		{"padded_named", " red ", color.RGB{255, 0, 0}, false},
		{"rgb_boundary_low", "rgb(0,0,0)", color.RGB{0, 0, 0}, false},
		{"rgb_missing_paren", "rgb(255,0,0", color.RGB{}, true},

		// Empty / whitespace
		{"empty_spec", "", color.RGB{}, true},
		{"whitespace_spec", "   ", color.RGB{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := color.Parse(tt.colorSpec)
			if (err != nil) != tt.wantErr {
				t.Fatalf(`Parse(%q) error = %v, wantErr %t`, tt.colorSpec, err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Fatalf(`Parse(%q) = %#v, want %#v`, tt.colorSpec, got, tt.want)
			}
		})
	}
}

func TestANSI(t *testing.T) {
	tests := []struct {
		name string
		rgb  color.RGB
		want string
	}{
		{"red", color.RGB{255, 0, 0}, "\033[38;2;255;0;0m"},
		{"green", color.RGB{0, 255, 0}, "\033[38;2;0;255;0m"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := color.ANSI(tt.rgb)
			if got != tt.want {
				t.Fatalf("ANSI(%#v) = %q, want %q", tt.rgb, got, tt.want)
			}

		})
	}
}
