# Development Journal

## 2026-01-17: Dynamic URL Baseline Positioning

### Problem

The `drawURL` function hardcoded the URL's Y position using a multiplier of 5 times the title baseline value:

```go
targetY := TextTopMargin + titleFontHeight + (titleFontHeight*LineSpacing)*5
```

This caused the URL to appear in incorrect positions for different image heights or when using fonts with varying x-heights.

### Solution

Updated `drawURL` to dynamically calculate the URL baseline position based on the typographic grid established by the title font:

1. **Load title font metrics** - Get the title font height to establish the baseline grid
2. **Calculate baseline grid** - First baseline at `TextTopMargin + titleFontHeight`, with subsequent baselines at intervals of `titleFontHeight * LineSpacing`
3. **Find last valid baseline** - Iterate through baselines to find the last one that fits within `height - BackgroundMargin`
4. **Place URL on grid** - Position the URL on this calculated baseline

```go
firstBaseline := TextTopMargin + titleFontHeight
baselineStep := titleFontHeight * LineSpacing
maxY := float64(height) - BackgroundMargin

targetY := firstBaseline
for y := firstBaseline; y <= maxY; y += baselineStep {
    targetY = y
}
```

### Key Insight

The baseline grid must be calculated from the **title font**, not the URL font. This maintains consistent vertical rhythm across the entire image regardless of which fonts are selected or their respective metrics.

### Tests Added

- `TestDrawURLPositionDynamic` - Verifies URL positioning works correctly for various image heights (400, 628, 1080 pixels)
- Tests confirm the URL sits on the baseline grid and remains within image bounds
