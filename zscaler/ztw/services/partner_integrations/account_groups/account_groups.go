package account_groups

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/common"
)

const (
	accountGroupsEndpoint = "/ztw/api/v1/accountGroups"
)

type AccountGroups struct {
	ID                   int                       `json:"id,omitempty"`
	Name                 string                    `json:"name,omitempty"`
	Description          string                    `json:"description,omitempty"`
	CloudType            string                    `json:"cloudType,omitempty"`
	PublicCloudAccounts  []common.IDNameExtensions `json:"publicCloudAccounts,omitempty"`
	CloudConnectorGroups []common.IDNameExtensions `json:"cloudConnectorGroups,omitempty"`
}

func GetAccountGroup(ctx context.Context, service *zscaler.Service, awsAccountID int) ([]AccountGroups, error) {
	var accountGroups []AccountGroups
	err := common.ReadAllPages(ctx, service.Client, fmt.Sprintf("%s/%d", accountGroupsEndpoint, awsAccountID), &accountGroups)
	if err != nil {
		return nil, err
	}

	service.Client.GetLogger().Printf("[DEBUG] Returning account groups from Get: %d items", len(accountGroups))
	return accountGroups, nil
}

func GetByName(ctx context.Context, service *zscaler.Service, accountGroupsName string) (*AccountGroups, error) {
	var accountGroups []AccountGroups
	err := common.ReadAllPages(ctx, service.Client, accountGroupsEndpoint, &accountGroups)
	if err != nil {
		return nil, err
	}
	for _, accountGroup := range accountGroups {
		if strings.EqualFold(accountGroup.Name, accountGroupsName) {
			return &accountGroup, nil
		}
	}
	return nil, fmt.Errorf("no account group found with name: %s", accountGroupsName)
}

func GetAccountGroupsLite(ctx context.Context, service *zscaler.Service) ([]AccountGroups, error) {
	var accountGroups []AccountGroups
	err := common.ReadAllPages(ctx, service.Client, accountGroupsEndpoint+"/lite", &accountGroups)

	service.Client.GetLogger().Printf("[DEBUG]Returning account groups lite from Get: %d items", len(accountGroups))

	return accountGroups, err
}

func CreateAccountGroups(ctx context.Context, service *zscaler.Service, accountGroup *AccountGroups) (*AccountGroups, error) {
	resp, err := service.Client.CreateResource(ctx, accountGroupsEndpoint, *accountGroup)
	if err != nil {
		return nil, err
	}

	createdAccountGroup, ok := resp.(*AccountGroups)
	if !ok {
		return nil, errors.New("object returned from api was not an account group pointer")
	}

	service.Client.GetLogger().Printf("[DEBUG]returning account group from create: %d", createdAccountGroup.ID)
	return createdAccountGroup, nil
}

func UpdateAccountGroups(ctx context.Context, service *zscaler.Service, awsAccountID int, accountGroup *AccountGroups) (*AccountGroups, error) {
	resp, err := service.Client.UpdateWithPutResource(ctx, fmt.Sprintf("%s/%d", accountGroupsEndpoint, awsAccountID), *accountGroup)
	if err != nil {
		return nil, err
	}

	updatedAccountGroup, ok := resp.(*AccountGroups)
	if !ok {
		return nil, errors.New("object returned from api was not an account group pointer")
	}

	service.Client.GetLogger().Printf("[DEBUG]returning updated account group from update: %d", updatedAccountGroup.ID)
	return updatedAccountGroup, nil
}

func DeleteAccountGroups(ctx context.Context, service *zscaler.Service, awsAccountID int) error {
	err := service.Client.DeleteResource(ctx, fmt.Sprintf("%s/%d", accountGroupsEndpoint, awsAccountID))

	service.Client.GetLogger().Printf("[DEBUG]deleted account group: %d", awsAccountID)

	return err
}

func GetAllAccountGroups(ctx context.Context, service *zscaler.Service) ([]AccountGroups, error) {
	var accountGroups []AccountGroups
	err := common.ReadAllPages(ctx, service.Client, accountGroupsEndpoint, &accountGroups)

	service.Client.GetLogger().Printf("[DEBUG]Returning account groups from Get: %d items", len(accountGroups))

	return accountGroups, err
}
