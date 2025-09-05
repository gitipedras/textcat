#!/bin/bash
# bundle.sh - Combine HTML, CSS, and JS into one file with a custom title

SRC_DIR="src"
BUILD_DIR="build"
OUTPUT="$BUILD_DIR/index.html"

# Set the title
TITLEVAR="Textcat Web Bundle"

# Make sure build directory exists
mkdir -p "$BUILD_DIR"

# Start building the HTML
awk -v myTitle="$TITLEVAR" -v cssfile="$SRC_DIR/style.css" '
/<head>/ { 
    print
    print "  <title>" myTitle "</title>"
    print "  <style>"
    while ((getline cssline < cssfile) > 0) print cssline
    print "  </style>"
    next
}
# Replace <h1>Login</h1> with actual title
{
    gsub(/<h1>Login<\/h1>/, "<h1>" myTitle "</h1>")
    print
}
# Skip existing <script> blocks
/<script>/,/<\/script>/ { next }
' "$SRC_DIR/index.html" > "$OUTPUT"

# Insert JS just before </body>
awk -v js1="$SRC_DIR/websocket.js" -v js2="$SRC_DIR/client.js" -v js3="$SRC_DIR/clientconfig.js" '
/<\/body>/ {
    print "<script>"
    while ((getline line < js1) > 0) print line
    print "</script>"
    print "<script>"
    while ((getline line < js2) > 0) print line
    print "</script>"
    print "<script>"
    while ((getline line < js3) > 0) print line
    print "</script>"
    print
    next
}
{ print }
' "$OUTPUT" > "$OUTPUT.tmp" && mv "$OUTPUT.tmp" "$OUTPUT"


echo "Bundled file created: $OUTPUT with title: $TITLEVAR"
