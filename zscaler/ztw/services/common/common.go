package common

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw"
)

const pageSize = 1000

type IDNameExtensions struct {
	ID         int                    `json:"id,omitempty"`
	Name       string                 `json:"name,omitempty"`
	Extensions map[string]interface{} `json:"extensions,omitempty"`
}

// General purpose object. This is an immutable reference to an entity, which mainly consists of ID and name.
type CommonIDNameExternalID struct {
	ID              int                    `json:"id,omitempty"`
	Name            string                 `json:"name,omitempty"`
	IsNameL10nTag   bool                   `json:"isNameL10nTag,omitempty"`
	Extensions      map[string]interface{} `json:"extensions,omitempty"`
	Deleted         bool                   `json:"deleted,omitempty"`
	ExternalId      string                 `json:"externalId,omitempty"`
	AssociationTime int                    `json:"associationTime,omitempty"`
}

type CommonZPAIDNameID struct {
	ID          int    `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	ExternalId  string `json:"externalId,omitempty"`
	ZPATenantId int    `json:"zpaTenantId,omitempty"`
}

type CommonIDName struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type ECVMLite struct {
	ID               int    `json:"id,omitempty"`
	Name             string `json:"name,omitempty"`
	UpgradeStatus    int    `json:"upgradeStatus,omitempty"`
	UpgradeStartTime int    `json:"upgradeStartTime,omitempty"`
	UpgradeEndTime   int    `json:"upgradeEndTime,omitempty"`
}

type ECVMs struct {
	ID                int           `json:"id,omitempty"`
	Name              string        `json:"name,omitempty"`
	Status            []string      `json:"status,omitempty"`
	OperationalStatus string        `json:"operationalStatus,omitempty"`
	FormFactor        string        `json:"formFactor,omitempty"`
	CityGeoId         int           `json:"cityGeoId,omitempty"`
	NATIP             string        `json:"natIp,omitempty"`
	ZiaGateway        string        `json:"ziaGateway,omitempty"`
	ZpaBroker         string        `json:"zpaBroker,omitempty"`
	BuildVersion      string        `json:"buildVersion,omitempty"`
	LastUpgradeTime   int           `json:"lastUpgradeTime,omitempty"`
	UpgradeStatus     int           `json:"upgradeStatus,omitempty"`
	UpgradeStartTime  int           `json:"upgradeStartTime,omitempty"`
	UpgradeEndTime    int           `json:"upgradeEndTime,omitempty"`
	UpgradeDayOfWeek  int           `json:"upgradeDayOfWeek,omitempty"`
	ManagementNw      *ManagementNw `json:"managementNw,omitempty"`
	ECInstances       []ECInstances `json:"ecInstances,omitempty"`
}

type ManagementNw struct {
	ID             int    `json:"id,omitempty"`
	IPStart        string `json:"ipStart,omitempty"`
	IPEnd          string `json:"ipEnd,omitempty"`
	Netmask        string `json:"netmask,omitempty"`
	DefaultGateway string `json:"defaultGateway,omitempty"`
	NWType         string `json:"nwType,omitempty"`
	DNS            *DNS   `json:"dns,omitempty"`
}

type DNS struct {
	ID      int      `json:"id,omitempty"`
	IPs     []string `json:"ips,omitempty"`
	DNSType string   `json:"dnsType,omitempty"`
}

type ECInstances struct {
	ID             int           `json:"id,omitempty"`
	ECInstanceType string        `json:"ecInstanceType,omitempty"`
	ServiceIPs     *CommonIPs    `json:"serviceIps,omitempty"`
	LBIPAddr       *CommonIPs    `json:"lbIpAddr,omitempty"`
	OutGwIp        string        `json:"outGwIp,omitempty"`
	NatIP          string        `json:"natIp,omitempty"`
	DNSIP          []string      `json:"dnsIp,omitempty"`
	ServiceNw      *ManagementNw `json:"serviceNw,omitempty"`
	VirtualNw      *ManagementNw `json:"virtualNw,omitempty"`
}

type CommonIPs struct {
	IPStart string `json:"ipStart,omitempty"`
	IPEnd   string `json:"ipEnd,omitempty"`
}

// GetPageSize returns the page size.
func GetPageSize() int {
	return pageSize
}

func ReadAllPages[T any](ctx context.Context, client *zscaler.Client, endpoint string, list *[]T) error {
	if list == nil {
		return nil
	}
	page := 1
	if !strings.Contains(endpoint, "?") {
		endpoint += "?"
	}

	for {
		pageItems := []T{}
		err := client.ReadResource(ctx, fmt.Sprintf("%s&pageSize=%d&page=%d", endpoint, pageSize, page), &pageItems)
		if err != nil {
			return err
		}
		*list = append(*list, pageItems...)
		if len(pageItems) < pageSize {
			break
		}
		page++
	}
	return nil
}

func ReadPage[T any](ctx context.Context, client *ztw.Client, endpoint string, page int, list *[]T) error {
	if list == nil {
		return nil
	}

	// Parse the endpoint into a URL.
	u, err := url.Parse(endpoint)
	if err != nil {
		return fmt.Errorf("could not parse endpoint URL: %w", err)
	}

	// Get the existing query parameters and add new ones.
	q := u.Query()
	q.Set("pageSize", fmt.Sprintf("%d", pageSize))
	q.Set("page", fmt.Sprintf("%d", page))

	// Set the URL's RawQuery to the encoded query parameters.
	u.RawQuery = q.Encode()

	// Convert the URL back to a string and read the page.
	pageItems := []T{}
	err = client.Read(ctx, u.String(), &pageItems)
	if err != nil {
		return err
	}
	*list = pageItems
	return nil
}
