package policy_export

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
)

const (
	policyExportEndpoint = "/zia/api/v1/exportPolicies"
)

// ExportPoliciesToJSON sends the policyTypes to the export API, unzips the returned data,
// and writes each JSON file to the specified outputDir (e.g. "./exported_policies").
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
		// Some files might be directories or non-JSON; typically they are "foo.json"
		if f.FileInfo().IsDir() {
			continue
		}

		// 5) Open each file in the ZIP
		zippedFile, err := f.Open()
		if err != nil {
			return fmt.Errorf("failed to open zipped file %q: %w", f.Name, err)
		}
		defer zippedFile.Close()

		// 6) Create a destination file on disk
		destPath := filepath.Join(outputDir, f.Name)
		outFile, err := os.Create(destPath)
		if err != nil {
			return fmt.Errorf("failed to create file %q: %w", destPath, err)
		}

		// 7) Copy the file contents
		if _, err := io.Copy(outFile, zippedFile); err != nil {
			outFile.Close()
			return fmt.Errorf("failed to write file %q: %w", destPath, err)
		}
		outFile.Close()
		service.Client.GetLogger().Printf("[INFO] Extracted %s", destPath)
	}

	return nil
}
