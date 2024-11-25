package departments

import (
	"context"
	"fmt"
	"strings"

	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler"
	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler/zia/services/common"
)

const (
	departmentEndpoint = "/zia/api/v1/departments"
)

type Department struct {
	// Department ID
	ID int `json:"id"`

	// Department name
	Name string `json:"name,omitempty"`

	// Identity provider (IdP) ID
	IdpID int `json:"idpId"`

	// Additional information about this department
	Comments string `json:"comments,omitempty"`
	Deleted  bool   `json:"deleted"`
}

func GetDepartments(ctx context.Context, service *zscaler.Service, departmentID int) (*Department, error) {
	var departments Department
	err := service.Client.Read(ctx, fmt.Sprintf("%s/%d", departmentEndpoint, departmentID), &departments)
	if err != nil {
		return nil, err
	}

	service.Client.GetLogger().Printf("[DEBUG]Returning departments from Get: %d", departments.ID)
	return &departments, nil
}

func GetDepartmentsByName(ctx context.Context, service *zscaler.Service, departmentName string) (*Department, error) {
	var departments []Department
	err := common.ReadAllPages(ctx, service.Client, departmentEndpoint+"?"+common.GetSortParams(common.SortField(service.SortBy), common.SortOrder(service.SortOrder)), &departments)
	if err != nil {
		return nil, err
	}
	for _, department := range departments {
		if strings.EqualFold(department.Name, departmentName) {
			return &department, nil
		}
	}
	return nil, fmt.Errorf("no department found with name: %s", departmentName)
}

func GetAll(ctx context.Context, service *zscaler.Service) ([]Department, error) {
	var departments []Department
	err := common.ReadAllPages(ctx, service.Client, departmentEndpoint+"?"+common.GetSortParams(common.SortField(service.SortBy), common.SortOrder(service.SortOrder)), &departments)
	return departments, err
}
