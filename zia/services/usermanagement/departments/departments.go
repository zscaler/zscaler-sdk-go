package departments

import (
	"fmt"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/common"
)

const (
	departmentEndpoint = "/departments"
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

func (service *Service) GetDepartments(departmentID int) (*Department, error) {
	var departments Department
	err := service.Client.Read(fmt.Sprintf("%s/%d", departmentEndpoint, departmentID), &departments)
	if err != nil {
		return nil, err
	}

	service.Client.Logger.Printf("[DEBUG]Returning departments from Get: %d", departments.ID)
	return &departments, nil
}

func (service *Service) GetDepartmentsByName(departmentName string) (*Department, error) {
	var departments []Department
	err := common.ReadAllPages(service.Client, departmentEndpoint, &departments)
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

func (service *Service) GetAll() ([]Department, error) {
	var departments []Department
	err := common.ReadAllPages(service.Client, departmentEndpoint, &departments)
	return departments, err
}
