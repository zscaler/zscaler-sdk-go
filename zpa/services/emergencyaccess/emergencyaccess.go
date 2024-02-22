package emergencyaccess

import (
	"fmt"
	"net/http"
)

const (
	mgmtConfig              = "/mgmtconfig/v1/admin/customers/"
	emergencyAccessEndpoint = "/emergencyAccess/user"
)

type EmergencyAccess struct {
	ActivatedOn       string `json:"activatedOn,omitempty"`
	AllowedActivate   bool   `json:"allowedActivate"`
	AllowedDeactivate bool   `json:"allowedDeactivate"`
	EmailID           string `json:"emailId,omitempty"`
	FirstName         string `json:"firstName,omitempty"`
	LastLoginTime     string `json:"lastLoginTime,omitempty"`
	LastName          string `json:"lastName,omitempty"`
	UpdateEnabled     bool   `json:"updateEnabled"`
	UserID            string `json:"userId,omitempty"`
	UserStatus        string `json:"userStatus,omitempty"`
	ActivateNow       bool   `json:"activateNow,omitempty" url:"activateNow,omitempty"`
}

func (service *Service) Get(profileID string) (*EmergencyAccess, *http.Response, error) {
	v := new(EmergencyAccess)
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.Config.CustomerID+emergencyAccessEndpoint, profileID)
	resp, err := service.Client.NewRequestDo("GET", relativeURL, nil, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

// func (service *Service) GetByEmailID(emailID string) (*EmergencyAccess, *http.Response, error) {
// 	// Use the GetAll function to retrieve all EmergencyAccess records
// 	list, resp, err := service.GetAll()
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	// Filter the retrieved list for the specific emailID
// 	for _, emgAccess := range list {
// 		if strings.EqualFold(emgAccess.EmailId, emailID) {
// 			return &emgAccess, resp, nil
// 		}
// 	}

// 	// If no matching record is found, return an error
// 	return nil, resp, fmt.Errorf("no emergency access record found with email ID '%s'", emailID)
// }

func (service *Service) Create(emergencyAccess *EmergencyAccess) (*EmergencyAccess, *http.Response, error) {
	// Constructing the relative URL correctly
	emergencyAccess.ActivateNow = false // Ensure this parameter is set as intended
	// Constructing the relative URL correctly
	relativeURL := fmt.Sprintf("%s%s%s", mgmtConfig, service.Client.Config.CustomerID, emergencyAccessEndpoint)
	v := new(EmergencyAccess)

	// Making the request without common.Filter
	resp, err := service.Client.NewRequestDo("POST", relativeURL, nil, emergencyAccess, v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func (service *Service) Update(userID string, emergencyAccess *EmergencyAccess) (*http.Response, error) {
	path := fmt.Sprintf("%v/%v", mgmtConfig+service.Client.Config.CustomerID+emergencyAccessEndpoint, userID)
	resp, err := service.Client.NewRequestDo("PUT", path, nil, emergencyAccess, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

// PUT - /mgmtconfig/v1/admin/customers/{customerId}/emergencyAccess/user/{userId}/activate
func (service *Service) Activate(userID string) (*http.Response, error) {
	path := fmt.Sprintf("%s/%s/activate", mgmtConfig+service.Client.Config.CustomerID+emergencyAccessEndpoint, userID)

	// Making the request without common.Filter
	resp, err := service.Client.NewRequestDo("PUT", path, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// PUT - /mgmtconfig/v1/admin/customers/{customerId}/emergencyAccess/user/{userId}/deactivate
func (service *Service) Deactivate(userID string) (*http.Response, error) {
	path := fmt.Sprintf("%s/%s/deactivate", mgmtConfig+service.Client.Config.CustomerID+emergencyAccessEndpoint, userID)

	// Making the request without common.Filter
	resp, err := service.Client.NewRequestDo("PUT", path, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// func (service *Service) GetAll() ([]EmergencyAccess, *http.Response, error) {
// 	relativeURL := fmt.Sprintf("%s%s%s", mgmtConfig, service.Client.Config.CustomerID, emergencyAccessEndpoint+"s") // Correct endpoint
// 	pageSize := 20                                                                                                  // Define the pageSize as needed
// 	initialPageId := ""                                                                                             // Start without a pageId or as required

// 	return GetAllEmergencyAccessUsers(service.Client, relativeURL, pageSize, initialPageId)
// }

// func fetchEmergencyAccessUsersPage(client *zpa.Client, fullURL string) (*http.Response, error) {
// 	// Directly use the fullURL, which is expected to be correctly formatted beforehand
// 	return client.NewRequestDo("GET", fullURL, nil, nil, nil)
// }

// func GetAllEmergencyAccessUsers(client *zpa.Client, baseRelativeURL string, pageSize int, initialPageId string) ([]EmergencyAccess, *http.Response, error) {
// 	var allUsers []EmergencyAccess
// 	var lastResponse *http.Response
// 	pageId := initialPageId

// 	for {
// 		// Construct the URL for each request to avoid duplication and encoding issues
// 		var fullURL string
// 		if pageId != "" {
// 			fullURL = fmt.Sprintf("%s?pageSize=%d&pageId=%s", baseRelativeURL, pageSize, pageId)
// 		} else {
// 			fullURL = fmt.Sprintf("%s?pageSize=%d", baseRelativeURL, pageSize)
// 		}

// 		resp, err := fetchEmergencyAccessUsersPage(client, fullURL)
// 		if err != nil {
// 			return nil, lastResponse, err
// 		}

// 		// Assume this struct matches the expected JSON response structure
// 		var pageData struct {
// 			Items    []EmergencyAccess `json:"items"`
// 			NextPage string            `json:"nextPage"`
// 		}

// 		bodyBytes, err := ioutil.ReadAll(resp.Body)
// 		resp.Body.Close() // Ensure the body is always closed after reading
// 		if err != nil {
// 			return nil, resp, fmt.Errorf("error reading response body: %w", err)
// 		}

// 		if err := json.Unmarshal(bodyBytes, &pageData); err != nil {
// 			return nil, resp, fmt.Errorf("error unmarshalling response: %w", err)
// 		}

// 		allUsers = append(allUsers, pageData.Items...)
// 		if pageData.NextPage == "" {
// 			break // Exit the loop if there are no more pages
// 		}

// 		// Update pageId for the next iteration
// 		pageId = pageData.NextPage
// 		lastResponse = resp
// 	}

// 	return allUsers, lastResponse, nil
// }
