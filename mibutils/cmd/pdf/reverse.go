package pdf

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/spf13/cobra"
)

var reverseCmd = &cobra.Command{
	Use:   "reverse",
	Short: "reverse order of pages",
	Long:  `order of pages in a pdf will be reversed and saved as a new pdf`,
	RunE:  runReverse,
}

func init() {
	reverseCmd.Flags().StringVar(&file1Path, "file", "", "Path to the PDF file")
	reverseCmd.Flags().StringVar(&outputPath, "output", "", "Path for the output PDF file")

	reverseCmd.MarkFlagRequired("file")
	reverseCmd.MarkFlagRequired("output")
}

func runReverse(cmd *cobra.Command, args []string) error {
	// Validate input file exists
	if err := validateFile(file1Path); err != nil {
		return fmt.Errorf("file1 error: %w", err)
	}

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	fmt.Printf("Loading PDF file...\n")

	// Get page count
	pageCount, err := getPageCount(file1Path)
	if err != nil {
		return fmt.Errorf("failed to get page count for %s: %w", file1Path, err)
	}

	fmt.Printf("PDF has %d pages\n", pageCount)

	// Create temporary directory for processing
	tempDir, err := os.MkdirTemp("", "pdf_reverse_*")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// Split PDF into individual pages
	fmt.Printf("Splitting PDF into individual pages...\n")

	pagesDir := filepath.Join(tempDir, "pages")
	if err := os.MkdirAll(pagesDir, 0755); err != nil {
		return fmt.Errorf("failed to create pages directory: %w", err)
	}

	// Split file into individual pages
	if err := api.SplitFile(file1Path, pagesDir, 1, nil); err != nil {
		return fmt.Errorf("failed to split file: %w", err)
	}

	// Get actual split file names by listing the directory
	pageFiles, err := getSplitFiles(pagesDir)
	if err != nil {
		return fmt.Errorf("failed to get split files: %w", err)
	}

	fmt.Printf("Found %d page files\n", len(pageFiles))

	// Reverse the order of pages
	fmt.Printf("Reversing page order...\n")
	var reversedPages []string
	for i := len(pageFiles) - 1; i >= 0; i-- {
		reversedPages = append(reversedPages, pageFiles[i])
	}

	// Merge pages in reversed order into final output
	fmt.Printf("Merging pages into final PDF...\n")
	if err := api.MergeCreateFile(reversedPages, outputPath, false, nil); err != nil {
		return fmt.Errorf("failed to merge final PDF: %w", err)
	}

	fmt.Printf("Successfully created reversed PDF with %d pages\n", len(reversedPages))
	return nil
}
