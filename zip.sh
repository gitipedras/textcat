#!/bin/bash

# Name of the zip file
ZIP_NAME="client.zip"

# Files to include
FILES=("src/client.js" "src/index.html" "src/README.md" "src/style.css" "src/websocket.js")

# Create the zip
zip -r "$ZIP_NAME" "${FILES[@]}"

echo "Created $ZIP_NAME with files."
