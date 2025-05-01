package ftp_control_policy

import (
	"context"
	"net/http"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
)

const (
	FTPSettingsEndpoint = "/zia/api/v1/ftpSettings"
)

type FTPControlPolicy struct {
	// Indicates whether to enable FTP over HTTP.
	FtpOverHttpEnabled bool `json:"ftpOverHttpEnabled,omitempty"`

	// Indicates whether to enable native FTP. When enabled, users can connect to native FTP sites and download files.
	FtpEnabled bool `json:"ftpEnabled,omitempty"`

	// List of URL categories that allow FTP traffic
	UrlCategories []string `json:"urlCategories,omitempty"`

	// Domains or URLs included for the FTP Control settings
	Urls []string `json:"urls,omitempty"`
}

func GetFTPControlPolicy(ctx context.Context, service *zscaler.Service) (*FTPControlPolicy, error) {
	var advSettings FTPControlPolicy
	err := service.Client.Read(ctx, FTPSettingsEndpoint, &advSettings)
	if err != nil {
		return nil, err
	}
	return &advSettings, nil
}

func UpdateFTPControlPolicy(ctx context.Context, service *zscaler.Service, ftpSettings *FTPControlPolicy) (*FTPControlPolicy, *http.Response, error) {
	resp, err := service.Client.UpdateWithPut(ctx, (FTPSettingsEndpoint), *ftpSettings)
	if err != nil {
		return nil, nil, err
	}
	updatedSettings, _ := resp.(*FTPControlPolicy)

	service.Client.GetLogger().Printf("[DEBUG]returning updates ftp control policy from update: %d", updatedSettings)
	return updatedSettings, nil, nil
}
