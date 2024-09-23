package common

const (
	DefaultPageSize = 30
	MaxPageSize     = 5000
)

type Pagination struct {
	PageSize int `json:"pagesize,omitempty" url:"pagesize,omitempty"`
	Page     int `json:"page,omitempty" url:"page,omitempty"`
}

// NewPagination creates a new Pagination struct with provided page size
// If page size is less than or equal to 0, it uses the default page size
// If page size is greater than the max page size, it uses the max page size
func NewPagination(pageSize int) Pagination {
	if pageSize <= 0 {
		pageSize = DefaultPageSize
	}
	if pageSize > MaxPageSize {
		pageSize = MaxPageSize
	}
	return Pagination{PageSize: pageSize}
}
