# Architecture & Implementation Notes

## Overview

This application generates social media images (Open Graph images) by combining:
- Command-line argument parsing for configuration
- Image rendering using the `gg` graphics library
- TrueType font rasterization
- PNG encoding and export

## Design Decisions

### Image Format
- **Resolution**: 1200x628 pixels (Twitter/Facebook standard)
- **Format**: PNG (lossless, widely supported)
- **Color Model**: RGB (8-bit per channel)

### Typography
- **Title Font**: Large bold font (72pt) for article titles
- **URL Font**: Medium font (40pt) for article URLs
- **Branding Font**: Medium font (50pt) for "OG Image" branding
- **Text Shadow**: 2px black offset for improved readability over varied backgrounds

### Visual Design
- **Background**: Solid color (customizable via hex)
- **Overlay**: Semi-transparent black rectangle (80% opacity) with 20px margin
- **Text Colors**: 
  - Title/Branding: White
  - URL: Light gray (200, 200, 200) at 220/255 opacity
  - Shadow: Black

### Layout
- **Margins**: 60px horizontal, 90px top for title
- **Text Wrapping**: Line height 1.5x, left-aligned
- **URL Positioning**: Centered horizontally, 40px from bottom
- **Branding**: Bottom right corner, 40px from right, 20px from bottom

## Technical Architecture

```
Command Line Arguments
    ↓
Parse Flags (title, url, output, dimensions, color)
    ↓
Create gg.Context
    ↓
Set Background Color & Clear
    ↓
Draw Overlay Rectangle
    ↓
Load Font (custom or system fallback)
    ↓
Render Title (shadow + text)
    ↓
Render URL
    ↓
Render Branding
    ↓
Save PNG
```

## Code Structure

### `main()`
Entry point that calls `run()` and handles exit codes

### `run()`
Main logic that:
1. Parses command-line flags
2. Validates required arguments
3. Initializes graphics context
4. Loads font with fallback system support
5. Renders all text elements
6. Saves output PNG

### `hexToRGB()`
Utility function to convert hex color strings to `color.RGBA`

## Font System

The application uses a three-tier font resolution strategy:

1. **Custom Fonts**: Check `fonts/OpenSans-Bold.ttf` in working directory
2. **System Fonts**: Try OS-specific font locations
3. **Error**: Fail with helpful message if no font found

This design ensures:
- Works on macOS, Linux, and Windows
- No external dependencies beyond Go standard library fonts
- Users can easily customize fonts

## Performance Characteristics

- **Font Loading**: ~10-20ms (system fonts faster than network fonts)
- **Text Rendering**: ~20-30ms per text element
- **PNG Encoding**: ~10-20ms
- **Total**: Typically <100ms for image generation

## Extensibility

Potential enhancements:

1. **GIF Support**: Render multiple frames and encode as animated GIF
2. **Image Backgrounds**: Load and composite background images
3. **Custom Logos**: Add organization branding/logos
4. **Multi-line Formatting**: Support additional text fields (author, category)
5. **Gradient Backgrounds**: Replace solid colors with gradients
6. **Template System**: Load SVG or other template formats
7. **Batch Processing**: Generate images for multiple articles

## Dependencies

### Direct
- `github.com/fogleman/gg` (v1.3.0): High-level graphics library

### Transitive
- `github.com/golang/freetype`: Font rasterization
- `golang.org/x/image`: Low-level image manipulation

## Known Limitations

1. **Font Format**: Only supports TrueType (.ttf) and OpenType (.otf) fonts
   - Some system fonts (.ttc, .dfont) not directly supported
2. **Text Overflow**: Long titles that don't fit are truncated, not wrapped off-screen
3. **Unicode**: Relies on font support for non-ASCII characters
4. **Color Format**: Only supports 6-digit hex colors (e.g., #RRGGBB)

## Testing Recommendations

When extending this application:

1. **Font Testing**: Test with various font files
2. **Title Length**: Test long titles, special characters
3. **URL Validation**: Test with various URL formats
4. **Color Formats**: Test hex color parsing edge cases
5. **Image Dimensions**: Test extreme aspect ratios
6. **Platform Testing**: Test on macOS, Linux, Windows

## Security Considerations

1. **Input Validation**: User-provided title and URL are used directly in rendering
   - Consider sanitizing for very long strings
   - Current implementation truncates naturally at image boundaries
2. **File Output**: Outputs to user-specified path
   - Could overwrite existing files
   - No permission checks implemented
3. **Command Injection**: Flag parsing is safe (uses `flag` package, not shell)

## References

- [Mat Ryer's Blog Post](https://pace.dev/blog/2020/03/02/dynamically-generate-social-images-in-golang-by-mat-ryer.html)
- [fogleman/gg Documentation](https://github.com/fogleman/gg)
- [Go image Package](https://golang.org/pkg/image/)
- [Twitter Card Specification](https://developer.twitter.com/en/docs/twitter-for-websites/cards/overview/abouts-cards)
- [Open Graph Protocol](https://ogp.me/)
