package sandbox_report

import "fmt"

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
	MD5Hash string `json:"md5Hash,omitempty"`
	Details string `json:"details,omitempty"`
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

	var report ReportMD5Hash
	err := service.Client.Read(endpoint, &report)
	if err != nil {
		return nil, err
	}

	service.Client.Logger.Printf("[DEBUG] Returning report for MD5 hash '%s' with details '%s': %+v", md5Hash, details, report)
	return &report, nil
}
