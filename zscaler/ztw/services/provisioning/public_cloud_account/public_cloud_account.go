package public_cloud_account

import (
	"context"
	"fmt"
	"log"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/common"
)

const (
	publicCloudEndpoint      = "/ztw/api/v1/publicCloudAccountDetails"
	publicCloudEndpointLite  = "/ztw/api/v1/publicCloudAccountDetails/lite"
	publicCloudAccountStatus = "/ztw/api/v1/publicCloudAccountIdStatus"
)

type PublicCloudAccountDetails struct {
	// Internal ID of public cloud account/subscription.
	ID int `json:"id,omitempty"`

	// Account or subscription ID of public cloud account.
	AccountID string `json:"accountId,omitempty"`

	// Public cloud platform (AWS or Azure)
	PlatformID string `json:"platformId,omitempty"`
}

type PublicCloudAccountIDStatus struct {
	// Indicates whether public cloud account is enabled.
	AccountIdEnabled bool `json:"accountIdEnabled,omitempty"`

	// Indicates whether public cloud subscription is enabled.
	SubIDEnabled bool `json:"subIdEnabled,omitempty"`

	// Indicates whether public cloud subscription is enabled.
	ProjectIdEnabled bool `json:"projectIdEnabled,omitempty"`
}

// GetAccountID remains the same
func GetAccountID(ctx context.Context, service *zscaler.Service, accountID int) (*PublicCloudAccountDetails, error) {
	var cloudAccount PublicCloudAccountDetails
	err := service.Client.ReadResource(ctx, fmt.Sprintf("%s/%d", publicCloudEndpoint, accountID), &cloudAccount)
	if err != nil {
		return nil, err
	}

	log.Printf("Returning Available Cloud Account from Get: %d", cloudAccount.ID)
	return &cloudAccount, nil
}

// GetLite returns all available accounts without filtering by ID
func GetLite(ctx context.Context, service *zscaler.Service) ([]PublicCloudAccountDetails, error) {
	var cloudAccounts []PublicCloudAccountDetails
	err := service.Client.ReadResource(ctx, publicCloudEndpointLite, &cloudAccounts)
	if err != nil {
		return nil, err
	}

	log.Printf("Returning all available cloud accounts from GetLite")
	return cloudAccounts, nil
}

// GetAccountStatus returns a status payload directly
func GetAccountStatus(ctx context.Context, service *zscaler.Service) (*PublicCloudAccountIDStatus, error) {
	var accountStatus PublicCloudAccountIDStatus
	err := service.Client.ReadResource(ctx, publicCloudAccountStatus, &accountStatus)
	if err != nil {
		return nil, err
	}

	log.Printf("Returning account status from GetAccountStatus")
	return &accountStatus, nil
}

func GetAll(ctx context.Context, service *zscaler.Service) ([]PublicCloudAccountDetails, error) {
	var accounts []PublicCloudAccountDetails
	err := common.ReadAllPages(ctx, service.Client, publicCloudEndpoint, &accounts)
	return accounts, err
}
