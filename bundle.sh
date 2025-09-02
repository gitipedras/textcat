#!/bin/bash
# bundle.sh - Combine HTML, CSS, and JS into one file

SRC_DIR="src"
BUILD_DIR="build"
OUTPUT="$BUILD_DIR/index.html"

# Make sure build directory exists
mkdir -p "$BUILD_DIR"

# Start the output file with the HTML head from index.html
# Up to the closing </head>
awk '/<\/head>/{print "<style>"; system("cat '"$SRC_DIR"/style.css'"); print "</style>"; next}1' "$SRC_DIR/index.html" > "$OUTPUT"

# Append JS files at the end of the body
echo "<script>" >> "$OUTPUT"
cat "$SRC_DIR/websocket.js" >> "$OUTPUT"
echo "</script>" >> "$OUTPUT"

echo "<script>" >> "$OUTPUT"
cat "$SRC_DIR/client.js" >> "$OUTPUT"
echo "</script>" >> "$OUTPUT"

# Add the closing </html> if not present
if ! grep -q "</html>" "$OUTPUT"; then
    echo "</html>" >> "$OUTPUT"
fi

echo "Bundled file created: $OUTPUT"
