package partner_integrations

import (
	"context"
	"fmt"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/common"
)

const (
	publicCloudInfoEndpoint = "/ztw/api/v1/publicCloudInfo"
)

type WorkloadDiscoverySettings struct {

	// The ID of the Zscaler AWS account that is required in the AWS trust policy.
	TrustedAccountId string `json:"trustedAccountId,omitempty"`

	// The name of the trusted role in the Zscaler AWS account.
	TrustedRoleName string `json:"trustedRoleName,omitempty"`
}

type DiscoveryPermissions struct {
	// The discovery role for AWS account verification
	DiscoveryRole string `json:"discoveryRole"`

	// The external ID for AWS account verification
	ExternalID string `json:"externalId"`
}

func GetSupportedRegions(ctx context.Context, service *zscaler.Service) ([]common.SupportedRegions, error) {
	var regions []common.SupportedRegions
	err := common.ReadAllPages(ctx, service.Client, publicCloudInfoEndpoint+"/supportedRegions", &regions)

	service.Client.GetLogger().Printf("[DEBUG]Returning public supported regions from Get: %d items", len(regions))

	return regions, err
}

func GetSupportedRegionsByName(ctx context.Context, service *zscaler.Service, regionName string) (*common.SupportedRegions, error) {
	var regions []common.SupportedRegions
	err := common.ReadAllPages(ctx, service.Client, publicCloudInfoEndpoint+"/supportedRegions", &regions)
	if err != nil {
		return nil, err
	}
	for _, region := range regions {
		if strings.EqualFold(region.Name, regionName) {
			return &region, nil
		}
	}
	return nil, fmt.Errorf("no supported region found with name: %s", regionName)
}

func GetCloudFormationTemplateURL(ctx context.Context, service *zscaler.Service, awsAccountID *int) (string, error) {
	var templateURL string
	endpoint := publicCloudInfoEndpoint + "/cloudFormationTemplate"

	// Add optional awsAccountID parameter if provided
	if awsAccountID != nil {
		endpoint = fmt.Sprintf("%s?awsAccountId=%d", endpoint, *awsAccountID)
	}

	err := service.Client.ReadTextResource(ctx, endpoint, &templateURL)

	service.Client.GetLogger().Printf("[DEBUG]Returning cloud formation template URL: %s", templateURL)

	return templateURL, err
}

func GetWorkloadDiscoverySettings(ctx context.Context, service *zscaler.Service) ([]WorkloadDiscoverySettings, error) {
	var settings []WorkloadDiscoverySettings
	err := service.Client.ReadResource(ctx, "/ztw/api/v1/discoveryService/workloadDiscoverySettings", &settings)

	service.Client.GetLogger().Printf("[DEBUG] Returning workload discovery settings: %d items", len(settings))

	return settings, err
}

func UpdateDiscoveryPermissions(ctx context.Context, service *zscaler.Service, awsAccountID int, permissions *DiscoveryPermissions) error {
	endpoint := fmt.Sprintf("/ztw/api/v1/discoveryService/%d/permissions", awsAccountID)

	_, err := service.Client.UpdateWithPutResource(ctx, endpoint, *permissions)
	if err != nil {
		return err
	}

	service.Client.GetLogger().Printf("[DEBUG] verified discovery permissions for AWS account: %d", awsAccountID)
	return nil
}
