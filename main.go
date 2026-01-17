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
	dc.DrawRectangle(margin, margin, float64(width)-(2*margin), float64(height)-(2*margin))
	dc.Fill()
}

func drawTitle(dc *gg.Context, title, fontPath string, width int) error {
	if err := dc.LoadFontFace(fontPath, 72); err != nil {
		return fmt.Errorf("load font: %w", err)
	}

	dc.SetColor(color.Black)
	textRightMargin := 60.0
	textTopMargin := 90.0
	maxWidth := float64(width) - (2 * textRightMargin)

	dc.DrawStringWrapped(title, textRightMargin+2, textTopMargin+2, 0, 0, maxWidth, 1.5, gg.AlignLeft)

	dc.SetColor(color.White)
	dc.DrawStringWrapped(title, textRightMargin, textTopMargin, 0, 0, maxWidth, 1.5, gg.AlignLeft)

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
	urlY := float64(height) - 70.0
	dc.DrawString(url, 60.0, urlY)

	return nil
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
