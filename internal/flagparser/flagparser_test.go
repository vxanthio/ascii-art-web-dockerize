package flagparser_test

import (
	"testing"

	"ascii-art-web/internal/flagparser"
)

func TestParseArgs(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "no arguments",
			args:    []string{"program"},
			wantErr: true,
		},
		{
			name:    "too many arguments",
			args:    []string{"program", "banner", "--color=red", "substring", "some text", "EXTRA"},
			wantErr: true,
		},
		{
			name:    "invalid color flag prefix",
			args:    []string{"program", "-color:black", "some text"},
			wantErr: true,
		},
		{
			name:    "invalid color flag format without equals",
			args:    []string{"program", "--color:red", "some text"},
			wantErr: true,
		},
		{
			name:    "valid single string",
			args:    []string{"program", "text"},
			wantErr: false,
		},
		{
			name:    "color flag with string and no substring",
			args:    []string{"program", "--color=red", "text"},
			wantErr: false,
		},
		{
			name:    "color flag with string and substring",
			args:    []string{"program", "--color=red", "text", "substring"},
			wantErr: false,
		},
		{
			name:    "missing string after color flag",
			args:    []string{"program", "--color=red"},
			wantErr: true,
		},
		{
			name:    "multiple color flags",
			args:    []string{"program", "--color=red", "--color=blue", "text"},
			wantErr: true,
		},
		{
			name:    "invalid color flag position",
			args:    []string{"program", "text", "--color=red"},
			wantErr: true,
		},
		{
			name:    "empty color value",
			args:    []string{"program", "--color=", "text"},
			wantErr: true,
		},
		{
			name:    "valid banner with color",
			args:    []string{"program", "--color=red", "text", "standard"},
			wantErr: false,
		},
		{
			name:    "valid banner without color",
			args:    []string{"program", "text", "standard"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := flagparser.ParseArgs(tt.args)

			if tt.wantErr && err == nil {
				t.Errorf("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
