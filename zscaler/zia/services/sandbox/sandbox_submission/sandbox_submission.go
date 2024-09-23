package sandbox_submission

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
	"mime"
	"net/url"
	"path/filepath"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
)

const (
	submitEndpoint = "/zia/api/v1/zscsb/submit"
	discanEndpoint = "/zia/api/v1/zscsb/discan"
)

// Information about the file inspection results
type ScanResult struct {
	Code              int    `json:"code,omitempty"`
	Message           string `json:"message,omitempty"`
	FileType          string `json:"fileType,omitempty"`
	Md5               string `json:"md5,omitempty"`
	SandboxSubmission string `json:"sandboxSubmission,omitempty"`
	VirusName         string `json:"virusName,omitempty"`
	VirusType         string `json:"virusType,omitempty"`
}

// Submit: Submits raw or archive files (e.g., ZIP) to Sandbox for analysis. You can submit up to 100 files per day and it supports all file types that are currently supported by Sandbox. To learn more, see About Sandbox. By default, files are scanned by Zscaler antivirus (AV) and submitted directly to the sandbox in order to obtain a verdict. However, if a verdict already exists for the file, you can use the 'force' parameter to make the sandbox to reanalyze it.
// You must have a Sandbox policy rule configured within the ZIA Admin Portal in order to analyze files that aren't present in the default policy rule. Ensure that you have explicitly added Sandbox policy rules that include the appropriate file types within your request. If not, an 'Unknown' message is shown in the response.
// To learn more, see Configuring the Sandbox Policy and Configuring the Default Sandbox Rule.
// After files are sent for analysis, you must use GET /sandbox/report/{md5Hash} in order to retrieve the verdict. You can get the Sandbox report 10 minutes after a file is sent for analysis.
// SubmitFile submits a file for scanning and returns the result of the scan.
func SubmitFile(service *zscaler.Service, filename string, file io.Reader, force string) (*ScanResult, error) {
	return scanFile(service, filename, file, force, submitEndpoint)
}

// Discan Submits raw or archive files (e.g., ZIP) to the Zscaler service for out-of-band file inspection to generate real-time verdicts for known and unknown files. It leverages capabilities such as Malware Prevention, Advanced Threat Prevention, Sandbox cloud effect, AI/ML-driven file analysis, and integrated third-party threat intelligence feeds to inspect files and classify them as benign or malicious instantaneously.
// All file types that are currently supported by the Malware Protection policy and Advanced Threat Protection policy are supported for inspection, and each file is limited to a size of 400 MB.
// Note: Dynamic file analysis is not included in out-of-band file inspection.
func Discan(service *zscaler.Service, filename string, file io.Reader) (*ScanResult, error) {
	return scanFile(service, filename, file, "", discanEndpoint)
}

func scanFile(service *zscaler.Service, filename string, file io.Reader, force, endpoint string) (*ScanResult, error) {
	// Add the API token and force parameter to the request URL query
	urlParams := url.Values{}
	urlParams.Set("api_token", service.Client.GetSandboxToken())
	if force != "" {
		urlParams.Set("force", force)
	}

	// Determine the Content-Type based on the file extension
	contentType := mime.TypeByExtension(filepath.Ext(filename))
	if contentType == "" {
		// If the content type cannot be determined, set it to a default value
		contentType = "application/octet-stream"
	}

	// Create a buffer to store the gzipped file
	var gzippedFile bytes.Buffer
	gz := gzip.NewWriter(&gzippedFile)

	// Copy the file content to the gzip writer
	_, err := io.Copy(gz, file)
	if err != nil {
		return nil, err
	}
	gz.Close() // Ensure to close the gzip writer to flush the buffer

	// Correct the argument order for ExecuteRequest
	data, _, err := service.Client.ExecuteRequest("POST", endpoint, &gzippedFile, urlParams, contentType) // Ignore the req value
	if err != nil {
		return nil, err
	}

	var result ScanResult
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
