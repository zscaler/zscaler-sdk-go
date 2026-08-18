/*
Copyright (c) 2023, Zscaler Inc.

Permission to use, copy, modify, and/or distribute this software for any
purpose with or without fee is hereby granted, provided that the above
copyright notice and this permission notice appear in all copies.

THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
*/

package browser_access_groups

import (
	"context"
	"fmt"
	"net/http"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
)

const (
	mgmtConfig                  = "/zpa/mgmtconfig/v1/admin/customers/"
	browserAccessGroupsEndpoint = "/browserAccessGroups"
)

type BrowserAccessGroups struct {
	ID                            string            `json:"id,omitempty"`
	Name                          string            `json:"name,omitempty"`
	Description                   string            `json:"description,omitempty"`
	Enabled                       bool              `json:"enabled"`
	City                          string            `json:"city,omitempty"`
	CityCountry                   string            `json:"cityCountry,omitempty"`
	CountryCode                   string            `json:"countryCode,omitempty"`
	CreationTime                  string            `json:"creationTime,omitempty"`
	GeoLocationID                 string            `json:"geoLocationId,omitempty"`
	Latitude                      string            `json:"latitude,omitempty"`
	Location                      string            `json:"location,omitempty"`
	Longitude                     string            `json:"longitude,omitempty"`
	ModifiedBy                    string            `json:"modifiedBy,omitempty"`
	ModifiedTime                  string            `json:"modifiedTime,omitempty"`
	OverrideVersionProfile        bool              `json:"overrideVersionProfile"`
	PrivateExporters              []PrivateExporter `json:"privateExporters,omitempty"`
	ReadOnly                      bool              `json:"readOnly,omitempty"`
	RestrictionType               string            `json:"restrictionType,omitempty"`
	MicrotenantID                 string            `json:"microtenantId,omitempty"`
	MicrotenantName               string            `json:"microtenantName,omitempty"`
	SelectedUpgradePriority       string            `json:"selectedUpgradePriority,omitempty"`
	EnrollmentCertID              string            `json:"enrollmentCertId,omitempty"`
	PrivateCloudID                string            `json:"privateCloudId,omitempty"`
	SiteName                      string            `json:"siteName,omitempty"`
	UpgradeDay                    string            `json:"upgradeDay,omitempty"`
	UpgradePriorities             []string          `json:"upgradePriorities,omitempty"`
	UpgradePriority               string            `json:"upgradePriority,omitempty"`
	UpgradeTimeInSecs             string            `json:"upgradeTimeInSecs,omitempty"`
	Version                       Version           `json:"version,omitempty"`
	VersionProfileID              string            `json:"versionProfileId,omitempty"`
	VersionProfileName            string            `json:"versionProfileName,omitempty"`
	VersionProfileVisibilityScope string            `json:"versionProfileVisibilityScope,omitempty"`
	ZscalerManaged                bool              `json:"zscalerManaged,omitempty"`
}

type PrivateExporter struct {
	ApplicationStartTime             string                       `json:"applicationStartTime,omitempty"`
	ControlChannelStatus             string                       `json:"controlChannelStatus,omitempty"`
	CreationTime                     string                       `json:"creationTime,omitempty"`
	CtrlBrokerName                   string                       `json:"ctrlBrokerName,omitempty"`
	CurrentVersion                   string                       `json:"currentVersion,omitempty"`
	Description                      string                       `json:"description,omitempty"`
	Enabled                          bool                         `json:"enabled,omitempty"`
	EnrollmentTime                   string                       `json:"enrollmentTime,omitempty"`
	ExpectedSargeVersion             string                       `json:"expectedSargeVersion,omitempty"`
	ExpectedUpgradeTime              string                       `json:"expectedUpgradeTime,omitempty"`
	ExpectedVersion                  string                       `json:"expectedVersion,omitempty"`
	Fingerprint                      string                       `json:"fingerprint,omitempty"`
	ID                               string                       `json:"id,omitempty"`
	IpAcl                            []string                     `json:"ipAcl,omitempty"`
	IssuedCertID                     string                       `json:"issuedCertId,omitempty"`
	LastBrokerConnectTime            string                       `json:"lastBrokerConnectTime,omitempty"`
	LastBrokerConnectTimeDuration    string                       `json:"lastBrokerConnectTimeDuration,omitempty"`
	LastBrokerDisconnectTime         string                       `json:"lastBrokerDisconnectTime,omitempty"`
	LastBrokerDisconnectTimeDuration string                       `json:"lastBrokerDisconnectTimeDuration,omitempty"`
	LastOsUpgradeTime                string                       `json:"lastOSUpgradeTime,omitempty"`
	LastSargeUpgradeTime             string                       `json:"lastSargeUpgradeTime,omitempty"`
	LastUpgradeTime                  string                       `json:"lastUpgradeTime,omitempty"`
	Latitude                         string                       `json:"latitude,omitempty"`
	ListenIps                        []string                     `json:"listenIps,omitempty"`
	Location                         string                       `json:"location,omitempty"`
	Longitude                        string                       `json:"longitude,omitempty"`
	ModifiedBy                       string                       `json:"modifiedBy,omitempty"`
	ModifiedTime                     string                       `json:"modifiedTime,omitempty"`
	Name                             string                       `json:"name,omitempty"`
	ProvisioningKeyID                string                       `json:"provisioningKeyId,omitempty"`
	ProvisioningKeyName              string                       `json:"provisioningKeyName,omitempty"`
	OsUpgradeEnabled                 bool                         `json:"osUpgradeEnabled,omitempty"`
	OsUpgradeStatus                  string                       `json:"osUpgradeStatus,omitempty"`
	Platform                         string                       `json:"platform,omitempty"`
	PlatformDetail                   string                       `json:"platformDetail,omitempty"`
	PlatformVersion                  string                       `json:"platformVersion,omitempty"`
	PreviousVersion                  string                       `json:"previousVersion,omitempty"`
	BrowserAccessGroupsID            string                       `json:"browserAccessGroupsId,omitempty"`
	PrivateExporterGroupName         string                       `json:"privateExporterGroupName,omitempty"`
	PrivateExporterVersion           PrivateExporterVersion       `json:"privateExporterVersion,omitempty"`
	PrivateIp                        string                       `json:"privateIp,omitempty"`
	PublicIp                         string                       `json:"publicIp,omitempty"`
	PublishIps                       []string                     `json:"publishIps,omitempty"`
	ReadOnly                         bool                         `json:"readOnly,omitempty"`
	RestrictionType                  string                       `json:"restrictionType,omitempty"`
	RuntimeOS                        string                       `json:"runtimeOS,omitempty"`
	SargeUpgradeAttempt              string                       `json:"sargeUpgradeAttempt,omitempty"`
	SargeUpgradeStatus               string                       `json:"sargeUpgradeStatus,omitempty"`
	SargeVersion                     string                       `json:"sargeVersion,omitempty"`
	MicrotenantID                    string                       `json:"microtenantId,omitempty"`
	MicrotenantName                  string                       `json:"microtenantName,omitempty"`
	EnrollmentCert                   map[string]interface{}       `json:"enrollmentCert,omitempty"`
	UpgradeAttempt                   string                       `json:"upgradeAttempt,omitempty"`
	UpgradeStatus                    string                       `json:"upgradeStatus,omitempty"`
	Version                          Version                      `json:"version,omitempty"`
	ZpnSubModuleUpgradeList          []common.ZPNSubModuleUpgrade `json:"zpnSubModuleUpgradeList,omitempty"`
	ZscalerManaged                   bool                         `json:"zscalerManaged,omitempty"`
}

type PrivateExporterVersion struct {
	ApplicationStartTime  string                       `json:"applicationStartTime,omitempty"`
	BrokerID              string                       `json:"brokerId,omitempty"`
	CreationTime          string                       `json:"creationTime,omitempty"`
	CtrlChannelStatus     string                       `json:"ctrlChannelStatus,omitempty"`
	CurrentVersion        string                       `json:"currentVersion,omitempty"`
	DisableAutoUpdate     bool                         `json:"disableAutoUpdate,omitempty"`
	ExpectedSargeVersion  string                       `json:"expectedSargeVersion,omitempty"`
	ExpectedVersion       string                       `json:"expectedVersion,omitempty"`
	ID                    string                       `json:"id,omitempty"`
	LastConnectTime       string                       `json:"lastConnectTime,omitempty"`
	LastDisconnectTime    string                       `json:"lastDisconnectTime,omitempty"`
	LastOsUpgradeTime     string                       `json:"lastOSUpgradeTime,omitempty"`
	LastSargeUpgradeTime  string                       `json:"lastSargeUpgradeTime,omitempty"`
	LastUpgradedTime      string                       `json:"lastUpgradedTime,omitempty"`
	LoneWarrior           bool                         `json:"loneWarrior,omitempty"`
	ModifiedBy            string                       `json:"modifiedBy,omitempty"`
	ModifiedTime          string                       `json:"modifiedTime,omitempty"`
	OsUpgradeEnabled      bool                         `json:"osUpgradeEnabled,omitempty"`
	OsUpgradeStatus       string                       `json:"osUpgradeStatus,omitempty"`
	Platform              string                       `json:"platform,omitempty"`
	PlatformDetail        string                       `json:"platformDetail,omitempty"`
	PlatformVersion       string                       `json:"platformVersion,omitempty"`
	PreviousVersion       string                       `json:"previousVersion,omitempty"`
	BrowserAccessGroupsID string                       `json:"browserAccessGroupsId,omitempty"`
	PrivateIp             string                       `json:"privateIp,omitempty"`
	PublicIp              string                       `json:"publicIp,omitempty"`
	RestartInstructions   string                       `json:"restartInstructions,omitempty"`
	RestartTimeInSec      string                       `json:"restartTimeInSec,omitempty"`
	RuntimeOS             string                       `json:"runtimeOS,omitempty"`
	SargeUpgradeAttempt   string                       `json:"sargeUpgradeAttempt,omitempty"`
	SargeUpgradeStatus    string                       `json:"sargeUpgradeStatus,omitempty"`
	SargeVersion          string                       `json:"sargeVersion,omitempty"`
	SystemStartTime       string                       `json:"systemStartTime,omitempty"`
	TunnelID              string                       `json:"tunnelId,omitempty"`
	UpgradeAttempt        string                       `json:"upgradeAttempt,omitempty"`
	UpgradeNowOnce        bool                         `json:"upgradeNowOnce,omitempty"`
	UpgradeStatus         string                       `json:"upgradeStatus,omitempty"`
	ZpnSubModuleUpgrade   []common.ZPNSubModuleUpgrade `json:"zpnSubModuleUpgrade,omitempty"`
}

type Version struct {
	ChildVersion       string `json:"childVersion,omitempty"`
	LatestPlatform     string `json:"latestPlatform,omitempty"`
	Platform           string `json:"platform,omitempty"`
	SargeVersion       string `json:"sargeVersion,omitempty"`
	VersionProfileName string `json:"versionProfileName,omitempty"`
	VersionProfileGid  string `json:"version_profile_gid,omitempty"`
}

func GetAll(ctx context.Context, service *zscaler.Service) ([]BrowserAccessGroups, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + browserAccessGroupsEndpoint
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[BrowserAccessGroups](ctx, service.Client, relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()})
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}

func Get(ctx context.Context, service *zscaler.Service, browserAccessGroupsID string) (*BrowserAccessGroups, *http.Response, error) {
	v := new(BrowserAccessGroups)
	resp, err := service.Client.NewRequestDo(ctx, "GET", fmt.Sprintf("%v/%v", mgmtConfig+service.Client.GetCustomerID()+browserAccessGroupsEndpoint, browserAccessGroupsID), common.Filter{MicroTenantID: service.MicroTenantID()}, nil, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}
