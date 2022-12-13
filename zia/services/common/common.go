package common

/*
const (

	DefaultPageSize = 1000

)

	type Pagination struct {
		PageSize int    `json:"pagesize,omitempty" url:"pagesize,omitempty"`
		Page     int    `json:"page,omitempty" url:"page,omitempty"`
		Search   string `json:"-" url:"-"`
		Search2  string `json:"search,omitempty" url:"search,omitempty"`
	}
*/
type IDNameExtensions struct {
	ID         int                    `json:"id,omitempty"`
	Name       string                 `json:"name,omitempty"`
	Extensions map[string]interface{} `json:"extensions,omitempty"`
}

type IDExtensions struct {
	ID         int                    `json:"id,omitempty"`
	Extensions map[string]interface{} `json:"extensions,omitempty"`
}

type UserGroups struct {
	ID       int    `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	IdpID    int    `json:"idp_id,omitempty"`
	Comments string `json:"comments,omitempty"`
}

type UserDepartment struct {
	ID       int    `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	IdpID    int    `json:"idp_id,omitempty"`
	Comments string `json:"comments,omitempty"`
	Deleted  bool   `json:"deleted,omitempty"`
}

type DeviceGroups struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type Devices struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

/*
func getAllPagesGeneric[T any](client *zia.Client, relativeURL string, page, pageSize int, searchQuery string) (int, []T, *http.Response, error) {
	var v struct {
		TotalPages interface{} `json:"totalPages"`
		List       []T         `json:"list"`
	}
	resp, err := client.Read("GET", relativeURL, Pagination{PageSize: pageSize, Page: page, Search2: searchQuery}, nil, &v)
	if err != nil {
		return 0, nil, resp, err
	}

	pages := fmt.Sprintf("%v", v.TotalPages)
	totalPages, _ := strconv.Atoi(pages)

	return totalPages, v.List, resp, nil
}

// GetAllPagesGeneric fetches all resources instead of just one single page
func GetAllPagesGeneric[T any](client *zia.Client, relativeURL, searchQuery string) ([]T, *http.Response, error) {
	totalPages, result, resp, err := getAllPagesGeneric[T](client, relativeURL, 1, DefaultPageSize, searchQuery)
	if err != nil {
		return nil, resp, err
	}
	var l []T
	for page := 2; page <= totalPages; page++ {
		totalPages, l, resp, err = getAllPagesGeneric[T](client, relativeURL, page, DefaultPageSize, searchQuery)
		if err != nil {
			return nil, resp, err
		}
		result = append(result, l...)
	}

	return result, resp, nil
}
*/
