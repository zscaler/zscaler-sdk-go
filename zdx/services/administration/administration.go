package administration

const (
	departmentsEndpoint = "/administration/departments"
	locationsEndpoint   = "/administration/locations"
)

/*
https://help.zscaler.com/zdx/administration-0#/administration/departments-get
Gets configured departments.

https://help.zscaler.com/zdx/administration-0#/administration/locations-get
Gets Zscaler locations. All configured Zscaler locations are returned if the search filters is not specified.
*/
type Departments struct {
	ID   int    `json:"id"`
	Name string `json:"name,omitempty"`
}

type Locations struct {
	ID   int    `json:"id"`
	Name string `json:"name,omitempty"`
}
