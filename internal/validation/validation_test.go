package validation

import (
	"strings"
	"testing"
)

func TestValidateText(t *testing.T) {
	tests := []struct {
		name    string
		text    string
		wantErr error
	}{
		{
			name:    "valid text",
			text:    "Hello World",
			wantErr: nil,
		},
		{
			name:    "valid with newline",
			text:    "Hello\nWorld",
			wantErr: nil,
		},
		{
			name:    "carriage return newline",
			text:    "Hello\r\nWorld",
			wantErr: ErrInvalidChars,
		},
		{
			name:    "valid with special chars",
			text:    "Hello! @#$%^&*()",
			wantErr: nil,
		},
		{
			name:    "empty string",
			text:    "",
			wantErr: ErrEmptyText,
		},
		{
			name:    "only whitespace",
			text:    "   \t\n  ",
			wantErr: ErrEmptyText,
		},
		{
			name:    "text too long",
			text:    strings.Repeat("a", MaxTextLength+1),
			wantErr: ErrTextTooLong,
		},
		{
			name:    "max length exactly",
			text:    strings.Repeat("a", MaxTextLength),
			wantErr: nil,
		},
		{
			name:    "non-printable char",
			text:    "Hello\x00World",
			wantErr: ErrInvalidChars,
		},
		{
			name:    "tab character",
			text:    "Hello\tWorld",
			wantErr: ErrInvalidChars,
		},
		{
			name:    "valid all printable ASCII",
			text:    " !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~",
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateText(tt.text)
			if err != tt.wantErr {
				t.Errorf("ValidateText() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateBanner(t *testing.T) {
	tests := []struct {
		name    string
		banner  string
		wantErr error
	}{
		{
			name:    "standard banner",
			banner:  "standard",
			wantErr: nil,
		},
		{
			name:    "shadow banner",
			banner:  "shadow",
			wantErr: nil,
		},
		{
			name:    "thinkertoy banner",
			banner:  "thinkertoy",
			wantErr: nil,
		},
		{
			name:    "invalid banner",
			banner:  "invalid",
			wantErr: ErrInvalidBanner,
		},
		{
			name:    "empty banner",
			banner:  "",
			wantErr: ErrInvalidBanner,
		},
		{
			name:    "uppercase banner",
			banner:  "STANDARD",
			wantErr: ErrInvalidBanner,
		},
		{
			name:    "path traversal attempt",
			banner:  "../../etc/passwd",
			wantErr: ErrInvalidBanner,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateBanner(tt.banner)
			if err != tt.wantErr {
				t.Errorf("ValidateBanner() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
