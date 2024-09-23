package zscaler

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
	microTenantID *string
	// for some resources
	SortOrder SortOrder
	SortBy    SortField
}

// NewService is a generic function to instantiate a Service with the Zscaler OneAPI Client
func NewService(client *Client) *Service {
	return &Service{
		Client:    client,
		SortOrder: ASCSortOrder,
		SortBy:    NameSortField,
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
