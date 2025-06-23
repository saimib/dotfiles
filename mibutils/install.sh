#!/bin/bash

# Build the application
echo "Building mibutils..."
go build -o mibutils .

# Check if build was successful
if [ $? -ne 0 ]; then
    echo "Build failed!"
    exit 1
fi

# Create bin directory in home if it doesn't exist
mkdir -p "$HOME/bin"

# Copy the binary to ~/bin
cp mibutils "$HOME/bin/"

# Check if ~/bin is in PATH
if [[ ":$PATH:" != *":$HOME/bin:"* ]]; then
    echo ""
    echo "WARNING: $HOME/bin is not in your PATH."
    echo "Add the following line to your shell profile (.bashrc, .zshrc, etc.):"
    echo "export PATH=\"\$HOME/bin:\$PATH\""
    echo ""
fi

echo "mibutils installed successfully to $HOME/bin/mibutils"
echo ""
echo "Usage examples:"
echo "  mibutils pdf overlay --file1 base.pdf --file2 overlay.pdf --output result.pdf"
echo ""
echo "Run 'mibutils --help' for more information."
