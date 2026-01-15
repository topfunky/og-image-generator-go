#!/bin/bash

# Example usage script for og-image-generator

# Build the application
echo "Building og-image-generator..."
go build -o og-image-generator

# Create output directory
mkdir -p out

echo ""
echo "=== Example 1: Basic usage ==="
./og-image-generator \
  -title "How to Build APIs in Go" \
  -url "https://example.com/go-apis" \
  -output out/social-image.png

echo ""
echo "=== Example 2: Custom output and background color ==="
./og-image-generator \
  -title "Mastering Concurrency in Go" \
  -url "https://example.com/concurrency" \
  -output out/dark-image.png \
  -bg "#0f0f1e"

echo ""
echo "=== Example 3: Custom dimensions ==="
./og-image-generator \
  -title "Building Production Systems" \
  -url "https://example.com/production" \
  -output out/wide-image.png \
  -width 1600 \
  -height 900

echo ""
echo "=== Example 4: Very long title (tests text wrapping) ==="
./og-image-generator \
  -title "Advanced Patterns for Building High-Performance, Scalable, and Maintainable Web Services with Go" \
  -url "https://example.com/advanced-patterns" \
  -output out/long-title.png \
  -bg "#2c3e50"

echo ""
echo "Generated images:"
ls -lh out/*.png
