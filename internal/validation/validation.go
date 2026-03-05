// Package validation provides input validation for user-submitted data.
// It validates text content and banner names to ensure safe processing.
package validation

import (
	"errors"
	"strings"
)

// Sentinel errors returned by validation functions.
// Callers can compare against these values using errors.Is.
var (
	ErrEmptyText     = errors.New("text cannot be empty")
	ErrTextTooLong   = errors.New("text exceeds maximum length")
	ErrInvalidChars  = errors.New("text contains non-printable characters")
	ErrInvalidBanner = errors.New("invalid banner name")
)

// MaxTextLength is the maximum number of characters allowed in a text submission.
// It must match the maxlength attribute on the HTML textarea.
const MaxTextLength = 1000

// ValidateText checks that the submitted text is non-empty, within the allowed
// length, and contains only printable ASCII characters (codes 32–126) or newlines.
//
// Parameters:
//   - text: The raw string submitted by the user.
//
// Returns:
//   - ErrEmptyText if text is empty or whitespace-only.
//   - ErrTextTooLong if text exceeds MaxTextLength characters.
//   - ErrInvalidChars if text contains non-printable or non-ASCII characters.
//   - nil if text is valid.
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

// ValidateBanner checks that the submitted banner name is one of the three
// supported options: standard, shadow, or thinkertoy.
//
// The whitelist approach implicitly prevents path traversal attacks — any
// string not in the whitelist is rejected regardless of its content.
//
// Parameters:
//   - banner: The banner name submitted by the user.
//
// Returns:
//   - ErrInvalidBanner if the name is not a recognized banner.
//   - nil if the banner name is valid.
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
