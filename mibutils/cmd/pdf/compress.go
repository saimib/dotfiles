package pdf

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/spf13/cobra"
)

var compressCmd = &cobra.Command{
	Use:   "compress",
	Short: "Compress PDF file to reduce size",
	Long:  `Compress a PDF file to reduce its size without significantly affecting quality. Uses optimization techniques to remove redundant data and compress images.`,
	RunE:  runCompress,
}

func init() {
	compressCmd.Flags().StringVar(&file1Path, "file", "", "Path to the PDF file to compress")
	compressCmd.Flags().StringVar(&outputPath, "output", "", "Path for the compressed output PDF file")

	compressCmd.MarkFlagRequired("file")
	compressCmd.MarkFlagRequired("output")
}

func runCompress(cmd *cobra.Command, args []string) error {
	// Validate input file exists
	if err := validateFile(file1Path); err != nil {
		return fmt.Errorf("input file error: %w", err)
	}

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	fmt.Printf("Loading PDF file: %s\n", file1Path)

	// Get original file size
	originalInfo, err := os.Stat(file1Path)
	if err != nil {
		return fmt.Errorf("failed to get original file info: %w", err)
	}
	originalSize := originalInfo.Size()

	// Get page count for information
	pageCount, err := getPageCount(file1Path)
	if err != nil {
		return fmt.Errorf("failed to get page count for %s: %w", file1Path, err)
	}

	fmt.Printf("Original PDF: %d pages, %.2f MB\n", pageCount, float64(originalSize)/(1024*1024))

	// Try multiple compression strategies
	bestSize := originalSize
	bestFile := file1Path

	// Create temporary directory for compression attempts
	tempDir, err := os.MkdirTemp("", "pdf_compress_*")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// Strategy 1: Basic optimization
	fmt.Printf("Trying basic optimization...\n")
	tempFile1 := filepath.Join(tempDir, "optimized.pdf")
	if err := api.OptimizeFile(file1Path, tempFile1, nil); err == nil {
		if info, err := os.Stat(tempFile1); err == nil && info.Size() < bestSize {
			bestSize = info.Size()
			bestFile = tempFile1
			fmt.Printf("Basic optimization: %.2f MB (%.1f%% reduction)\n",
				float64(bestSize)/(1024*1024),
				float64(originalSize-bestSize)/float64(originalSize)*100)
		}
	}

	// Strategy 2: Split and merge approach
	fmt.Printf("Trying split-merge compression...\n")
	pagesDir := filepath.Join(tempDir, "pages")
	if err := os.MkdirAll(pagesDir, 0755); err == nil {
		if err := api.SplitFile(bestFile, pagesDir, 1, nil); err == nil {
			if pageFiles, err := getSplitFiles(pagesDir); err == nil {
				tempFile2 := filepath.Join(tempDir, "merged.pdf")
				if err := api.MergeCreateFile(pageFiles, tempFile2, false, nil); err == nil {
					if info, err := os.Stat(tempFile2); err == nil && info.Size() < bestSize {
						bestSize = info.Size()
						bestFile = tempFile2
						fmt.Printf("Split-merge compression: %.2f MB (%.1f%% reduction)\n",
							float64(bestSize)/(1024*1024),
							float64(originalSize-bestSize)/float64(originalSize)*100)
					}
				}
			}
		}
	}

	// Strategy 3: Multiple optimization passes
	fmt.Printf("Trying multiple optimization passes...\n")
	currentFile := bestFile
	for i := 0; i < 3; i++ {
		tempFileN := filepath.Join(tempDir, fmt.Sprintf("pass%d.pdf", i+1))
		if err := api.OptimizeFile(currentFile, tempFileN, nil); err == nil {
			if info, err := os.Stat(tempFileN); err == nil {
				if info.Size() < bestSize {
					bestSize = info.Size()
					bestFile = tempFileN
					fmt.Printf("Pass %d: %.2f MB (%.1f%% reduction)\n",
						i+1,
						float64(bestSize)/(1024*1024),
						float64(originalSize-bestSize)/float64(originalSize)*100)
				} else if info.Size() >= bestSize {
					// No improvement, stop trying
					break
				}
				currentFile = tempFileN
			}
		}
	}

	// Copy the best result to output
	fmt.Printf("Finalizing compression...\n")
	if bestFile == file1Path {
		// No compression achieved, just copy the original
		if err := copyFile(file1Path, outputPath); err != nil {
			return fmt.Errorf("failed to copy file: %w", err)
		}
	} else {
		// Copy the best compressed version
		if err := copyFile(bestFile, outputPath); err != nil {
			return fmt.Errorf("failed to copy compressed file: %w", err)
		}
	}

	// Calculate final compression ratio
	compressionRatio := float64(originalSize-bestSize) / float64(originalSize) * 100

	fmt.Printf("Compression completed!\n")
	fmt.Printf("Original size: %.2f MB\n", float64(originalSize)/(1024*1024))
	fmt.Printf("Compressed size: %.2f MB\n", float64(bestSize)/(1024*1024))
	if compressionRatio > 0 {
		fmt.Printf("Size reduction: %.1f%%\n", compressionRatio)
	} else {
		fmt.Printf("No size reduction achieved - PDF may already be optimized\n")
	}
	fmt.Printf("Output saved to: %s\n", outputPath)

	return nil
}

// Helper function to copy files
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = destFile.ReadFrom(sourceFile)
	return err
}
