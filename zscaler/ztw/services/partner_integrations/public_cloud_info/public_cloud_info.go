package public_cloud_info

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/common"
)

const (
	publicCloudInfoEndpoint = "/ztw/api/v1/publicCloudInfo"
)

type PublicCloudInfo struct {
	// The unique ID of the AWS account.
	ID int `json:"id,omitempty"`

	// The name of the AWS account. The name must be non-null, non-empty, unique, and 128 characters or fewer in length.
	Name string `json:"name,omitempty"`

	// The cloud type. The default and mandatory value is AWS.
	CloudType string `json:"cloudType,omitempty"`

	// A unique external ID for the AWS account.
	ExternalID string `json:"externalId,omitempty"`

	// The date and time when the AWS account was last modified.
	LastModTime int `json:"lastModTime,omitempty"`

	// The last time the AWS account was synced.
	LastSyncTime int `json:"lastSyncTime,omitempty"`

	// // Indicates whether the provided credentials (external ID and AWS role name) have permission to access the AWS account.
	// PermissionStatus string `json:"permissionStatus,omitempty"`

	// An immutable reference to an entity, which consists of ID, name, etc.
	AccountGroups []common.IDNameExtensions `json:"accountGroups,omitempty"`

	// Automatically populated with the current ZIA admin user, after a successful POST or PUT request.
	LastModUser *common.CommonIDNameExternalID `json:"lastModUser,omitempty"`

	// The status and configuration details of the region where the workloads are deployed.
	RegionStatus []common.RegionStatus `json:"regionStatus,omitempty"`

	// Regions supported by Zscaler’s Tag Discovery Service.
	SupportedRegions []common.SupportedRegions `json:"supportedRegions,omitempty"`

	// The AWS account details.
	AccountDetails *AccountDetails `json:"accountDetails,omitempty"`
}

type AccountDetails struct {
	// The AWS account ID where workloads are deployed. The ID is non-null, non-empty, and unique, and contains 12 digits.
	Name string `json:"name,omitempty"`

	// The AWS account ID where workloads are deployed. The ID is non-null, non-empty, and unique, and contains 12 digits.
	AwsAccountID string `json:"awsAccountId,omitempty"`

	// The AWS trusting role in your account. The name is non-null, non-empty, and 64 characters or fewer in length.
	AwsRoleName string `json:"awsRoleName,omitempty"`

	// The resource name (ARN) of the AWS CloudWatch log group.
	CloudWatchGroupArn string `json:"cloudWatchGroupArn,omitempty"`

	// The name of the event bus that sends notifications to the Zscaler service using EventBridge.
	EventBusName string `json:"eventBusName,omitempty"`

	// (Optional) The unique external ID for the AWS account. If provided, it must match the externalId specified outside of accountDetails.
	ExternalID string `json:"externalId,omitempty"`

	// The type of log information. Supported types are INFO and ERROR. Supported Values: "INFO", "ERROR"
	LogInfoType string `json:"logInfoType,omitempty"`

	// Indicates whether logging is enabled for troubleshooting purposes.
	TroubleShootingLogging bool `json:"troubleShootingLogging,omitempty"`

	// (Optional) The ID of the Zscaler AWS account.
	TrustedAccountID string `json:"trustedAccountId,omitempty"`

	// (Optional) The name of the trusted role in the Zscaler AWS account.
	TrustedRole string `json:"trustedRole,omitempty"`
}

type PublicCloudInfoLite struct {

	// The unique ID of the supported region.
	ID int `json:"id,omitempty"`

	// The name of the supported region.
	Name string `json:"name,omitempty"`

	// The AWS account ID where workloads are deployed.
	AccountId string `json:"accountId,omitempty"`

	// The cloud type. The default and mandatory value is AWS. Supported Values: "AWS", "AZURE", "GCP"
	CloudType string `json:"cloudType,omitempty"`
}

func GetPublicCloudInfo(ctx context.Context, service *zscaler.Service, cloudID int) (*PublicCloudInfo, error) {
	var cloudInfo PublicCloudInfo
	err := service.Client.ReadResource(ctx, fmt.Sprintf("%s/%d", publicCloudInfoEndpoint, cloudID), &cloudInfo)
	if err != nil {
		return nil, err
	}

	service.Client.GetLogger().Printf("[DEBUG]Returning public cloud info from Get: %d", cloudInfo.ID)
	return &cloudInfo, nil
}

func GetByName(ctx context.Context, service *zscaler.Service, publicAccountName string) (*PublicCloudInfo, error) {
	var accountNames []PublicCloudInfo
	err := common.ReadAllPages(ctx, service.Client, publicCloudInfoEndpoint, &accountNames)
	if err != nil {
		return nil, err
	}
	for _, accountName := range accountNames {
		if strings.EqualFold(accountName.Name, publicAccountName) {
			return &accountName, nil
		}
	}
	return nil, fmt.Errorf("no public account info found with name: %s", publicAccountName)
}

func GetPublicCloudInfoLite(ctx context.Context, service *zscaler.Service) ([]PublicCloudInfoLite, error) {
	var locations []PublicCloudInfoLite
	err := common.ReadAllPages(ctx, service.Client, publicCloudInfoEndpoint+"/lite", &locations)

	service.Client.GetLogger().Printf("[DEBUG]Returning public cloud info from Get: %d items", len(locations))

	return locations, err
}

func GetAllPublicCloudInfo(ctx context.Context, service *zscaler.Service) ([]PublicCloudInfo, error) {
	var cloudInfo []PublicCloudInfo
	err := common.ReadAllPages(ctx, service.Client, publicCloudInfoEndpoint, &cloudInfo)

	service.Client.GetLogger().Printf("[DEBUG]Returning public cloud info from GetAll: %d items", len(cloudInfo))

	return cloudInfo, err
}

func GetPublicCloudInfoCount(ctx context.Context, service *zscaler.Service) (int, error) {
	var count int
	err := service.Client.ReadResource(ctx, publicCloudInfoEndpoint+"/count", &count)

	service.Client.GetLogger().Printf("[DEBUG]Returning public cloud info count: %d", count)

	return count, err
}

func CreatePublicCloudInfo(ctx context.Context, service *zscaler.Service, cloudInfo *PublicCloudInfo) (*PublicCloudInfo, error) {
	resp, err := service.Client.CreateResource(ctx, publicCloudInfoEndpoint, *cloudInfo)
	if err != nil {
		return nil, err
	}

	createdCloudInfo, ok := resp.(*PublicCloudInfo)
	if !ok {
		return nil, errors.New("object returned from api was not a cloud info pointer")
	}

	service.Client.GetLogger().Printf("[DEBUG]returning public cloud info from create: %d", createdCloudInfo.ID)
	return createdCloudInfo, nil
}

func UpdatePublicCloudInfo(ctx context.Context, service *zscaler.Service, awsAccountID int, cloudInfo *PublicCloudInfo) (*PublicCloudInfo, error) {
	resp, err := service.Client.UpdateWithPutResource(ctx, fmt.Sprintf("%s/%d", publicCloudInfoEndpoint, awsAccountID), *cloudInfo)
	if err != nil {
		return nil, err
	}

	updatedCloudInfo, ok := resp.(*PublicCloudInfo)
	if !ok {
		return nil, errors.New("object returned from api was not a public cloud info pointer")
	}

	service.Client.GetLogger().Printf("[DEBUG]returning updated public cloud info from update: %d", updatedCloudInfo.ID)
	return updatedCloudInfo, nil
}

func DeletePublicCloudInfo(ctx context.Context, service *zscaler.Service, awsAccountID int) error {
	err := service.Client.DeleteResource(ctx, fmt.Sprintf("%s/%d", publicCloudInfoEndpoint, awsAccountID))

	service.Client.GetLogger().Printf("[DEBUG] deleted public cloud info: %d", awsAccountID)

	return err
}

func UpdatePublicCloudChangeState(ctx context.Context, service *zscaler.Service, awsAccountID int, enable bool) error {
	endpoint := fmt.Sprintf("%s/%d/changeState?enable=%t", publicCloudInfoEndpoint, awsAccountID, enable)

	_, err := service.Client.UpdateWithPutResource(ctx, endpoint, nil)
	if err != nil {
		return err
	}

	service.Client.GetLogger().Printf("[DEBUG]changed state for public cloud info %d to enabled=%t", awsAccountID, enable)
	return nil
}

func GenerateExternalID(ctx context.Context, service *zscaler.Service, accountDetails *AccountDetails) (*AccountDetails, error) {
	resp, err := service.Client.CreateResource(ctx, publicCloudInfoEndpoint+"/generateExternalId", *accountDetails)
	if err != nil {
		return nil, err
	}

	generatedAccountDetails, ok := resp.(*AccountDetails)
	if !ok {
		return nil, errors.New("object returned from api was not an account details pointer")
	}

	service.Client.GetLogger().Printf("[DEBUG]generated external ID for AWS account: %s", generatedAccountDetails.AwsAccountID)
	return generatedAccountDetails, nil
}
