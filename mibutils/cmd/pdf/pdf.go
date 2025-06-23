/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package pdf

import (
	"github.com/spf13/cobra"
)

var (
	file1Path  string
	file2Path  string
	outputPath string
)

// pdfCmd represents the pdf command
var PDFCmd = &cobra.Command{
	Use:   "pdf",
	Short: "PDF manipulation utilities",
	Long:  `PDF manipulation utilities including overlay, split, and merge operations.`,
}

func init() {
	PDFCmd.AddCommand(overlayCmd)
	PDFCmd.AddCommand(reverseCmd)
	PDFCmd.AddCommand(compressCmd)
}
