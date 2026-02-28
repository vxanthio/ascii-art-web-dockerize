// Package validation provides input validation for user-submitted data.
// It validates text content and banner names to ensure safe processing.
package validation

import (
	"errors"
	"strings"
)

var (
	ErrEmptyText     = errors.New("text cannot be empty")
	ErrTextTooLong   = errors.New("text exceeds maximum length")
	ErrInvalidChars  = errors.New("text contains non-printable characters")
	ErrInvalidBanner = errors.New("invalid banner name")
)

const MaxTextLength = 1000

func ValidateText(text string) error {
	if strings.TrimSpace(text) == "" {
		return ErrEmptyText
	}

	if len(text) > MaxTextLength {
		return ErrTextTooLong
	}

	for _, ch := range text {
		if ch != '\n' && (ch < 32 || ch > 126) {
			return ErrInvalidChars
		}
	}

	return nil
}

func ValidateBanner(banner string) error {
	validBanners := map[string]bool{
		"standard":   true,
		"shadow":     true,
		"thinkertoy": true,
	}

	if !validBanners[banner] {
		return ErrInvalidBanner
	}

	return nil
}
