#!/bin/bash

# Name of the zip file
ZIP_NAME="client.zip"

# Files to include
FILES=("client.js" "index.html" "README.md" "style.css" "websocket.js")

# Create the zip
zip -r "$ZIP_NAME" "${FILES[@]}"

echo "Created $ZIP_NAME with specified files."
