#!/usr/bin/env nu
# Example script to download emoji assets and run the emoji composite example
#
# Usage: nu run_example.nu

def main [] {
    print "=== Emoji Composite Example ==="
    print ""

    # Create emoji directory
    mkdir emoji

    # Emoji map: codepoint -> name
    let emoji_map = [
        { codepoint: "1f3c8", name: "football 🏈" }
        { codepoint: "26bd", name: "soccer ⚽" }
        { codepoint: "1f3c0", name: "basketball 🏀" }
        { codepoint: "26be", name: "baseball ⚾" }
        { codepoint: "1f3be", name: "tennis 🎾" }
        { codepoint: "2b50", name: "star ⭐" }
        { codepoint: "1f525", name: "fire 🔥" }
        { codepoint: "2764", name: "heart ❤️" }
        { codepoint: "1f44d", name: "thumbsup 👍" }
        { codepoint: "1f680", name: "rocket 🚀" }
    ]

    # Twemoji base URL (72x72 resolution)
    let twemoji_base = "https://raw.githubusercontent.com/twitter/twemoji/master/assets/72x72"

    print "Downloading emoji assets from Twemoji..."
    print ""

    for emoji in $emoji_map {
        let output = $"emoji/($emoji.codepoint).png"
        
        if ($output | path exists) {
            print $"  ✓ ($emoji.name) \(already exists\)"
        } else {
            let url = $"($twemoji_base)/($emoji.codepoint).png"
            try {
                http get $url | save $output
                print $"  ✓ ($emoji.name)"
            } catch {
                print $"  ✗ ($emoji.name) \(download failed\)"
            }
        }
    }

    print ""
    print "Running Go example..."
    print ""

    go run main.go

    print ""
    print "Done! Check output.png"
}

main
