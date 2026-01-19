// Example: Pre-render emoji as images and composite using dc.DrawImage()
//
// This example demonstrates how to include emoji in gg-rendered images by:
// 1. Loading pre-rendered emoji PNG files
// 2. Compositing them inline with text
//
// To run this example:
//   1. Download emoji PNGs (e.g., from https://github.com/twitter/twemoji)
//   2. Place them in an "emoji" subdirectory (e.g., emoji/1f3c8.png for 🏈)
//   3. Run: go run main.go
//
// Emoji files are typically named by their Unicode codepoint in hex.
// For example: 🏈 (U+1F3C8) -> 1f3c8.png

package main

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/fogleman/gg"
	"golang.org/x/image/draw"
)

const (
	imageWidth  = 800
	imageHeight = 400
	fontSize    = 48.0
	emojiSize   = 48 // emoji will be scaled to match font size
)

// EmojiMap maps shortcodes to Unicode codepoints (hex, lowercase)
var emojiMap = map[string]string{
	":football:":   "1f3c8",
	":soccer:":     "26bd",
	":basketball:": "1f3c0",
	":baseball:":   "26be",
	":tennis:":     "1f3be",
	":star:":       "2b50",
	":fire:":       "1f525",
	":heart:":      "2764",
	":thumbsup:":   "1f44d",
	":rocket:":     "1f680",
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	// Create the drawing context
	dc := gg.NewContext(imageWidth, imageHeight)

	// Fill background
	dc.SetRGB(0.1, 0.1, 0.2)
	dc.Clear()

	// Load a font (adjust path as needed for your system)
	fontPath := findFont()
	if err := dc.LoadFontFace(fontPath, fontSize); err != nil {
		return fmt.Errorf("load font: %w", err)
	}

	// Example text with emoji shortcodes
	text := "Game day! :football: Let's go! :fire:"

	// Draw text with inline emoji
	x := 50.0
	y := 200.0

	if err := drawTextWithEmoji(dc, text, x, y, fontPath, fontSize); err != nil {
		return fmt.Errorf("draw text with emoji: %w", err)
	}

	// Save the result
	output := "output.png"
	if err := dc.SavePNG(output); err != nil {
		return fmt.Errorf("save png: %w", err)
	}

	fmt.Printf("Generated: %s\n", output)
	return nil
}

// drawTextWithEmoji renders text with emoji shortcodes replaced by images
func drawTextWithEmoji(dc *gg.Context, text string, x, y float64, fontPath string, fontSize float64) error {
	// Ensure font is loaded
	if err := dc.LoadFontFace(fontPath, fontSize); err != nil {
		return err
	}

	dc.SetRGB(1, 1, 1) // white text

	// Parse text into segments (text and emoji)
	segments := parseTextWithEmoji(text)

	currentX := x
	for _, seg := range segments {
		if seg.isEmoji {
			// Load and draw emoji image
			emojiImg, err := loadEmoji(seg.codepoint)
			if err != nil {
				// If emoji not found, draw the shortcode as text
				fmt.Printf("Warning: emoji %s not found, using text fallback\n", seg.content)
				dc.DrawString(seg.content, currentX, y)
				w, _ := dc.MeasureString(seg.content)
				currentX += w
				continue
			}

			// Scale emoji to match font size
			scaled := scaleImage(emojiImg, emojiSize, emojiSize)

			// Draw emoji aligned with text baseline
			// Adjust Y position: DrawImage uses top-left corner,
			// but text Y is the baseline, so move up by ~80% of emoji height
			emojiY := y - float64(emojiSize)*0.8
			dc.DrawImage(scaled, int(currentX), int(emojiY))
			currentX += float64(emojiSize) + 4 // small gap after emoji
		} else {
			// Draw regular text
			dc.DrawString(seg.content, currentX, y)
			w, _ := dc.MeasureString(seg.content)
			currentX += w
		}
	}

	return nil
}

// segment represents a piece of text or an emoji
type segment struct {
	content   string // original text or shortcode
	isEmoji   bool
	codepoint string // Unicode codepoint if emoji
}

// parseTextWithEmoji splits text into segments of plain text and emoji shortcodes
func parseTextWithEmoji(text string) []segment {
	var segments []segment
	remaining := text

	for len(remaining) > 0 {
		// Find next emoji shortcode
		startIdx := strings.Index(remaining, ":")
		if startIdx == -1 {
			// No more colons, rest is plain text
			if remaining != "" {
				segments = append(segments, segment{content: remaining, isEmoji: false})
			}
			break
		}

		// Add text before the colon
		if startIdx > 0 {
			segments = append(segments, segment{content: remaining[:startIdx], isEmoji: false})
		}

		// Look for closing colon
		endIdx := strings.Index(remaining[startIdx+1:], ":")
		if endIdx == -1 {
			// No closing colon, rest is plain text
			segments = append(segments, segment{content: remaining[startIdx:], isEmoji: false})
			break
		}
		endIdx += startIdx + 1 // adjust for offset

		// Extract potential shortcode
		shortcode := remaining[startIdx : endIdx+1]

		// Check if it's a known emoji
		if codepoint, ok := emojiMap[shortcode]; ok {
			segments = append(segments, segment{
				content:   shortcode,
				isEmoji:   true,
				codepoint: codepoint,
			})
		} else {
			// Not a known emoji, treat as plain text
			segments = append(segments, segment{content: shortcode, isEmoji: false})
		}

		remaining = remaining[endIdx+1:]
	}

	return segments
}

// loadEmoji loads an emoji PNG from the emoji directory
func loadEmoji(codepoint string) (image.Image, error) {
	// Try common emoji asset locations
	paths := []string{
		filepath.Join("emoji", codepoint+".png"),
		filepath.Join("assets", "emoji", codepoint+".png"),
		filepath.Join("twemoji", codepoint+".png"),
	}

	for _, path := range paths {
		f, err := os.Open(path)
		if err != nil {
			continue
		}
		defer f.Close()

		img, err := png.Decode(f)
		if err != nil {
			return nil, fmt.Errorf("decode emoji png: %w", err)
		}
		return img, nil
	}

	return nil, fmt.Errorf("emoji file not found for codepoint %s", codepoint)
}

// scaleImage scales an image to the specified dimensions
func scaleImage(src image.Image, width, height int) image.Image {
	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.CatmullRom.Scale(dst, dst.Bounds(), src, src.Bounds(), draw.Over, nil)
	return dst
}

// findFont returns a path to an available system font
func findFont() string {
	paths := []string{
		"../../fonts/OpenSans-Bold.ttf",
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

	return paths[0] // fallback, will error if not found
}
