package scimgroup

import (
	"github.com/zscaler/zscaler-sdk-go/v2/zpa"
)

type SortOrder string
type SortField string

const (
	ASCSortOrder          SortOrder = "ASC"
	DESCSortOrder                   = "DESC"
	IDSortField           SortField = "id"
	NameSortField                   = "name"
	CreationTimeSortField           = "creationTime"
	ModifiedTimeSortField           = "modifiedTime"
)

type Service struct {
	Client    *zpa.Client
	sortOrder SortOrder
	sortBy    SortField
}

func New(c *zpa.Client) *Service {
	return &Service{
		Client:    c,
		sortOrder: ASCSortOrder,
		sortBy:    NameSortField,
	}
}

func (service *Service) WithSort(sortBy SortField, sortOrder SortOrder) *Service {
	c := Service{
		Client:    service.Client,
		sortOrder: service.sortOrder,
		sortBy:    service.sortBy,
	}
	if sortBy == IDSortField || sortBy == NameSortField || sortBy == CreationTimeSortField || sortBy == ModifiedTimeSortField {
		c.sortBy = sortBy
	}

	if sortOrder == ASCSortOrder || sortOrder == DESCSortOrder {
		c.sortOrder = sortOrder
	}
	return &c
}
