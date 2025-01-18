package zscaler

import (
	"log"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zcc"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa"
)

type (
	SortOrder string
	SortField string
)

const (
	ASCSortOrder          SortOrder = "ASC"
	DESCSortOrder                   = "DESC"
	IDSortField           SortField = "id"
	NameSortField                   = "name"
	CreationTimeSortField           = "creationTime"
	ModifiedTimeSortField           = "modifiedTime"
)

// Service defines the structure that contains the common client
type Service struct {
	Client        *Client // use the common Zscaler OneAPI Client here
	LegacyClient  *LegacyClient
	microTenantID *string
	// for some resources
	SortOrder SortOrder
	SortBy    SortField
}

// NewService is a generic function to instantiate a Service with the Zscaler OneAPI Client
func NewService(client *Client, legacyClient *LegacyClient) *Service {
	return &Service{
		Client:       client,
		LegacyClient: legacyClient,
		SortOrder:    ASCSortOrder,
		SortBy:       NameSortField,
	}
}

func (service *Service) WithMicroTenant(microTenantID string) *Service {
	var mid *string
	if microTenantID != "" {
		mid_ := microTenantID
		mid = &mid_
	}
	return &Service{
		Client:        service.Client,
		microTenantID: mid,
	}
}

func (service *Service) MicroTenantID() *string {
	return service.microTenantID
}

func (service *Service) WithSort(sortBy SortField, sortOrder SortOrder) *Service {
	c := Service{
		Client:    service.Client,
		SortOrder: service.SortOrder,
		SortBy:    service.SortBy,
	}
	if sortBy == IDSortField || sortBy == NameSortField || sortBy == CreationTimeSortField || sortBy == ModifiedTimeSortField {
		c.SortBy = sortBy
	}

	if sortOrder == ASCSortOrder || sortOrder == DESCSortOrder {
		c.SortOrder = sortOrder
	}
	return &c
}

func newLegacyHelper(conf ...ConfigSetter) (*Service, error) {
	cfg, err := NewConfiguration(
		conf...,
	)
	if err != nil {
		log.Fatalf("Error creating Zscaler configuration: %v", err)
		return nil, err
	}

	// Initialize the OneAPI client
	service, err := NewOneAPIClient(cfg)
	if err != nil {
		log.Fatalf("Error creating OneAPI client: %v", err)
		return nil, err
	}

	return service, nil
}

func NewLegacyZiaClient(config *zia.Configuration) (*Service, error) {
	ziaClient, err := zia.NewClient(config)
	if err != nil {
		log.Fatalf("Error creating ZIA client: %v", err)
		return nil, err
	}

	return newLegacyHelper(
		WithLegacyClient(true),
		WithZiaLegacyClient(ziaClient),
		WithDebug(config.Debug),
		// add other config mapping, if necessary
	)
}

func NewLegacyZccClient(config *zcc.Configuration) (*Service, error) {
	zccClient, err := zcc.NewClient(config)
	if err != nil {
		log.Fatalf("Error creating ZCC client: %v", err)
		return nil, err
	}

	return newLegacyHelper(
		WithLegacyClient(true),
		WithZccLegacyClient(zccClient),
		WithDebug(config.Debug),
		// add other config mapping, if necessary
	)
}

func NewLegacyZpaClient(config *zpa.Configuration) (*Service, error) {
	zpaClient, err := zpa.NewClient(config)
	if err != nil {
		log.Fatalf("Error creating ZPA client: %v", err)
		return nil, err
	}

	return newLegacyHelper(
		WithLegacyClient(true),
		WithZpaLegacyClient(zpaClient),
		WithDebug(config.Debug),
		// add other config mapping, if necessary
	)
}

type ScimService struct {
	ScimClient *zpa.ScimClient
}

// NewScimService initializes a SCIM-based ZPA Service with *zpa.ScimConfig
func NewScimService(scimClient *zpa.ScimClient) *ScimService {
	return &ScimService{ScimClient: scimClient}
}
