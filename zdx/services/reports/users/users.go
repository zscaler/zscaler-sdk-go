package devices

const (
	usersEndpoint = "/users"
)

/*
https://help.zscaler.com/zdx/reports#/apps/{appid}/users-get
Gets the list of all users and their devices that were used to access an application.
The endpoint allows to get the list of users based on the score category (i.e., Poor, Okay, or Good), location, department, or geoloaction.
If the time range is not specified, the endpoint defaults to the last 2 hours.
*/

type Users struct {
	ID         int    `json:"id"`
	Name       string `json:"name,omitempty"`
	Email      string `json:"email,omitempty"`
	NextOffSet string `json:"next_offset,omitempty"`
}

type Devices struct {
	ID           int            `json:"id"`
	Name         string         `json:"name,omitempty"`
	UserLocation []UserLocation `json:"geo_loc,omitempty"`
	ZSLocation   []ZSLocation   `json:"zs_loc,omitempty"`
}

type UserLocation struct {
	ID      int    `json:"id"`
	City    string `json:"city,omitempty"`
	Region  string `json:"region,omitempty"`
	Country string `json:"country,omitempty"`
	GeoType string `json:"geo_type,omitempty"`
}

type ZSLocation struct {
	ID   int    `json:"id"`
	Name string `json:"name,omitempty"`
}
