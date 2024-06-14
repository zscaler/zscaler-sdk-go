package emergencyaccess

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/common"
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

func Get(service *services.Service, userID string) (*EmergencyAccess, *http.Response, error) {
	v := new(EmergencyAccess)
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.Config.CustomerID+emergencyAccessEndpoint, userID)
	resp, err := service.Client.NewRequestDo("GET", relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()}, nil, v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func GetByEmailID(service *services.Service, emailID string) (*EmergencyAccess, *http.Response, error) {
	// Use the GetAll function to retrieve all EmergencyAccess records
	list, resp, err := GetAll(service)
	if err != nil {
		return nil, nil, err
	}

	// Filter the retrieved list for the specific emailID
	for _, emgAccess := range list {
		if strings.EqualFold(emgAccess.EmailID, emailID) {
			return &emgAccess, resp, nil
		}
	}
	return nil, resp, fmt.Errorf("no emergency access record found with email ID '%s'", emailID)
}

func Create(service *services.Service, emergencyAccess *EmergencyAccess) (*EmergencyAccess, *http.Response, error) {
	emergencyAccess.ActivateNow = false
	relativeURL := fmt.Sprintf("%s%s%s", mgmtConfig, service.Client.Config.CustomerID, emergencyAccessEndpoint)
	v := new(EmergencyAccess)
	resp, err := service.Client.NewRequestDo("POST", relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()}, emergencyAccess, v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func Update(service *services.Service, userID string, emergencyAccess *EmergencyAccess) (*http.Response, error) {
	path := fmt.Sprintf("%v/%v", mgmtConfig+service.Client.Config.CustomerID+emergencyAccessEndpoint, userID)
	resp, err := service.Client.NewRequestDo("PUT", path, common.Filter{MicroTenantID: service.MicroTenantID()}, emergencyAccess, nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// PUT - /mgmtconfig/v1/admin/customers/{customerId}/emergencyAccess/user/{userId}/activate
func Activate(service *services.Service, userID string) (*http.Response, error) {
	path := fmt.Sprintf("%s/%s/activate", mgmtConfig+service.Client.Config.CustomerID+emergencyAccessEndpoint, userID)
	resp, err := service.Client.NewRequestDo("PUT", path, common.Filter{MicroTenantID: service.MicroTenantID()}, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// PUT - /mgmtconfig/v1/admin/customers/{customerId}/emergencyAccess/user/{userId}/deactivate
func Deactivate(service *services.Service, userID string) (*http.Response, error) {
	path := fmt.Sprintf("%s/%s/deactivate", mgmtConfig+service.Client.Config.CustomerID+emergencyAccessEndpoint, userID)
	resp, err := service.Client.NewRequestDo("PUT", path, common.Filter{MicroTenantID: service.MicroTenantID()}, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func GetAll(service *services.Service) ([]EmergencyAccess, *http.Response, error) {
	relativeURL := fmt.Sprintf("%s%s%ss", mgmtConfig, service.Client.Config.CustomerID, emergencyAccessEndpoint) // Correct endpoint
	pageSize := 500                                                                                              // Define the pageSize as needed
	initialPageId := ""                                                                                          // Start without a pageId or as required

	return GetAllEmergencyAccessUsers(service, relativeURL, pageSize, initialPageId)
}

func fetchEmergencyAccessUsersPage(service *services.Service, fullURL string) (*http.Response, error) {
	return service.Client.NewRequestDo("GET", fullURL, common.Filter{MicroTenantID: service.MicroTenantID()}, nil, nil)
}

func GetAllEmergencyAccessUsers(service *services.Service, baseRelativeURL string, pageSize int, initialPageId string) ([]EmergencyAccess, *http.Response, error) {
	var allUsers []EmergencyAccess
	var lastResponse *http.Response
	pageId := initialPageId

	for {
		// Construct the URL for each request to avoid duplication and encoding issues
		var fullURL string
		if pageId != "" {
			fullURL = fmt.Sprintf("%s?pageSize=%d&pageId=%s", baseRelativeURL, pageSize, pageId)
		} else {
			fullURL = fmt.Sprintf("%s?pageSize=%d", baseRelativeURL, pageSize)
		}

		resp, err := fetchEmergencyAccessUsersPage(service, fullURL)
		if err != nil {
			return nil, lastResponse, err
		}
		// Assume this struct matches the expected JSON response structure
		var pageData struct {
			Items    []EmergencyAccess `json:"items"`
			NextPage string            `json:"nextPage"`
		}

		bodyBytes, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close() // Ensure the body is always closed after reading
		if err != nil {
			return nil, resp, fmt.Errorf("error reading response body: %w", err)
		}

		if err := json.Unmarshal(bodyBytes, &pageData); err != nil {
			return nil, resp, fmt.Errorf("error unmarshalling response: %w", err)
		}

		allUsers = append(allUsers, pageData.Items...)
		if pageData.NextPage == "" {
			break // Exit the loop if there are no more pages
		}

		// Update pageId for the next iteration
		pageId = pageData.NextPage
		lastResponse = resp
	}

	return allUsers, lastResponse, nil
}
