# OG Image Generator

A Go application for generating dynamic social share images (Open Graph images) with customizable titles and URLs.

Based on the approach from [Mat Ryer's blog post](https://pace.dev/blog/2020/03/02/dynamically-generate-social-images-in-golang-by-mat-ryer.html) on programmatically generating beautiful social media images in Go.

![OpenGraph social image](https://topfunky.com/2026/easy-tools-you-should-use/og-image.out.png)

## Features

- **Optimized Dimensions**: Generates 1200x628px images optimized for Twitter/Facebook (configurable)
- **Typography**: Applies a text layout algorithm to ensure that text is balanced (no orphans; single words on the last visible line)
- **Baseline grid**: Places text on a consistent grid, calculated from actual title font, size, and line height
- **Text Rendering**: Displays article titles with text shadows for improved readability
- **Visual Design**: Semi-transparent overlays and customizable background colors
- **Automatic font sizing**: URLs of any length will be sized to fit the card dimensions
- **Responsive Layout**: Text wrapping and positioning works across different image sizes
- **System Font Fallback**: Automatically uses available system fonts if custom fonts aren't provided
- **PNG Output**: Generates high-quality PNG files ready for web use

## Installation

### Homebrew (macOS and Linux)

```bash
brew tap topfunky/tap
brew install topfunky/tap/og-image-generator
```

### Build from Source

#### Prerequisites

- Go 1.21 or later

```bash
go mod download
go build -o og-image-generator
```

## Usage

### Basic Command

```bash
./og-image-generator -title "Your Article Title" -url "https://example.com/article"
```

### All Options

```bash
./og-image-generator \
  -title "Article Title" \
  -url "https://example.com/article" \
  -output social-image.png \
  -width 1200 \
  -height 628 \
  -bg "#1a1a2e"
```

Alternate example

```bash
./og-image-generator \
	-title "This is a test many lines of text one after the other" \
  -title-font fonts/MonaspaceKryptonFrozen-ExtraBold.ttf \
	-url "https://example.com" \
	-url-font fonts/MonaspaceKryptonFrozen-Medium.ttf \
	-output out/krypton.png \
	-bg "#002052"
```

### Command-Line Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-title` | *required* | Article title to display on image |
| `-url` | *required* | Article URL to display at bottom |
| `-output` | `social-image.png` | Output file path |
| `-width` | `1200` | Image width in pixels |
| `-height` | `628` | Image height in pixels |
| `-bg` | `#1a1a2e` | Background color in hex format (e.g., `#2c3e50`) |

### Examples

**Basic usage:**
```bash
./og-image-generator \
  -title "How to Build APIs in Go" \
  -url "https://example.com/go-apis"
```

**Custom output and colors:**
```bash
./og-image-generator \
  -title "Mastering Concurrency" \
  -url "https://example.com/concurrency" \
  -output my-image.png \
  -bg "#2c3e50"
```

**Custom dimensions (16:9):**
```bash
./og-image-generator \
  -title "My Blog Post" \
  -url "https://myblog.com/post" \
  -output wide-image.png \
  -width 1600 \
  -height 900
```

**Using a long title:**
```bash
./og-image-generator \
  -title "Advanced Patterns for Building High-Performance, Scalable, and Maintainable Systems" \
  -url "https://example.com/patterns" \
  -output pattern-image.png
```

## Font Configuration

### Using Custom Fonts

To use a custom TrueType font:

1. **Place font in the `fonts/` directory:**
   ```bash
   cp ~/Downloads/MyFont-Bold.ttf ./fonts/
   ```

1. **Edit `main.go` and update the `fontPath` variable:**
   ```go
   fontPath := filepath.Join("fonts", "MyFont-Bold.ttf")
   ```

1. **Rebuild the application:**
   ```bash
   go build -o og-image-generator
   ```

### System Fonts

The application automatically uses system fonts if a custom font isn't provided. It checks for fonts in this order:
- `/System/Library/Fonts/SFCompact.ttf` (macOS)
- `/System/Library/Fonts/SFNSDisplay.ttf` (macOS)
- `/usr/share/fonts/truetype/dejavu/DejaVuSans-Bold.ttf` (Linux)
- `/usr/share/fonts/truetype/liberation/LiberationSans-Bold.ttf` (Linux)
- `C:\Windows\Fonts\arial.ttf` (Windows)

### Recommended Free Fonts

Download from [Google Fonts](https://fonts.google.com/):
- **Open Sans** - Clean, readable sans-serif
- **Roboto** - Modern and versatile
- **Lato** - Warm and friendly

## Generated Image Format

Each generated image includes:

- **Background**: Customizable solid color (hex format)
- **Overlay**: Semi-transparent dark rectangle for contrast
- **Title**: Large text with shadow effect, supports text wrapping
- **URL**: Medium text centered at bottom
- **Branding**: "OG Image" text at bottom right

## Integration

### HTML Meta Tags

Add to your HTML `<head>` to use the generated image for social sharing:

```html
<meta property="og:image" content="/path/to/social-image.png">
<meta property="og:image:width" content="1200">
<meta property="og:image:height" content="628">
<meta name="twitter:image" content="/path/to/social-image.png">
<meta name="twitter:card" content="summary_large_image">
```

### Go Web Server Integration

```go
package main

import (
	"fmt"
	"net/http"
	"os/exec"
)

func generateImage(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Query().Get("title")
	url := r.URL.Query().Get("url")
	
	cmd := exec.Command("./og-image-generator",
		"-title", title,
		"-url", url,
		"-output", "social-image.png")
	
	if err := cmd.Run(); err != nil {
		http.Error(w, "Failed to generate image", http.StatusInternalServerError)
		return
	}
	
	http.ServeFile(w, r, "social-image.png")
}

func main() {
	http.HandleFunc("/generate-image", generateImage)
	http.ListenAndServe(":8080", nil)
}
```

## Dependencies

- `github.com/fogleman/gg` - 2D graphics library providing a simple API on top of the Go standard library

## Performance Notes

- Image generation is fast (typically < 100ms)
- PNG encoding adds ~10-20ms
- File sizes typically 30-50KB for standard 1200x628 images
- Suitable for on-demand generation in web services

## Troubleshooting

### Font not found error
- Ensure a TTF font file is in the `fonts/` directory, OR
- Install system fonts (DejaVu Sans on Linux, San Francisco fonts on macOS)
- Check that the font path in `main.go` is correct

### Text is cut off or wrapping strangely
- Use the `-width` and `-height` flags to adjust image dimensions
- Reduce title length or increase font size by modifying the code

### Generated image looks blurry
- Ensure you're using a proper TTF font file
- Avoid very large font sizes on small images
- Save as PNG for lossless compression

## License

This example is provided for educational purposes based on Mat Ryer's blog post.

## Further Reading

- [Mat Ryer's original blog post](https://pace.dev/blog/2020/03/02/dynamically-generate-social-images-in-golang-by-mat-ryer.html)
- [fogleman/gg documentation](https://github.com/fogleman/gg)
- [Twitter Card documentation](https://developer.twitter.com/en/docs/twitter-for-websites/cards/overview/abouts-cards)
- [Open Graph protocol](https://ogp.me/)
