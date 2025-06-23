# mibutils

A collection of utility commands for various tasks including PDF manipulation.

## Installation

1. Make sure you have Go installed (version 1.24.4 or later)
2. Clone or download this repository
3. Run the installation script:
   ```bash
   ./install.sh
   ```

This will build the application and install it to `~/bin/mibutils`.

## Commands

### PDF Overlay

Overlay pages from two PDF files. Pages from the second file will be overlaid onto pages from the first file.

**Usage:**
```bash
mibutils pdf overlay --file1 base.pdf --file2 overlay.pdf --output result.pdf
```

**Parameters:**
- `--file1`: Path to the first PDF file (base layer)
- `--file2`: Path to the second PDF file (overlay layer)  
- `--output`: Path for the output PDF file

**Behavior:**
- Pages from file2 are overlaid on top of pages from file1
- If PDFs have different page counts, remaining pages from the longer PDF are included as-is
- Original page sizes are preserved without any scaling or modification

**Examples:**
```bash
# Basic overlay
mibutils pdf overlay --file1 document.pdf --file2 watermark.pdf --output watermarked.pdf

# Overlay with different page counts
mibutils pdf overlay --file1 5-page-doc.pdf --file2 2-page-overlay.pdf --output result.pdf
# Result: Pages 1-2 overlaid, pages 3-5 from first file added as-is
```

## Development

### Building
```bash
go build -o mibutils .
```

### Adding Dependencies
```bash
go mod tidy
```

### Project Structure
```
mibutils/
├── main.go           # Entry point
├── cmd/
│   ├── root.go       # Root command definition
│   └── pdf.go        # PDF manipulation commands
├── go.mod            # Go module definition
├── install.sh        # Installation script
└── README.md         # This file
```

## Dependencies

- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [pdfcpu](https://github.com/pdfcpu/pdfcpu) - Open-source PDF processing library

## License

See LICENSE file for details.
