package report

import "fmt"

const (
	reportQuotaEndpoint = "/report/quota"
	reportMd5Endpoint   = "/report"
)

type SandBoxReport struct {
	StartTime int    `json:"startTime"`
	Used      int    `json:"used,omitempty"`
	Allowed   int    `json:"allowed,omitempty"`
	Scale     string `json:"scale,omitempty"`
	Unused    int    `json:"unused,omitempty"`
}

func (service *Service) GetQuota() (*SandBoxReport, error) {
	var quota SandBoxReport
	err := service.Client.Read(reportQuotaEndpoint, &quota)
	if err != nil {
		return nil, err
	}

	return &quota, nil
}

func (service *Service) GetMD5Hash(md5Hash string) (*SandBoxReport, error) {
	var md5 SandBoxReport
	err := service.Client.Read(fmt.Sprintf("%s/%s", reportMd5Endpoint, md5Hash), &md5)
	if err != nil {
		return nil, err
	}

	service.Client.Logger.Printf("[DEBUG]Returning firewall rule from Get: %d", md5)
	return &md5, nil
}
