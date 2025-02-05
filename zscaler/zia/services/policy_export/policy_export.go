package policy_export

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
)

const (
	policyExportEndpoint = "/zia/api/v1/exportPolicies"
)

// ExportPolicies sends the policyTypes to the export API, unzips the returned data,
// and writes each JSON file to outputDir (e.g., "./exported_policies").
//
// This version includes a fix for the "zip slip" vulnerability:
//  1. Cleans and checks paths to prevent extraction outside outputDir.
//  2. Creates subdirectories, if any, within the archive.
func ExportPolicies(ctx context.Context, service *zscaler.Service, policyTypes []string, outputDir string) error {
	// 1) Call CreateWithSlicePayload to POST the slice of policyTypes and get the ZIP data
	respBody, err := service.Client.CreateWithSlicePayload(ctx, policyExportEndpoint, policyTypes)
	if err != nil {
		return fmt.Errorf("failed to export policies: %w", err)
	}

	service.Client.GetLogger().Printf("[INFO] Successfully triggered export. Unzipping response...")

	// 2) Create output directory if it doesn't exist
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// 3) Convert respBody bytes into a zip reader
	r, err := zip.NewReader(bytes.NewReader(respBody), int64(len(respBody)))
	if err != nil {
		return fmt.Errorf("failed to read ZIP from response: %w", err)
	}

	// 4) Iterate through each file in the ZIP
	for _, f := range r.File {
		// If it's a directory, just create it (if preserving subfolders)
		if f.FileInfo().IsDir() {
			// We'll sanitize the folder path too
			folderName := filepath.Clean(f.Name)
			destDir := filepath.Join(outputDir, folderName)

			// Check for zip slip by ensuring the result is still inside outputDir
			if !strings.HasPrefix(destDir, filepath.Clean(outputDir)+string(os.PathSeparator)) {
				return fmt.Errorf("zip slip attempted with directory %q", f.Name)
			}

			if err := os.MkdirAll(destDir, 0o755); err != nil {
				return fmt.Errorf("failed to create directory %q: %w", destDir, err)
			}
			continue
		}

		// 5) Clean the file name to remove any ../ or other path tricks
		cleanedName := filepath.Clean(f.Name)

		// 6) Combine it with outputDir to get the final path
		destPath := filepath.Join(outputDir, cleanedName)

		// 7) Check for zip slip
		if !strings.HasPrefix(destPath, filepath.Clean(outputDir)+string(os.PathSeparator)) {
			return fmt.Errorf("zip slip attempted with file %q", f.Name)
		}

		// 8) Open each file in the ZIP
		zippedFile, err := f.Open()
		if err != nil {
			return fmt.Errorf("failed to open zipped file %q: %w", f.Name, err)
		}
		defer zippedFile.Close()

		// Create subdirectories if needed (in case the file is inside a subfolder)
		if err := os.MkdirAll(filepath.Dir(destPath), 0o755); err != nil {
			return fmt.Errorf("failed to create subdirectory for %q: %w", destPath, err)
		}

		// 9) Create a destination file on disk
		outFile, err := os.Create(destPath)
		if err != nil {
			return fmt.Errorf("failed to create file %q: %w", destPath, err)
		}

		// 10) Copy the file contents
		if _, err := io.Copy(outFile, zippedFile); err != nil {
			outFile.Close()
			return fmt.Errorf("failed to write file %q: %w", destPath, err)
		}
		outFile.Close()

		service.Client.GetLogger().Printf("[INFO] Extracted %s", destPath)
	}

	return nil
}
