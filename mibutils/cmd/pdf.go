/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/unidoc/unipdf/v3/creator"
	"github.com/unidoc/unipdf/v3/model"
)

var (
	file1Path  string
	file2Path  string
	outputPath string
)

// pdfCmd represents the pdf command
var pdfCmd = &cobra.Command{
	Use:   "pdf",
	Short: "PDF manipulation utilities",
	Long:  `PDF manipulation utilities including overlay, split, and merge operations.`,
}

// overlayCmd represents the overlay command
var overlayCmd = &cobra.Command{
	Use:   "overlay",
	Short: "Overlay two PDF files",
	Long: `Overlay pages from two PDF files. Pages from file2 will be overlaid onto pages from file1.
If the PDFs have different page counts, remaining pages from the longer PDF will be included as-is.`,
	RunE: runOverlay,
}

func init() {
	rootCmd.AddCommand(pdfCmd)
	pdfCmd.AddCommand(overlayCmd)

	overlayCmd.Flags().StringVar(&file1Path, "file1", "", "Path to the first PDF file (base layer)")
	overlayCmd.Flags().StringVar(&file2Path, "file2", "", "Path to the second PDF file (overlay layer)")
	overlayCmd.Flags().StringVar(&outputPath, "output", "", "Path for the output PDF file")

	overlayCmd.MarkFlagRequired("file1")
	overlayCmd.MarkFlagRequired("file2")
	overlayCmd.MarkFlagRequired("output")
}

func runOverlay(cmd *cobra.Command, args []string) error {
	// Validate input files exist
	if err := validateFile(file1Path); err != nil {
		return fmt.Errorf("file1 error: %w", err)
	}
	if err := validateFile(file2Path); err != nil {
		return fmt.Errorf("file2 error: %w", err)
	}

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	fmt.Printf("Loading PDF files...\n")
	
	// Load both PDF files
	pages1, err := loadPDFPages(file1Path)
	if err != nil {
		return fmt.Errorf("failed to load %s: %w", file1Path, err)
	}
	
	pages2, err := loadPDFPages(file2Path)
	if err != nil {
		return fmt.Errorf("failed to load %s: %w", file2Path, err)
	}

	fmt.Printf("File1: %d pages, File2: %d pages\n", len(pages1), len(pages2))

	// Create output PDF
	c := creator.New()
	
	// Process pages
	minPages := min(len(pages1), len(pages2))
	
	// Overlay pages where both files have pages
	for i := 0; i < minPages; i++ {
		fmt.Printf("Processing page %d (overlay)...\n", i+1)
		if err := overlayPages(c, pages1[i], pages2[i]); err != nil {
			return fmt.Errorf("failed to overlay page %d: %w", i+1, err)
		}
	}
	
	// Add remaining pages from the longer file
	if len(pages1) > minPages {
		fmt.Printf("Adding remaining pages from file1...\n")
		for i := minPages; i < len(pages1); i++ {
			fmt.Printf("Processing page %d (file1 only)...\n", i+1)
			if err := addSinglePage(c, pages1[i]); err != nil {
				return fmt.Errorf("failed to add page %d from file1: %w", i+1, err)
			}
		}
	} else if len(pages2) > minPages {
		fmt.Printf("Adding remaining pages from file2...\n")
		for i := minPages; i < len(pages2); i++ {
			fmt.Printf("Processing page %d (file2 only)...\n", i+1)
			if err := addSinglePage(c, pages2[i]); err != nil {
				return fmt.Errorf("failed to add page %d from file2: %w", i+1, err)
			}
		}
	}

	// Write output file
	fmt.Printf("Saving output to %s...\n", outputPath)
	if err := c.WriteToFile(outputPath); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	fmt.Printf("Successfully created overlaid PDF with %d pages\n", len(pages1)+len(pages2)-minPages)
	return nil
}

func validateFile(path string) error {
	if path == "" {
		return fmt.Errorf("file path cannot be empty")
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("file does not exist: %s", path)
	}
	return nil
}

func loadPDFPages(filePath string) ([]*model.PdfPage, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	pdfReader, err := model.NewPdfReader(f)
	if err != nil {
		return nil, err
	}

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return nil, err
	}

	var pages []*model.PdfPage
	for i := 1; i <= numPages; i++ {
		page, err := pdfReader.GetPage(i)
		if err != nil {
			return nil, fmt.Errorf("failed to get page %d: %w", i, err)
		}
		pages = append(pages, page)
	}

	return pages, nil
}

func overlayPages(c *creator.Creator, basePage, overlayPage *model.PdfPage) error {
	// Get page dimensions from base page
	bbox, err := basePage.GetMediaBox()
	if err != nil {
		return err
	}

	// Create a new page with base page dimensions
	c.SetPageSize(creator.PageSize{bbox.Width(), bbox.Height()})
	c.NewPage()

	// Add base page content as a block
	baseBlock, err := creator.NewBlockFromPage(basePage)
	if err != nil {
		return err
	}
	baseBlock.SetPos(0, 0)
	c.Draw(baseBlock)

	// Add overlay page content as a block
	overlayBlock, err := creator.NewBlockFromPage(overlayPage)
	if err != nil {
		return err
	}
	overlayBlock.SetPos(0, 0)
	c.Draw(overlayBlock)

	return nil
}

func addSinglePage(c *creator.Creator, pdfPage *model.PdfPage) error {
	// Get page dimensions
	bbox, err := pdfPage.GetMediaBox()
	if err != nil {
		return err
	}

	// Create a new page
	c.SetPageSize(creator.PageSize{bbox.Width(), bbox.Height()})
	c.NewPage()

	// Add page content as a block
	block, err := creator.NewBlockFromPage(pdfPage)
	if err != nil {
		return err
	}
	block.SetPos(0, 0)
	c.Draw(block)

	return nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
