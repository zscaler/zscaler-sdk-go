package provisioningurl

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v2/zcon/services"
	"github.com/zscaler/zscaler-sdk-go/v2/zcon/services/common"
	"github.com/zscaler/zscaler-sdk-go/v2/zcon/services/ecgroup"
	"github.com/zscaler/zscaler-sdk-go/v2/zcon/services/locationmanagement/locationtemplate"
)

const (
	provUrlEndpoint = "/provUrl"
)

type ECProvDetails struct {
	ID             int              `json:"id,omitempty"`
	Name           string           `json:"name,omitempty"`
	Desc           string           `json:"desc,omitempty"`
	ProvURL        string           `json:"provUrl,omitempty"`
	ProvUrlType    string           `json:"provUrlType,omitempty"`
	Status         string           `json:"status,omitempty"`
	LastModTime    int              `json:"lastModTime,omitempty"`
	ProvUrlData    *ProvUrlData     `json:"provUrlData,omitempty"`
	UsedInECGroups []common.UIDName `json:"usedInEcGroups,omitempty"`
	LastModUid     *common.UIDName  `json:"lastModUid,omitempty"`
}

type ProvUrlData struct {
	ZSCloudDomain      string                             `json:"zsCloudDomain,omitempty"`
	OrgID              int                                `json:"orgId,omitempty"`
	ConfigServer       string                             `json:"configServer,omitempty"`
	RegistrationServer string                             `json:"registrationServer,omitempty"`
	ApiServer          string                             `json:"apiServer,omitempty"`
	PacServer          string                             `json:"pacServer,omitempty"`
	CloudProviderType  string                             `json:"cloudProviderType,omitempty"`
	FormFactor         string                             `json:"formFactor,omitempty"`
	HyperVisors        string                             `json:"hyperVisors,omitempty"`
	CloudProvider      *common.UIDName                    `json:"cloudProvider,omitempty"`
	Location           *common.UIDName                    `json:"location,omitempty"`
	BCGroup            *ecgroup.EcGroup                   `json:"bcGroup,omitempty"`
	LocationTemplate   *locationtemplate.LocationTemplate `json:"locationTemplate,omitempty"`
}

type LBIPAddr struct {
	IPStart string `json:"ipStart,omitempty"`
	IPEnd   string `json:"ipEnd,omitempty"`
}

func Get(service *services.Service, provUrlID int) (*ECProvDetails, error) {
	var provUrl ECProvDetails
	err := service.Client.Read(fmt.Sprintf("%s/%d", provUrlEndpoint, provUrlID), &provUrl)
	if err != nil {
		return nil, err
	}

	log.Printf("Returning Provisining URL from Get: %d", provUrl.ID)
	return &provUrl, nil
}

func GetByName(service *services.Service, provUrlName string) (*ECProvDetails, error) {
	var provUrls []ECProvDetails
	// We are assuming this provisioning url name will be in the firsy 1000 obejcts
	err := service.Client.Read(fmt.Sprintf("%s?page=1&pageSize=1000", provUrlEndpoint), &provUrls)
	if err != nil {
		return nil, err
	}
	for _, provUr := range provUrls {
		if strings.EqualFold(provUr.Name, provUrlName) {
			return &provUr, nil
		}
	}
	return nil, fmt.Errorf("no provisioning url found with name: %s", provUrlName)
}

func Create(service *services.Service, provUrls *ECProvDetails) (*ECProvDetails, error) {
	resp, err := service.Client.Create(provUrlEndpoint, *provUrls)
	if err != nil {
		return nil, err
	}

	createdProvUrl, ok := resp.(*ECProvDetails)
	if !ok {
		return nil, errors.New("object returned from api was not a provisioning url pointer")
	}

	log.Printf("returning provisioning url from create: %d", createdProvUrl.ID)
	return createdProvUrl, nil
}

func Update(service *services.Service, provUrlID int, provUrls *ECProvDetails) (*ECProvDetails, *http.Response, error) {
	resp, err := service.Client.UpdateWithPut(fmt.Sprintf("%s/%d", provUrlEndpoint, provUrlID), *provUrls)
	if err != nil {
		return nil, nil, err
	}
	updatedProvUrl, _ := resp.(*ECProvDetails)

	log.Printf("returning provisioning url from Update: %d", updatedProvUrl.ID)
	return updatedProvUrl, nil, nil
}

func Delete(service *services.Service, provUrlID int) (*http.Response, error) {
	err := service.Client.Delete(fmt.Sprintf("%s/%d", provUrlEndpoint, provUrlID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}
