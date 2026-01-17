package main

import (
	"flag"
	"fmt"
	"image/color"
	"os"
	"path/filepath"
	"strings"

	"github.com/fogleman/gg"
)

var (
	version = "dev"
	commit  = "unknown"
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
		if err := dc.LoadFontFace(titleFontPath, 72); err != nil {
			return fmt.Errorf("load font for debug: %w", err)
		}
		_, fontHeight := dc.MeasureString("Mg")
		textTopMargin := 90.0
		lineSpacing := 1.5
		drawDebugBaselines(dc, fontHeight, lineSpacing, textTopMargin, opts.Width, opts.Height)
	}

	if err := drawURL(dc, opts.URL, urlFontPath, opts.Width, opts.Height); err != nil {
		return err
	}

	if err := dc.SavePNG(opts.Output); err != nil {
		return fmt.Errorf("save png: %w", err)
	}

	fmt.Printf("Social image generated: %s\n", opts.Output)
	return nil
}

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

	margin := 20.0
	dc.SetColor(color.RGBA{0, 0, 0, 100})
	drawRoundedTopRect(dc, margin, margin, float64(width)-(2*margin), float64(height)-(2*margin), 20.0)
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
	var currentLine string

	for _, word := range words {
		testLine := currentLine
		if testLine != "" {
			testLine += " "
		}
		testLine += word

		w, _ := dc.MeasureString(testLine)
		if w > maxWidth && currentLine != "" {
			lines = append(lines, currentLine)
			currentLine = word
		} else {
			currentLine = testLine
		}
	}
	if currentLine != "" {
		lines = append(lines, currentLine)
	}

	// Prevent orphans: if last line has only one word and there's a previous line,
	// move the last word from the previous line to the last line
	lines = preventOrphans(lines)

	return lines
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
	// Work backwards from the modified line (len(lines)-2)
	lines = balanceLinesUpward(lines, len(lines)-2)

	return lines
}

// balanceLinesUpward checks if the line at modifiedIdx can be balanced with the line above it.
// If the line above ends with two words that both start after the length of the modified line,
// move one word down to balance. This process continues upward as needed.
func balanceLinesUpward(lines []string, modifiedIdx int) []string {
	// Need at least a line above the modified line
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
		// Build the line without the last word to find where the second-to-last word starts
		lineWithoutLastWord := strings.Join(aboveWords[:len(aboveWords)-1], " ")
		secondToLastWordStart := len(strings.Join(aboveWords[:len(aboveWords)-2], " "))
		if len(aboveWords) > 2 {
			secondToLastWordStart++ // account for space before the word
		}

		// If the second-to-last word starts at or after the current line's length,
		// both trailing words are "hanging" past the current line, so move one down
		if secondToLastWordStart >= currentLen {
			// Move the last word from the line above to the current line
			wordToMove := aboveWords[len(aboveWords)-1]
			newAboveLine := lineWithoutLastWord
			newCurrentLine := wordToMove + " " + currentLine

			lines[idx-1] = newAboveLine
			lines[idx] = newCurrentLine
			// Continue checking upward since we modified line idx-1
		} else {
			// No balancing needed at this level, stop propagating
			break
		}
	}

	return lines
}

func drawTitle(dc *gg.Context, title, fontPath string, width int) error {
	if err := dc.LoadFontFace(fontPath, 72); err != nil {
		return fmt.Errorf("load font: %w", err)
	}

	textRightMargin := 60.0
	textTopMargin := 90.0
	maxWidth := float64(width) - (2 * textRightMargin)
	lineSpacing := 1.5

	lines := wrapText(dc, title, maxWidth)

	// Get font metrics for line height calculation
	_, fontHeight := dc.MeasureString("Mg") // Use typical characters for height
	verticalOffset := fontHeight

	// Draw shadow
	dc.SetColor(color.Black)
	for i, line := range lines {
		y := textTopMargin + 2 + float64(i)*fontHeight*lineSpacing + verticalOffset
		dc.DrawString(line, textRightMargin+2, y)
	}

	// Draw text
	dc.SetColor(color.White)
	for i, line := range lines {
		y := textTopMargin + float64(i)*fontHeight*lineSpacing + verticalOffset
		dc.DrawString(line, textRightMargin, y)
	}

	return nil
}

func drawURL(dc *gg.Context, url, fontPath string, width, height int) error {
	maxWidth := float64(width) - 120.0
	fontSize := 40.0
	minFontSize := 16.0

	for fontSize >= minFontSize {
		if err := dc.LoadFontFace(fontPath, fontSize); err != nil {
			return fmt.Errorf("load font for url: %w", err)
		}

		textWidth, _ := dc.MeasureString(url)
		if textWidth <= maxWidth {
			break
		}

		fontSize -= 2.0
	}

	mutedColor := color.RGBA{R: 200, G: 200, B: 200, A: 220}
	dc.SetColor(mutedColor)

	// Align URL to the typographic baseline grid established by the title
	// Title uses: textTopMargin=90, fontHeight from 72pt, lineSpacing=1.5
	// We need to temporarily load the title font to get its metrics
	titleFontHeight := 72.0 * 1.2 // Approximate font height for 72pt (ascent + descent)
	textTopMargin := 90.0
	lineSpacing := 1.5
	verticalOffset := titleFontHeight
	firstBaseline := textTopMargin + verticalOffset

	// Find the baseline closest to the bottom of the image (with some margin)
	targetY := float64(height) - 70.0
	baselineInterval := titleFontHeight * lineSpacing

	// Calculate which baseline number we're closest to
	n := (targetY - firstBaseline) / baselineInterval
	// Round to nearest baseline and go up one to ensure it's above the target
	nearestBaseline := firstBaseline + float64(int(n))*baselineInterval

	dc.DrawString(url, 60.0, nearestBaseline)

	return nil
}

// drawDebugBaselines draws hairline red lines at each typographic baseline
func drawDebugBaselines(dc *gg.Context, fontHeight, lineSpacing, textTopMargin float64, width, height int) {
	dc.SetColor(color.RGBA{255, 0, 0, 255}) // Red
	dc.SetLineWidth(1)                      // Hairline

	verticalOffset := fontHeight
	firstBaseline := textTopMargin + verticalOffset

	// Draw baselines at each line height interval until we reach the bottom
	for y := firstBaseline; y < float64(height); y += fontHeight * lineSpacing {
		dc.DrawLine(0, y, float64(width), y)
		dc.Stroke()
	}
}

// hexToRGB converts hex color string to color.RGBA
func hexToRGB(hexColor string) color.Color {
	hexColor = strings.TrimPrefix(hexColor, "#")
	if len(hexColor) != 6 {
		return color.RGBA{26, 26, 46, 255} // default dark blue
	}

	var r, g, b uint8
	_, err := fmt.Sscanf(hexColor, "%02x%02x%02x", &r, &g, &b)
	if err != nil {
		return color.RGBA{26, 26, 46, 255} // default on error
	}

	return color.RGBA{r, g, b, 255}
}
