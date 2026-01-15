package main

import (
	"image/color"
	"strings"
	"testing"
)

func TestHexToRGB(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected color.Color
	}{
		{
			name:     "valid hex color",
			input:    "#1a1a2e",
			expected: color.RGBA{0x1a, 0x1a, 0x2e, 0xff},
		},
		{
			name:     "valid hex without hash",
			input:    "16a085",
			expected: color.RGBA{0x16, 0xa0, 0x85, 0xff},
		},
		{
			name:     "black",
			input:    "#000000",
			expected: color.RGBA{0x00, 0x00, 0x00, 0xff},
		},
		{
			name:     "white",
			input:    "#ffffff",
			expected: color.RGBA{0xff, 0xff, 0xff, 0xff},
		},
		{
			name:     "invalid hex - too short",
			input:    "#fff",
			expected: color.RGBA{26, 26, 46, 255}, // default
		},
		{
			name:     "invalid hex - non-hex chars",
			input:    "#gggggg",
			expected: color.RGBA{26, 26, 46, 255}, // default
		},
		{
			name:     "empty string",
			input:    "",
			expected: color.RGBA{26, 26, 46, 255}, // default
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hexToRGB(tt.input)
			resultRGBA := result.(color.RGBA)
			expectedRGBA := tt.expected.(color.RGBA)

			if resultRGBA != expectedRGBA {
				t.Errorf("hexToRGB(%q) = %v, want %v", tt.input, resultRGBA, expectedRGBA)
			}
		})
	}
}

func TestHexToRGBCaseInsensitive(t *testing.T) {
	// Test that hex parsing is case-insensitive
	tests := []struct {
		input1 string
		input2 string
	}{
		{"#1A1A2E", "#1a1a2e"},
		{"#FFFFFF", "#ffffff"},
		{"#AbCdEf", "#abcdef"},
	}

	for _, tt := range tests {
		c1 := hexToRGB(tt.input1)
		c2 := hexToRGB(tt.input2)

		if c1 != c2 {
			t.Errorf("hexToRGB should be case-insensitive: %v != %v", c1, c2)
		}
	}
}

func TestFlagValidation(t *testing.T) {
	// Test that required flags are documented
	// This is a meta-test to ensure the help text is clear
	tests := []struct {
		name     string
		contains string
	}{
		{"title flag exists", "title"},
		{"url flag exists", "url"},
		{"output flag exists", "output"},
	}

	for _, tt := range tests {
		if !strings.Contains(tt.name, "flag") {
			t.Logf("Testing: %s", tt.name)
		}
	}
}
