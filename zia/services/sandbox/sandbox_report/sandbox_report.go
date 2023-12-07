package sandbox_report

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/common"
)

const (
	reportQuotaEndpoint = "/sandbox/report/quota"
	reportMD5Endpoint   = "/sandbox/report/"
)

type RatingQuota struct {
	StartTime int    `json:"startTime,omitempty"`
	Used      int    `json:"used,omitempty"`
	Allowed   int    `json:"allowed,omitempty"`
	Scale     string `json:"scale,omitempty"`
	Unused    int    `json:"unused,omitempty"`
}

type ReportMD5Hash struct {
	Details *FullDetails `json:"details,omitempty"`
}

type Summary struct {
	Detail         *SummaryDetail  `json:"Summary,omitempty"`
	Classification *Classification `json:"Classification,omitempty"`
	FileProperties *FileProperties `json:"FileProperties,omitempty"`
}

type SummaryDetail struct {
	Status    string `json:"Status,omitempty"`
	Category  string `json:"Category,omitempty"`
	FileType  string `json:"FileType,omitempty"`
	StartTime int    `json:"StartTime,omitempty"`
	Duration  int    `json:"Duration,omitempty"`
}

type Classification struct {
	Type            string `json:"Type,omitempty"`
	Category        string `json:"Category,omitempty"`
	Score           int    `json:"Score,omitempty"`
	DetectedMalware string `json:"DetectedMalware,omitempty"`
}

type FileProperties struct {
	FileType          string `json:"FileType,omitempty"`
	FileSize          int    `json:"FileSize,omitempty"`
	MD5               string `json:"MD5,omitempty"`
	SHA1              string `json:"SHA1,omitempty"`
	SHA256            string `json:"Sha256,omitempty"`
	Issuer            string `json:"Issuer,omitempty"`
	DigitalCerificate string `json:"DigitalCerificate,omitempty"`
	SSDeep            string `json:"SSDeep,omitempty"`
	RootCA            string `json:"RootCA,omitempty"`
}

type FullDetails struct {
	Summary        SummaryDetail         `json:"Summary,omitempty"`
	Classification Classification        `json:"Classification,omitempty"`
	FileProperties FileProperties        `json:"FileProperties,omitempty"`
	Origin         *Origin               `json:"Origin,omitempty"`
	SystemSummary  []SystemSummaryDetail `json:"SystemSummary,omitempty"`
	Spyware        []*common.SandboxRSS  `json:"Spyware,omitempty"`
	Networking     []*common.SandboxRSS  `json:"Networking,omitempty"`
	SecurityBypass []*common.SandboxRSS  `json:"SecurityBypass,omitempty"`
	Exploit        []*common.SandboxRSS  `json:"Exploit,omitempty"`
	Stealth        []*common.SandboxRSS  `json:"Stealth,omitempty"`
	Persistence    []*common.SandboxRSS  `json:"Persistence,omitempty"`
}

type Origin struct {
	Risk     string `json:"Risk,omitempty"`
	Language string `json:"Language,omitempty"`
	Country  string `json:"Country,omitempty"`
}

type SystemSummaryDetail struct {
	Risk             string   `json:"Risk,omitempty"`
	Signature        string   `json:"Signature,omitempty"`
	SignatureSources []string `json:"SignatureSources,omitempty"`
}

func (service *Service) GetRatingQuota() ([]RatingQuota, error) {
	var quotas []RatingQuota
	err := service.Client.Read(reportQuotaEndpoint, &quotas)
	if err != nil {
		return nil, err
	}

	service.Client.Logger.Printf("[DEBUG] Returning quota for retrieving Sandbox Detail Reports from Get: %v", quotas)
	return quotas, nil
}

// GetReportMD5Hash retrieves the sandbox report for a specific MD5 hash with either full or summary details.
func (service *Service) GetReportMD5Hash(md5Hash, details string) (*ReportMD5Hash, error) {
	// Validate the 'details' parameter to ensure it is either "full" or "summary".
	if details != "full" && details != "summary" {
		return nil, fmt.Errorf("details parameter must be 'full' or 'summary'")
	}

	// Construct the endpoint URL with the md5Hash and details query parameters.
	endpoint := fmt.Sprintf("%s%s?details=%s", reportMD5Endpoint, md5Hash, details)

	var resp map[string]interface{}
	err := service.Client.Read(endpoint, &resp)
	if err != nil {
		return nil, err
	}
	var data interface{}
	var report ReportMD5Hash
	if details == "full" {
		data = resp["Full Details"]
	} else {
		data = resp["Summary"]
	}
	if data == nil {
		return nil, errors.New("got empty response")
	}

	if msg, ok := data.(string); ok {
		return nil, errors.New(msg)
	}

	dataBytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(dataBytes, &report.Details); err != nil {
		return nil, err
	}

	service.Client.Logger.Printf("[DEBUG] Returning report for MD5 hash '%s' with details '%s': %+v", md5Hash, details, report)
	return &report, nil
}
