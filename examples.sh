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
echo "=== Example 1 alt: Basic usage ==="
./og-image-generator \
  -title "Do first, and then understand later" \
  -url "https://example.com/go-apis" \
  -output out/do-first-understand-later.png

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
echo "=== Example 4 alt: Very long title alt (tests text wrapping) ==="
./og-image-generator \
  -title "Advanced Patterns for Building High Performance, Scalable, and Maintainable Web Services with Go" \
  -url "https://example.com/advanced-patterns" \
  -output out/long-title-alt.png \
  -bg "#cccc00"

echo ""
echo "=== Example 4 short: Very long title with short words (tests text wrapping) ==="
./og-image-generator \
  -title "The as via or with can alt vip run task bib bit lip too lid not eql mut var let const" \
  -url "https://example.com/advanced-patterns" \
  -output out/short-words.png \
  -bg "#cc00cc"

echo ""
echo "=== Example 4 debug: Activate debug mode ==="
./og-image-generator \
  -title "The as via or with can alt vip run task bib bit lip too lid not eql mut var let const" \
  -url "https://example.com/advanced-patterns" \
  -output out/debug-baselines.png \
  -debug \
  -bg "#00cccc"

echo ""
echo "=== Example 5: Orphan prevention - title that would end with single word ==="
# Without orphan prevention, "Today" would appear alone on the last line
# With orphan prevention, "Services Today" appears together
./og-image-generator \
  -title "Building High-Performance Web Services Today" \
  -url "https://example.com/web-services" \
  -output out/orphan-prevented.png \
  -bg "#1e3a5f"

echo ""
echo "=== Example 6: Another orphan prevention example ==="
# Without orphan prevention, "Go" would appear alone on the last line
# With orphan prevention, "with Go" appears together
./og-image-generator \
  -title "Modern API Development with Go" \
  -url "https://example.com/api-dev" \
  -output out/orphan-prevented-2.png \
  -bg "#2d3436"

echo ""
echo "=== Example 7: Three-line title with orphan prevention ==="
# Tests orphan prevention when text wraps to three lines
./og-image-generator \
  -title "Understanding Distributed Systems and Building Reliable Microservices Architecture" \
  -url "https://example.com/distributed" \
  -output out/three-line-orphan.png \
  -bg "#6c5ce7"

echo ""
echo "Generated images:"
ls -lh out/*.png
