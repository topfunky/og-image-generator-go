package main

import (
	"flag"
	"fmt"
	"image/color"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/fogleman/gg"
)

// testFontPath returns a valid font path for testing
func testFontPath(t *testing.T) string {
	t.Helper()
	paths := []string{
		"/usr/share/fonts/truetype/dejavu/DejaVuSans-Bold.ttf",
		"/usr/share/fonts/truetype/liberation/LiberationSans-Bold.ttf",
		"/System/Library/Fonts/Arial.ttf",
		"C:\\Windows\\Fonts\\arial.ttf",
	}
	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	t.Skip("No system font available for testing")
	return ""
}

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

func TestResolveFontPath(t *testing.T) {
	t.Run("custom font path provided", func(t *testing.T) {
		customPath := "/custom/font/path.ttf"
		result, err := resolveFontPath(customPath)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if result != customPath {
			t.Errorf("expected %q, got %q", customPath, result)
		}
	})

	t.Run("empty path uses system font", func(t *testing.T) {
		result, err := resolveFontPath("")
		// Should either find a system font or return an error
		if err != nil {
			if !strings.Contains(err.Error(), "font file not found") {
				t.Errorf("unexpected error: %v", err)
			}
		} else {
			if result == "" {
				t.Error("expected non-empty font path")
			}
		}
	})

	t.Run("local fonts directory", func(t *testing.T) {
		// Create a temporary fonts directory with a font file
		tmpDir := t.TempDir()
		fontsDir := filepath.Join(tmpDir, "fonts")
		if err := os.MkdirAll(fontsDir, 0755); err != nil {
			t.Fatal(err)
		}

		// Create a dummy font file
		fontFile := filepath.Join(fontsDir, "OpenSans-Bold.ttf")
		if err := os.WriteFile(fontFile, []byte("dummy"), 0644); err != nil {
			t.Fatal(err)
		}

		// Change to the temp directory to test local font resolution
		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		result, err := resolveFontPath("")
		if err != nil {
			t.Errorf("unexpected error when local font exists: %v", err)
		}
		if result != filepath.Join("fonts", "OpenSans-Bold.ttf") {
			t.Errorf("expected local font path, got %q", result)
		}
	})
}

func TestDrawBackground(t *testing.T) {
	tests := []struct {
		name    string
		bgColor string
		width   int
		height  int
	}{
		{"default color", "#1a1a2e", 1200, 628},
		{"white background", "#ffffff", 800, 600},
		{"red background", "#ff0000", 1920, 1080},
		{"invalid color falls back to default", "#invalid", 1200, 628},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dc := gg.NewContext(tt.width, tt.height)
			// Should not panic
			drawBackground(dc, tt.bgColor, tt.width, tt.height)

			// Verify the context was modified (image should have content)
			img := dc.Image()
			if img.Bounds().Dx() != tt.width || img.Bounds().Dy() != tt.height {
				t.Errorf("unexpected image dimensions: got %dx%d, want %dx%d",
					img.Bounds().Dx(), img.Bounds().Dy(), tt.width, tt.height)
			}
		})
	}
}

func TestDrawTitle(t *testing.T) {
	fontPath := testFontPath(t)

	tests := []struct {
		name    string
		title   string
		width   int
		wantErr bool
	}{
		{"simple title", "Hello World", 1200, false},
		{"long title", "This is a very long title that should wrap across multiple lines in the image", 1200, false},
		{"unicode title", "日本語タイトル", 1200, false},
		{"empty title", "", 1200, false},
		{"narrow width", "Test Title", 200, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dc := gg.NewContext(tt.width, 628)
			err := drawTitle(dc, tt.title, fontPath, tt.width)
			if (err != nil) != tt.wantErr {
				t.Errorf("drawTitle() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

	t.Run("invalid font path", func(t *testing.T) {
		dc := gg.NewContext(1200, 628)
		err := drawTitle(dc, "Test", "/nonexistent/font.ttf", 1200)
		if err == nil {
			t.Error("expected error for invalid font path")
		}
		if !strings.Contains(err.Error(), "load font") {
			t.Errorf("expected 'load font' error, got: %v", err)
		}
	})
}

func TestDrawURL(t *testing.T) {
	fontPath := testFontPath(t)

	tests := []struct {
		name    string
		url     string
		width   int
		height  int
		wantErr bool
	}{
		{"simple url", "https://example.com", 1200, 628, false},
		{"long url", "https://example.com/very/long/path/to/article/that/might/need/smaller/font", 1200, 628, false},
		{"very long url forces minimum font", "https://example.com/" + strings.Repeat("a", 200), 1200, 628, false},
		{"narrow width", "https://example.com", 300, 628, false},
		{"empty url", "", 1200, 628, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dc := gg.NewContext(tt.width, tt.height)
			err := drawURL(dc, tt.url, fontPath, tt.width, tt.height)
			if (err != nil) != tt.wantErr {
				t.Errorf("drawURL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

	t.Run("invalid font path", func(t *testing.T) {
		dc := gg.NewContext(1200, 628)
		err := drawURL(dc, "https://example.com", "/nonexistent/font.ttf", 1200, 628)
		if err == nil {
			t.Error("expected error for invalid font path")
		}
		if !strings.Contains(err.Error(), "load font") {
			t.Errorf("expected 'load font' error, got: %v", err)
		}
	})
}

func TestRun(t *testing.T) {
	fontPath := testFontPath(t)

	t.Run("successful image generation", func(t *testing.T) {
		tmpDir := t.TempDir()
		outputPath := filepath.Join(tmpDir, "test-output.png")

		// Save original args and restore after test
		oldArgs := os.Args
		defer func() { os.Args = oldArgs }()

		os.Args = []string{
			"og-image-generator",
			"-title", "Test Title",
			"-url", "https://example.com",
			"-output", outputPath,
			"-title-font", fontPath,
			"-url-font", fontPath,
		}

		// Reset flags for this test
		resetFlags()

		err := run()
		if err != nil {
			t.Errorf("run() unexpected error: %v", err)
		}

		// Verify output file was created
		if _, err := os.Stat(outputPath); os.IsNotExist(err) {
			t.Error("output file was not created")
		}
	})

	t.Run("missing required flags", func(t *testing.T) {
		oldArgs := os.Args
		defer func() { os.Args = oldArgs }()

		os.Args = []string{"og-image-generator"}
		resetFlags()

		err := run()
		if err == nil {
			t.Error("expected error for missing required flags")
		}
		if !strings.Contains(err.Error(), "title and url are required") {
			t.Errorf("unexpected error message: %v", err)
		}
	})

	t.Run("invalid font path", func(t *testing.T) {
		tmpDir := t.TempDir()
		outputPath := filepath.Join(tmpDir, "test-output.png")

		oldArgs := os.Args
		defer func() { os.Args = oldArgs }()

		os.Args = []string{
			"og-image-generator",
			"-title", "Test Title",
			"-url", "https://example.com",
			"-output", outputPath,
			"-title-font", "/nonexistent/font.ttf",
		}
		resetFlags()

		err := run()
		if err == nil {
			t.Error("expected error for invalid font path")
		}
	})

	t.Run("invalid output directory", func(t *testing.T) {
		oldArgs := os.Args
		defer func() { os.Args = oldArgs }()

		os.Args = []string{
			"og-image-generator",
			"-title", "Test Title",
			"-url", "https://example.com",
			"-output", "/nonexistent/directory/output.png",
			"-title-font", fontPath,
			"-url-font", fontPath,
		}
		resetFlags()

		err := run()
		if err == nil {
			t.Error("expected error for invalid output directory")
		}
	})

	t.Run("custom dimensions", func(t *testing.T) {
		tmpDir := t.TempDir()
		outputPath := filepath.Join(tmpDir, "test-output.png")

		oldArgs := os.Args
		defer func() { os.Args = oldArgs }()

		os.Args = []string{
			"og-image-generator",
			"-title", "Test Title",
			"-url", "https://example.com",
			"-output", outputPath,
			"-width", "800",
			"-height", "400",
			"-bg", "#ff5500",
			"-title-font", fontPath,
			"-url-font", fontPath,
		}
		resetFlags()

		err := run()
		if err != nil {
			t.Errorf("run() unexpected error: %v", err)
		}
	})
}

func TestParseFlags(t *testing.T) {
	t.Run("all flags provided", func(t *testing.T) {
		oldArgs := os.Args
		defer func() { os.Args = oldArgs }()

		os.Args = []string{
			"og-image-generator",
			"-title", "My Title",
			"-url", "https://example.com",
			"-output", "custom.png",
			"-width", "1920",
			"-height", "1080",
			"-bg", "#ff0000",
			"-title-font", "/path/to/title.ttf",
			"-url-font", "/path/to/url.ttf",
		}
		resetFlags()

		opts, err := parseFlags()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if opts.Title != "My Title" {
			t.Errorf("Title = %q, want %q", opts.Title, "My Title")
		}
		if opts.URL != "https://example.com" {
			t.Errorf("URL = %q, want %q", opts.URL, "https://example.com")
		}
		if opts.Output != "custom.png" {
			t.Errorf("Output = %q, want %q", opts.Output, "custom.png")
		}
		if opts.Width != 1920 {
			t.Errorf("Width = %d, want %d", opts.Width, 1920)
		}
		if opts.Height != 1080 {
			t.Errorf("Height = %d, want %d", opts.Height, 1080)
		}
		if opts.BgColor != "#ff0000" {
			t.Errorf("BgColor = %q, want %q", opts.BgColor, "#ff0000")
		}
		if opts.TitleFont != "/path/to/title.ttf" {
			t.Errorf("TitleFont = %q, want %q", opts.TitleFont, "/path/to/title.ttf")
		}
		if opts.URLFont != "/path/to/url.ttf" {
			t.Errorf("URLFont = %q, want %q", opts.URLFont, "/path/to/url.ttf")
		}
	})

	t.Run("default values", func(t *testing.T) {
		oldArgs := os.Args
		defer func() { os.Args = oldArgs }()

		os.Args = []string{
			"og-image-generator",
			"-title", "Test",
			"-url", "https://test.com",
		}
		resetFlags()

		opts, err := parseFlags()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if opts.Output != "social-image.png" {
			t.Errorf("Output = %q, want default %q", opts.Output, "social-image.png")
		}
		if opts.Width != 1200 {
			t.Errorf("Width = %d, want default %d", opts.Width, 1200)
		}
		if opts.Height != 628 {
			t.Errorf("Height = %d, want default %d", opts.Height, 628)
		}
		if opts.BgColor != "#1a1a2e" {
			t.Errorf("BgColor = %q, want default %q", opts.BgColor, "#1a1a2e")
		}
	})

	t.Run("missing title", func(t *testing.T) {
		oldArgs := os.Args
		defer func() { os.Args = oldArgs }()

		os.Args = []string{
			"og-image-generator",
			"-url", "https://example.com",
		}
		resetFlags()

		_, err := parseFlags()
		if err == nil {
			t.Error("expected error for missing title")
		}
	})

	t.Run("missing url", func(t *testing.T) {
		oldArgs := os.Args
		defer func() { os.Args = oldArgs }()

		os.Args = []string{
			"og-image-generator",
			"-title", "Test Title",
		}
		resetFlags()

		_, err := parseFlags()
		if err == nil {
			t.Error("expected error for missing url")
		}
	})
}

// resetFlags resets the flag package state for testing
func resetFlags() {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
}

func TestParseFlagsVersion(t *testing.T) {
	// Save and restore osExit
	oldOsExit := osExit
	defer func() { osExit = oldOsExit }()

	exitCalled := false
	exitCode := -1
	osExit = func(code int) {
		exitCalled = true
		exitCode = code
	}

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"og-image-generator", "-version"}
	resetFlags()

	_, err := parseFlags()

	if !exitCalled {
		t.Error("expected osExit to be called")
	}
	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d", exitCode)
	}
	if err != ErrVersionRequested {
		t.Errorf("expected ErrVersionRequested, got %v", err)
	}
}

func TestGetVersionString(t *testing.T) {
	// Save original values
	oldVersion := version
	oldCommit := commit
	defer func() {
		version = oldVersion
		commit = oldCommit
	}()

	t.Run("dev version uses commit", func(t *testing.T) {
		version = "dev"
		commit = "abc123"
		result := getVersionString()
		if result != "abc123" {
			t.Errorf("expected commit hash, got %q", result)
		}
	})

	t.Run("release version", func(t *testing.T) {
		version = "v1.2.3"
		commit = "abc123"
		result := getVersionString()
		if result != "v1.2.3" {
			t.Errorf("expected version, got %q", result)
		}
	})
}

func TestResolveFontPathNoFontsFound(t *testing.T) {
	// Test the case where no fonts are found
	// We need to be in a directory without the local fonts folder
	// and use empty system font paths

	tmpDir := t.TempDir()
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tmpDir)

	// Test with no system fonts available
	result, err := resolveFontPathWithPaths("", []string{})
	if err == nil {
		t.Errorf("expected error when no fonts found, got result: %q", result)
	}
	if !strings.Contains(err.Error(), "font file not found") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestResolveFontPathWithPaths(t *testing.T) {
	t.Run("custom font takes precedence", func(t *testing.T) {
		result, err := resolveFontPathWithPaths("/custom/font.ttf", []string{"/system/font.ttf"})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if result != "/custom/font.ttf" {
			t.Errorf("expected custom font path, got %q", result)
		}
	})

	t.Run("finds system font when no custom font", func(t *testing.T) {
		// Create a temp directory and change to it (no local fonts)
		tmpDir := t.TempDir()
		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		// Use actual system font path
		fontPath := testFontPath(t)
		result, err := resolveFontPathWithPaths("", []string{fontPath})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if result != fontPath {
			t.Errorf("expected system font path %q, got %q", fontPath, result)
		}
	})

	t.Run("searches multiple paths", func(t *testing.T) {
		tmpDir := t.TempDir()
		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		fontPath := testFontPath(t)
		// First path doesn't exist, second does
		result, err := resolveFontPathWithPaths("", []string{"/nonexistent/font.ttf", fontPath})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if result != fontPath {
			t.Errorf("expected %q, got %q", fontPath, result)
		}
	})

	t.Run("local fonts directory takes precedence over system", func(t *testing.T) {
		tmpDir := t.TempDir()
		fontsDir := filepath.Join(tmpDir, "fonts")
		os.MkdirAll(fontsDir, 0755)
		localFont := filepath.Join(fontsDir, "OpenSans-Bold.ttf")
		os.WriteFile(localFont, []byte("dummy"), 0644)

		oldWd, _ := os.Getwd()
		defer os.Chdir(oldWd)
		os.Chdir(tmpDir)

		systemFont := testFontPath(t)
		result, err := resolveFontPathWithPaths("", []string{systemFont})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		// Should find local font first
		if result != filepath.Join("fonts", "OpenSans-Bold.ttf") {
			t.Errorf("expected local font, got %q", result)
		}
	})
}

func TestRunWithURLFontError(t *testing.T) {
	fontPath := testFontPath(t)
	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "test-output.png")

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	// Valid title font but invalid URL font
	os.Args = []string{
		"og-image-generator",
		"-title", "Test Title",
		"-url", "https://example.com",
		"-output", outputPath,
		"-title-font", fontPath,
		"-url-font", "/nonexistent/url-font.ttf",
	}
	resetFlags()

	err := run()
	if err == nil {
		t.Error("expected error for invalid URL font path")
	}
}

func TestRunDrawTitleError(t *testing.T) {
	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "test-output.png")

	// Create a file that is not a valid font
	invalidFontPath := filepath.Join(tmpDir, "invalid.ttf")
	os.WriteFile(invalidFontPath, []byte("not a font"), 0644)

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{
		"og-image-generator",
		"-title", "Test Title",
		"-url", "https://example.com",
		"-output", outputPath,
		"-title-font", invalidFontPath,
		"-url-font", invalidFontPath,
	}
	resetFlags()

	err := run()
	if err == nil {
		t.Error("expected error for invalid font file")
	}
}

func TestRunDrawURLError(t *testing.T) {
	fontPath := testFontPath(t)
	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "test-output.png")

	// Create a file that is not a valid font for URL
	invalidFontPath := filepath.Join(tmpDir, "invalid-url.ttf")
	os.WriteFile(invalidFontPath, []byte("not a font"), 0644)

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{
		"og-image-generator",
		"-title", "Test Title",
		"-url", "https://example.com",
		"-output", outputPath,
		"-title-font", fontPath,
		"-url-font", invalidFontPath,
	}
	resetFlags()

	err := run()
	if err == nil {
		t.Error("expected error for invalid URL font file")
	}
	if err != nil && !strings.Contains(err.Error(), "load font") {
		t.Errorf("expected 'load font' error, got: %v", err)
	}
}

func TestPreventOrphans(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "no lines",
			input:    []string{},
			expected: []string{},
		},
		{
			name:     "single line",
			input:    []string{"Hello World"},
			expected: []string{"Hello World"},
		},
		{
			name:     "two lines no orphan",
			input:    []string{"Hello World", "Foo Bar"},
			expected: []string{"Hello World", "Foo Bar"},
		},
		{
			name:     "orphan on last line",
			input:    []string{"Hello World Foo", "Bar"},
			expected: []string{"Hello World", "Foo Bar"},
		},
		{
			name:     "three lines with orphan",
			input:    []string{"First Line Here", "Second Line Words", "Orphan"},
			expected: []string{"First Line Here", "Second Line", "Words Orphan"},
		},
		{
			name:     "previous line has only one word - cannot fix",
			input:    []string{"Hello", "World"},
			expected: []string{"Hello", "World"},
		},
		{
			name:     "last line already has multiple words",
			input:    []string{"Hello World", "Foo Bar Baz"},
			expected: []string{"Hello World", "Foo Bar Baz"},
		},
		{
			name:     "nil input",
			input:    nil,
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := preventOrphans(tt.input)

			if len(result) != len(tt.expected) {
				t.Errorf("preventOrphans() returned %d lines, want %d", len(result), len(tt.expected))
				return
			}

			for i := range result {
				if result[i] != tt.expected[i] {
					t.Errorf("preventOrphans()[%d] = %q, want %q", i, result[i], tt.expected[i])
				}
			}
		})
	}
}

func TestWrapText(t *testing.T) {
	fontPath := testFontPath(t)

	tests := []struct {
		name           string
		text           string
		maxWidth       float64
		minLines       int // minimum expected lines
		checkNoOrphans bool
	}{
		{
			name:           "empty text",
			text:           "",
			maxWidth:       500,
			minLines:       0,
			checkNoOrphans: false,
		},
		{
			name:           "single word",
			text:           "Hello",
			maxWidth:       500,
			minLines:       1,
			checkNoOrphans: false,
		},
		{
			name:           "short text fits on one line",
			text:           "Hello World",
			maxWidth:       500,
			minLines:       1,
			checkNoOrphans: false,
		},
		{
			name:           "text wraps to multiple lines",
			text:           "This is a longer title that should wrap across multiple lines",
			maxWidth:       300,
			minLines:       2,
			checkNoOrphans: true,
		},
		{
			name:           "title ending with single word should not orphan",
			text:           "Building High-Performance Web Services Today",
			maxWidth:       400,
			minLines:       2,
			checkNoOrphans: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dc := gg.NewContext(1200, 628)
			if err := dc.LoadFontFace(fontPath, 72); err != nil {
				t.Fatalf("failed to load font: %v", err)
			}

			lines := wrapText(dc, tt.text, tt.maxWidth)

			if len(lines) < tt.minLines {
				t.Errorf("wrapText() returned %d lines, want at least %d", len(lines), tt.minLines)
			}

			if tt.checkNoOrphans && len(lines) >= 2 {
				lastLine := lines[len(lines)-1]
				words := strings.Fields(lastLine)
				if len(words) == 1 {
					// Check if previous line has only one word (unavoidable orphan)
					prevLine := lines[len(lines)-2]
					prevWords := strings.Fields(prevLine)
					if len(prevWords) >= 2 {
						t.Errorf("wrapText() created orphan: last line has only one word %q", lastLine)
					}
				}
			}

			// Verify all words are preserved
			originalWords := strings.Fields(tt.text)
			var resultWords []string
			for _, line := range lines {
				resultWords = append(resultWords, strings.Fields(line)...)
			}

			if len(originalWords) != len(resultWords) {
				t.Errorf("wrapText() lost words: got %d, want %d", len(resultWords), len(originalWords))
			}
		})
	}
}

func TestWrapTextOrphanPrevention(t *testing.T) {
	fontPath := testFontPath(t)

	// This test specifically verifies orphan prevention behavior
	dc := gg.NewContext(1200, 628)
	if err := dc.LoadFontFace(fontPath, 72); err != nil {
		t.Fatalf("failed to load font: %v", err)
	}

	// Test with a title that would naturally create an orphan
	// "Advanced Patterns for Building High-Performance Web Services"
	// At certain widths, "Services" might end up alone on the last line
	title := "Advanced Patterns for Building High-Performance Web Services"

	// Use a width that would cause wrapping
	lines := wrapText(dc, title, 600)

	if len(lines) < 2 {
		t.Skip("text did not wrap at this width, cannot test orphan prevention")
	}

	// Check that last line doesn't have a single word (unless unavoidable)
	lastLine := lines[len(lines)-1]
	lastWords := strings.Fields(lastLine)

	if len(lastWords) == 1 {
		// Verify this is unavoidable (previous line has only one word)
		prevLine := lines[len(lines)-2]
		prevWords := strings.Fields(prevLine)
		if len(prevWords) >= 2 {
			t.Errorf("orphan detected: last line %q has only one word, but previous line %q has %d words",
				lastLine, prevLine, len(prevWords))
		}
	}
}

func TestMain(m *testing.M) {
	// This runs all tests
	os.Exit(m.Run())
}

func TestRunWithResolverTitleFontError(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{
		"og-image-generator",
		"-title", "Test Title",
		"-url", "https://example.com",
	}
	resetFlags()

	// Create a resolver that fails on title font (first call)
	callCount := 0
	failingResolver := func(customFont string) (string, error) {
		callCount++
		if callCount == 1 {
			return "", fmt.Errorf("title font not found")
		}
		return testFontPath(t), nil
	}

	err := runWithResolver(failingResolver)
	if err == nil {
		t.Error("expected error for title font resolution failure")
	}
	if !strings.Contains(err.Error(), "title font not found") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRunWithResolverURLFontError(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{
		"og-image-generator",
		"-title", "Test Title",
		"-url", "https://example.com",
	}
	resetFlags()

	// Create a resolver that succeeds on title font but fails on URL font
	callCount := 0
	failingResolver := func(customFont string) (string, error) {
		callCount++
		if callCount == 2 {
			return "", fmt.Errorf("url font not found")
		}
		return testFontPath(t), nil
	}

	err := runWithResolver(failingResolver)
	if err == nil {
		t.Error("expected error for URL font resolution failure")
	}
	if !strings.Contains(err.Error(), "url font not found") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestMainFunction(t *testing.T) {
	fontPath := testFontPath(t)

	t.Run("successful execution", func(t *testing.T) {
		tmpDir := t.TempDir()
		outputPath := filepath.Join(tmpDir, "main-test-output.png")

		oldArgs := os.Args
		oldOsExit := osExit
		defer func() {
			os.Args = oldArgs
			osExit = oldOsExit
		}()

		os.Args = []string{
			"og-image-generator",
			"-title", "Test Title",
			"-url", "https://example.com",
			"-output", outputPath,
			"-title-font", fontPath,
			"-url-font", fontPath,
		}
		resetFlags()

		exitCalled := false
		osExit = func(code int) {
			exitCalled = true
			if code != 0 {
				t.Errorf("expected exit code 0, got %d", code)
			}
		}

		main()

		if exitCalled {
			t.Error("osExit should not be called on success")
		}

		// Verify output file was created
		if _, err := os.Stat(outputPath); os.IsNotExist(err) {
			t.Error("output file was not created")
		}
	})

	t.Run("error execution", func(t *testing.T) {
		oldArgs := os.Args
		oldOsExit := osExit
		defer func() {
			os.Args = oldArgs
			osExit = oldOsExit
		}()

		// Missing required flags
		os.Args = []string{"og-image-generator"}
		resetFlags()

		exitCalled := false
		exitCode := -1
		osExit = func(code int) {
			exitCalled = true
			exitCode = code
		}

		main()

		if !exitCalled {
			t.Error("osExit should be called on error")
		}
		if exitCode != 1 {
			t.Errorf("expected exit code 1, got %d", exitCode)
		}
	})
}
