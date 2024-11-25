package public_cloud_account

import (
	"fmt"
	"log"

	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler/zcon/services"
	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler/zcon/services/common"
)

const (
	publicCloudEndpoint      = "/publicCloudAccountDetails"
	publicCloudEndpointLite  = "/publicCloudAccountDetails/lite"
	publicCloudAccountStatus = "/publicCloudAccountIdStatus"
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
func GetAccountID(service *services.Service, accountID int) (*PublicCloudAccountDetails, error) {
	var cloudAccount PublicCloudAccountDetails
	err := service.Client.Read(fmt.Sprintf("%s/%d", publicCloudEndpoint, accountID), &cloudAccount)
	if err != nil {
		return nil, err
	}

	log.Printf("Returning Available Cloud Account from Get: %d", cloudAccount.ID)
	return &cloudAccount, nil
}

// GetLite returns all available accounts without filtering by ID
func GetLite(service *services.Service) ([]PublicCloudAccountDetails, error) {
	var cloudAccounts []PublicCloudAccountDetails
	err := service.Client.Read(publicCloudEndpointLite, &cloudAccounts)
	if err != nil {
		return nil, err
	}

	log.Printf("Returning all available cloud accounts from GetLite")
	return cloudAccounts, nil
}

// GetAccountStatus returns a status payload directly
func GetAccountStatus(service *services.Service) (*PublicCloudAccountIDStatus, error) {
	var accountStatus PublicCloudAccountIDStatus
	err := service.Client.Read(publicCloudAccountStatus, &accountStatus)
	if err != nil {
		return nil, err
	}

	log.Printf("Returning account status from GetAccountStatus")
	return &accountStatus, nil
}

func GetAll(service *services.Service) ([]PublicCloudAccountDetails, error) {
	var accounts []PublicCloudAccountDetails
	err := common.ReadAllPages(service.Client, publicCloudEndpoint, &accounts)
	return accounts, err
}
