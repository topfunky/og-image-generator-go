package main

import (
	"flag"
	"fmt"
	"image/color"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/fogleman/gg"
)

var (
	version = "dev"
	commit  = "unknown"
)

// Typographic constants
const (
	// Font sizes
	TitleFontSize  = 72.0
	URLFontSize    = 40.0
	URLMinFontSize = 16.0

	// Spacing and margins
	TextTopMargin  = 135.0
	TextSideMargin = 60.0
	LineSpacing    = 1.5
	ShadowOffset   = 2.0

	// Background
	BackgroundMargin       = 20.0
	BackgroundCornerRadius = 20.0
	BackgroundOverlayAlpha = 100
)

// Default colors
var (
	defaultBgColor = color.RGBA{26, 26, 46, 255}
	shadowColor    = color.Black
	textColor      = color.White
	mutedTextColor = color.RGBA{R: 200, G: 200, B: 200, A: 220}
	debugColor     = color.RGBA{255, 0, 0, 255}
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		osExit(1)
	}
}

// fontResolver is a function type for resolving font paths
type fontResolver func(customFont string) (string, error)

// defaultFontResolver is the default font resolver
var defaultFontResolver fontResolver = resolveFontPath

func run() error {
	return runWithResolver(defaultFontResolver)
}

func runWithResolver(resolver fontResolver) error {
	opts, err := parseFlags()
	if err != nil {
		return err
	}

	titleFontPath, err := resolver(opts.TitleFont)
	if err != nil {
		return err
	}

	urlFontPath, err := resolver(opts.URLFont)
	if err != nil {
		return err
	}

	dc := gg.NewContext(opts.Width, opts.Height)

	drawBackground(dc, opts.BgColor, opts.Width, opts.Height)

	if err := drawTitle(dc, opts.Title, titleFontPath, opts.Width); err != nil {
		return err
	}

	if opts.Debug {
		// Load font to get metrics for debug baselines
		if err := dc.LoadFontFace(titleFontPath, TitleFontSize); err != nil {
			return fmt.Errorf("load font for debug: %w", err)
		}
		fontHeight := measureFontHeight(dc)
		drawDebugBaselines(dc, fontHeight, LineSpacing, TextTopMargin, opts.Width, opts.Height)
	}

	if err := drawURL(dc, opts.URL, titleFontPath, urlFontPath, opts.Width, opts.Height); err != nil {
		return err
	}

	if err := dc.SavePNG(opts.Output); err != nil {
		return fmt.Errorf("save png: %w", err)
	}

	fmt.Printf("Social image generated: %s\n", opts.Output)
	return nil
}

// Options holds the configuration for image generation
type Options struct {
	Title     string
	URL       string
	Output    string
	Width     int
	Height    int
	BgColor   string
	TitleFont string
	URLFont   string
	Debug     bool
}

// ErrVersionRequested is returned when the -version flag is passed
var ErrVersionRequested = fmt.Errorf("version requested")

// osExit is a variable to allow testing of os.Exit calls
var osExit = os.Exit

func parseFlags() (*Options, error) {
	title := flag.String("title", "", "Article title (required)")
	url := flag.String("url", "", "Article URL (required)")
	output := flag.String("output", "social-image.png", "Output file path")
	width := flag.Int("width", 1200, "Image width in pixels")
	height := flag.Int("height", 628, "Image height in pixels")
	bgColor := flag.String("bg", "#1a1a2e", "Background color (hex)")
	titleFont := flag.String("title-font", "", "Title font file path (TTF)")
	urlFont := flag.String("url-font", "", "URL font file path (TTF)")
	versionFlag := flag.Bool("version", false, "Print version and exit")
	debug := flag.Bool("debug", false, "Draw debug baselines")

	flag.Parse()

	if *versionFlag {
		fmt.Println("og-image-generator version " + getVersionString())
		osExit(0)
		return nil, ErrVersionRequested
	}

	if *title == "" || *url == "" {
		flag.PrintDefaults()
		return nil, fmt.Errorf("title and url are required")
	}

	return &Options{
		Title:     *title,
		URL:       *url,
		Output:    *output,
		Width:     *width,
		Height:    *height,
		BgColor:   *bgColor,
		TitleFont: *titleFont,
		URLFont:   *urlFont,
		Debug:     *debug,
	}, nil
}

func getVersionString() string {
	if version == "dev" {
		return commit
	}
	return version
}

// defaultSystemFontPaths contains the default system font paths to search
var defaultSystemFontPaths = []string{
	"/System/Library/Fonts/SFCompact.ttf",
	"/System/Library/Fonts/SFNSDisplay.ttf",
	"/System/Library/Fonts/Arial.ttf",
	"/usr/share/fonts/truetype/dejavu/DejaVuSans-Bold.ttf",
	"/usr/share/fonts/truetype/liberation/LiberationSans-Bold.ttf",
	"C:\\Windows\\Fonts\\arial.ttf",
}

func resolveFontPath(customFont string) (string, error) {
	return resolveFontPathWithPaths(customFont, defaultSystemFontPaths)
}

func resolveFontPathWithPaths(customFont string, systemPaths []string) (string, error) {
	if customFont != "" {
		return customFont, nil
	}

	fontPath := filepath.Join("fonts", "OpenSans-Bold.ttf")
	if _, err := os.Stat(fontPath); err == nil {
		return fontPath, nil
	}

	for _, p := range systemPaths {
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
	}

	return "", fmt.Errorf("font file not found at %s and no system fonts found. Please provide a TTF font file in the fonts/ directory", fontPath)
}

func drawBackground(dc *gg.Context, bgColorStr string, width, height int) {
	bgRGB := hexToRGB(bgColorStr)
	dc.SetColor(bgRGB)
	dc.Clear()

	dc.SetColor(color.RGBA{0, 0, 0, BackgroundOverlayAlpha})
	drawRoundedTopRect(dc, BackgroundMargin, BackgroundMargin, float64(width)-(2*BackgroundMargin), float64(height)-(2*BackgroundMargin), BackgroundCornerRadius)
	dc.Fill()
}

// drawRoundedTopRect draws a rectangle with rounded corners on top and square corners on bottom
func drawRoundedTopRect(dc *gg.Context, x, y, w, h, radius float64) {
	// Start at bottom-left corner (square)
	dc.MoveTo(x, y+h)
	// Line to bottom-right corner (square)
	dc.LineTo(x+w, y+h)
	// Line up to where top-right curve starts
	dc.LineTo(x+w, y+radius)
	// Top-right rounded corner
	dc.DrawArc(x+w-radius, y+radius, radius, 0, -gg.Radians(90))
	// Line to where top-left curve starts
	dc.LineTo(x+radius, y)
	// Top-left rounded corner
	dc.DrawArc(x+radius, y+radius, radius, gg.Radians(270), gg.Radians(180))
	// Close path back to bottom-left
	dc.LineTo(x, y+h)
	dc.ClosePath()
}

// wrapText wraps text to fit within maxWidth and prevents orphans.
// An orphan is when the last line contains only one word.
// If an orphan is detected, the last word from the previous line is moved
// to the last line so the final line has at least two words.
func wrapText(dc *gg.Context, text string, maxWidth float64) []string {
	words := strings.Fields(text)
	if len(words) == 0 {
		return nil
	}

	var lines []string
	var currentLine strings.Builder

	for _, word := range words {
		testLine := currentLine.String()
		if testLine != "" {
			testLine += " "
		}
		testLine += word

		w, _ := dc.MeasureString(testLine)
		if w > maxWidth && currentLine.Len() > 0 {
			lines = append(lines, currentLine.String())
			currentLine.Reset()
			currentLine.WriteString(word)
		} else {
			if currentLine.Len() > 0 {
				currentLine.WriteString(" ")
			}
			currentLine.WriteString(word)
		}
	}
	if currentLine.Len() > 0 {
		lines = append(lines, currentLine.String())
	}

	return preventOrphans(lines)
}

// preventOrphans checks if the last line has only one word and if so,
// moves the last word from the previous line to create a more balanced layout.
// After fixing an orphan, it also checks if the line before the modified line
// can be balanced by moving a word down.
func preventOrphans(lines []string) []string {
	if len(lines) < 2 {
		return lines
	}

	lastLine := lines[len(lines)-1]
	lastLineWords := strings.Fields(lastLine)

	// Only fix if last line has exactly one word (orphan)
	if len(lastLineWords) != 1 {
		return lines
	}

	prevLine := lines[len(lines)-2]
	prevLineWords := strings.Fields(prevLine)

	// Only move a word if the previous line has at least 2 words
	if len(prevLineWords) < 2 {
		return lines
	}

	// Move the last word from previous line to the last line
	wordToMove := prevLineWords[len(prevLineWords)-1]
	newPrevLine := strings.Join(prevLineWords[:len(prevLineWords)-1], " ")
	newLastLine := wordToMove + " " + lastLine

	lines[len(lines)-2] = newPrevLine
	lines[len(lines)-1] = newLastLine

	// Now check if we need to balance lines above the modified line
	return balanceLinesUpward(lines, len(lines)-2)
}

// balanceLinesUpward checks if the line at modifiedIdx can be balanced with the line above it.
// If the line above ends with two words that both start after the length of the modified line,
// move one word down to balance. This process continues upward as needed.
func balanceLinesUpward(lines []string, modifiedIdx int) []string {
	if modifiedIdx < 1 {
		return lines
	}

	for idx := modifiedIdx; idx >= 1; idx-- {
		currentLine := lines[idx]
		aboveLine := lines[idx-1]

		currentLen := len(currentLine)
		aboveWords := strings.Fields(aboveLine)

		// Need at least 2 words in the line above to consider balancing
		if len(aboveWords) < 2 {
			continue
		}

		// Check if the last two words of the line above both start after the current line's length
		lineWithoutLastWord := strings.Join(aboveWords[:len(aboveWords)-1], " ")
		secondToLastWordStart := len(strings.Join(aboveWords[:len(aboveWords)-2], " "))
		if len(aboveWords) > 2 {
			secondToLastWordStart++ // account for space before the word
		}

		// If the second-to-last word starts at or after the current line's length,
		// both trailing words are "hanging" past the current line, so move one down
		if secondToLastWordStart >= currentLen {
			wordToMove := aboveWords[len(aboveWords)-1]
			lines[idx-1] = lineWithoutLastWord
			lines[idx] = wordToMove + " " + currentLine
			// Continue checking upward since we modified line idx-1
		} else {
			// No balancing needed at this level, stop propagating
			break
		}
	}

	return lines
}

// drawTextWithShadow draws text with a shadow effect at the specified position
func drawTextWithShadow(dc *gg.Context, text string, x, y float64) {
	// Draw shadow
	dc.SetColor(shadowColor)
	dc.DrawString(text, x+ShadowOffset, y+ShadowOffset)

	// Draw text
	dc.SetColor(textColor)
	dc.DrawString(text, x, y)
}

func drawTitle(dc *gg.Context, title, fontPath string, width int) error {
	if err := dc.LoadFontFace(fontPath, TitleFontSize); err != nil {
		return fmt.Errorf("load font: %w", err)
	}

	maxWidth := float64(width) - (2 * TextSideMargin)
	lines := wrapText(dc, title, maxWidth)

	fontHeight := measureFontHeight(dc)
	verticalOffset := fontHeight

	for i, line := range lines {
		y := TextTopMargin + float64(i)*fontHeight*LineSpacing + verticalOffset
		drawTextWithShadow(dc, line, TextSideMargin, y)
	}

	return nil
}

func drawURL(dc *gg.Context, url string, titleFontPath string, urlFontPath string, width, height int) error {
	maxWidth := float64(width) - (2 * TextSideMargin)

	// Find the appropriate font size that fits the URL
	urlFontSize := URLFontSize
	for urlFontSize >= URLMinFontSize {
		if err := dc.LoadFontFace(urlFontPath, urlFontSize); err != nil {
			return fmt.Errorf("load font for url: %w", err)
		}

		textWidth, _ := dc.MeasureString(url)
		if textWidth <= maxWidth {
			break
		}
		urlFontSize -= 2.0
	}

	// Ensure font is loaded at final size
	if err := dc.LoadFontFace(urlFontPath, urlFontSize); err != nil {
		return fmt.Errorf("load font for url: %w", err)
	}

	dc.SetColor(mutedTextColor)

	// Calculate the baseline grid using the title font metrics
	titleFontHeight, err := getFontHeight(titleFontPath, TitleFontSize, width, height)
	if err != nil {
		return fmt.Errorf("load title font for baseline: %w", err)
	}

	// Find the last baseline that fits within the image bounds
	// The baseline grid starts at TextTopMargin + titleFontHeight (first baseline)
	// and increments by titleFontHeight * LineSpacing
	// Leave space equal to TextTopMargin at the bottom of the image
	firstBaseline := TextTopMargin + titleFontHeight
	baselineStep := titleFontHeight * LineSpacing
	maxY := float64(height) - TextTopMargin/2.0

	// Find the last baseline that doesn't exceed the bottom margin
	targetY := firstBaseline
	for y := firstBaseline; y <= maxY; y += baselineStep {
		targetY = y
	}

	dc.DrawString(url, TextSideMargin, targetY)

	return nil
}

// getFontHeight returns the height of a font at a given size
func getFontHeight(fontPath string, fontSize float64, width, height int) (float64, error) {
	tempDc := gg.NewContext(width, height)
	if err := tempDc.LoadFontFace(fontPath, fontSize); err != nil {
		return 0, err
	}
	return measureFontHeight(tempDc), nil
}

// drawDebugBaselines draws hairline red lines at each typographic baseline
func drawDebugBaselines(dc *gg.Context, fontHeight, lineSpacing, textTopMargin float64, width, height int) {
	dc.SetColor(debugColor)
	dc.SetLineWidth(2)

	firstBaseline := textTopMargin + fontHeight

	// Draw top margin line
	dc.DrawLine(0, textTopMargin, float64(width), textTopMargin)
	dc.Stroke()

	// Draw baselines at each line height interval until we reach the bottom
	for y := firstBaseline; y < float64(height); y += fontHeight * lineSpacing {
		roundedY := math.Round(y*2) / 2
		dc.DrawLine(0, roundedY, float64(width), roundedY)
		dc.Stroke()
	}
}

// measureFontHeight returns the height of the currently loaded font.
// It uses "Mg" as reference characters to capture both ascenders and descenders.
func measureFontHeight(dc *gg.Context) float64 {
	_, height := dc.MeasureString("Mg")
	return height
}

// hexToRGB converts hex color string to color.RGBA
func hexToRGB(hexColor string) color.Color {
	hexColor = strings.TrimPrefix(hexColor, "#")
	if len(hexColor) != 6 {
		return defaultBgColor
	}

	val, err := strconv.ParseUint(hexColor, 16, 32)
	if err != nil {
		return defaultBgColor
	}

	return color.RGBA{
		R: uint8(val >> 16),
		G: uint8(val >> 8),
		B: uint8(val),
		A: 255,
	}
}
