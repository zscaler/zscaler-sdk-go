package departments

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/common"
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

	Deleted bool `json:"deleted"`
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

func GetDepartmentLite(ctx context.Context, service *zscaler.Service, departmentID int) (*Department, error) {
	var departments Department
	err := service.Client.Read(ctx, fmt.Sprintf("%s/%d", departmentEndpoint+"/lite", departmentID), &departments)
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

func Create(ctx context.Context, service *zscaler.Service, departmentID *Department) (*Department, *http.Response, error) {
	resp, err := service.Client.Create(ctx, departmentEndpoint, *departmentID)
	if err != nil {
		return nil, nil, err
	}

	createdDept, ok := resp.(*Department)
	if !ok {
		return nil, nil, errors.New("object returned from api was not a department pointer")
	}

	service.Client.GetLogger().Printf("[DEBUG]returning new department from create: %d", createdDept.ID)
	return createdDept, nil, nil
}

func Update(ctx context.Context, service *zscaler.Service, departmentID int, depts *Department) (*Department, *http.Response, error) {
	resp, err := service.Client.UpdateWithPut(ctx, fmt.Sprintf("%s/%d", departmentEndpoint, departmentID), *depts)
	if err != nil {
		return nil, nil, err
	}
	updatedDept, _ := resp.(*Department)

	service.Client.GetLogger().Printf("[DEBUG]returning updates department from update: %d", updatedDept.ID)
	return updatedDept, nil, nil
}

func Delete(ctx context.Context, service *zscaler.Service, departmentID int) (*http.Response, error) {
	err := service.Client.Delete(ctx, fmt.Sprintf("%s/%d", departmentEndpoint, departmentID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func GetAll(ctx context.Context, service *zscaler.Service) ([]Department, error) {
	var departments []Department
	err := common.ReadAllPages(ctx, service.Client, departmentEndpoint+"?"+common.GetSortParams(common.SortField(service.SortBy), common.SortOrder(service.SortOrder)), &departments)
	return departments, err
}

func GetAllLite(ctx context.Context, service *zscaler.Service) ([]Department, error) {
	var depts []Department
	err := common.ReadAllPages(ctx, service.Client, departmentEndpoint+"/lite", &depts)
	return depts, err
}
