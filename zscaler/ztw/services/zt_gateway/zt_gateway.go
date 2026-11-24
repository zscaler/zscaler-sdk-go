package zt_gateway

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/common"
)

const (
	ztGatewayEndpoint                 = "/ztw/api/v1/ztGateway"
	ztGatewayCountEndpoint            = "/ztw/api/v1/ztGateway/count"
	ztGatewayLicenseInfoEndpoint      = "/ztw/api/v1/ztGateway/licenseInfo"
	ztGatewayLiteEndpoint             = "/ztw/api/v1/ztGateway/lite"
	ztGatewaySupportedRegionsEndpoint = "/ztw/api/v1/ztGateway/supportedRegions"
	ztGatewayTrafficTestEndpoint      = "/ztw/api/v1/ztGateway/trafficTest"
	ztGatewayConfigCompareEndpoint    = "/ztw/api/v1/ztGateway/config/compare"
	ztGatewayConfigStatusEndpoint     = "/ztw/api/v1/ztGateway/config/status"
	ztGatewayEventsEndpoint           = "/ztw/api/v1/ztGateway/events"
	ztGatewayStatusEndpoint           = "/ztw/api/v1/ztGateway/status"
	ztGatewayTrafficTestEnvEndpoint   = "/ztw/api/v1/ztGateway/trafficTest/env"
	ztGatewayAnalyticMetricsEndpoint  = "/ztw/api/v1/ztGateway/analytic/metrics"
)

type ZtGateway struct {
	AvailabilityZoneIds          []string                 `json:"availabilityZoneIds,omitempty"`
	CreateTime                   int                      `json:"createTime,omitempty"`
	EgressIPs                    map[string]interface{}   `json:"egressIPs,omitempty"`
	HealthStatus                 string                   `json:"healthStatus,omitempty"`
	ID                           int                      `json:"id"`
	LastModTime                  int                      `json:"lastModTime,omitempty"`
	LastModUid                   LastModUid               `json:"lastModUid,omitempty"`
	Name                         string                   `json:"name,omitempty"`
	Platform                     string                   `json:"platform,omitempty"`
	ProvData                     ProvData                 `json:"provData,omitempty"`
	Region                       string                   `json:"region,omitempty"`
	EndpointServiceName          string                   `json:"endpointServiceName,omitempty"`
	Endpoints                    []map[string]interface{} `json:"endpoints,omitempty"`
	EndpointsCount               string                   `json:"endpointsCount,omitempty"`
	InterceptDeploymentGroup     string                   `json:"interceptDeploymentGroup,omitempty"`
	InterceptEndpointGroups      []map[string]interface{} `json:"interceptEndpointGroups,omitempty"`
	InterceptEndpointGroupsCount string                   `json:"interceptEndpointGroupsCount,omitempty"`
}

type LastModUid struct {
	Deleted bool   `json:"deleted,omitempty"`
	ID      int    `json:"id"`
	Name    string `json:"name,omitempty"`
}

type ProvData struct {
	AccountGroups         []AccountGroup   `json:"accountGroups,omitempty"`
	AdditionalAwsAccounts []string         `json:"additionalAwsAccounts,omitempty"`
	AllowedAccounts       []AllowedAccount `json:"allowedAccounts,omitempty"`
	DrTestingMode         bool             `json:"drTestingMode,omitempty"`
	GlobalProfile         GlobalProfile    `json:"globalProfile,omitempty"`
	LocationName          string           `json:"locationName,omitempty"`
	LocationTemplate      LocationTemplate `json:"locationTemplate,omitempty"`
	IamPrincipals         []IamPrincipal   `json:"iamPrincipals,omitempty"`
}

type AccountGroup struct {
	Deleted bool   `json:"deleted,omitempty"`
	ID      int    `json:"id"`
	Name    string `json:"name,omitempty"`
}

type AllowedAccount struct {
	Deleted bool   `json:"deleted,omitempty"`
	ID      int    `json:"id"`
	Name    string `json:"name,omitempty"`
}

type GlobalProfile struct {
	Deleted bool   `json:"deleted,omitempty"`
	ID      int    `json:"id"`
	Name    string `json:"name,omitempty"`
}

type LocationTemplate struct {
	ID          int                            `json:"id"`
	Name        string                         `json:"name,omitempty"`
	Desc        string                         `json:"desc,omitempty"`
	Template    Template                       `json:"template,omitempty"`
	Editable    bool                           `json:"editable,omitempty"`
	LastModUid  *common.CommonIDNameExternalID `json:"lastModUid,omitempty"`
	LastModTime int                            `json:"lastModTime,omitempty"`
}

type Template struct {
	TemplatePrefix                      string `json:"templatePrefix,omitempty"`
	XffForwardEnabled                   bool   `json:"xffForwardEnabled,omitempty"`
	AuthRequired                        bool   `json:"authRequired,omitempty"`
	CautionEnabled                      bool   `json:"cautionEnabled,omitempty"`
	AupEnabled                          bool   `json:"aupEnabled,omitempty"`
	AupTimeoutInDays                    int    `json:"aupTimeoutInDays,omitempty"`
	OfwEnabled                          bool   `json:"ofwEnabled,omitempty"`
	IpsControl                          bool   `json:"ipsControl,omitempty"`
	EnforceBandwidthControl             bool   `json:"enforceBandwidthControl,omitempty"`
	UpBandwidth                         int    `json:"upBandwidth,omitempty"`
	DnBandwidth                         int    `json:"dnBandwidth,omitempty"`
	DisplayTimeUnit                     string `json:"displayTimeUnit,omitempty"`
	IdleTimeInMinutes                   int    `json:"idleTimeInMinutes,omitempty"`
	SurrogateIPEnforcedForKnownBrowsers bool   `json:"surrogateIPEnforcedForKnownBrowsers,omitempty"`
	SurrogateRefreshTimeUnit            string `json:"surrogateRefreshTimeUnit,omitempty"`
	SurrogateRefreshTimeInMinutes       int    `json:"surrogateRefreshTimeInMinutes,omitempty"`
	SurrogateIP                         bool   `json:"surrogateIP,omitempty"`
}

type IamPrincipal struct {
	AddedOn int    `json:"addedOn,omitempty"`
	Type    string `json:"type,omitempty"`
	Value   string `json:"value,omitempty"`
}

type RunRequest struct {
	TestID string `json:"testId,omitempty"`
}

type MetricsRequest struct {
	EndTime   string `json:"endTime,omitempty"`
	Period    string `json:"period,omitempty"`
	StartTime string `json:"startTime,omitempty"`
	Stat      string `json:"stat,omitempty"`
	Type      string `json:"type,omitempty"`
}

func GetAll(ctx context.Context, service *zscaler.Service) ([]ZtGateway, error) {
	var result []ZtGateway
	err := common.ReadAllPages(ctx, service.Client, ztGatewayEndpoint, &result)
	return result, err
}

func GetAllCount(ctx context.Context, service *zscaler.Service) ([]ZtGateway, error) {
	var result []ZtGateway
	err := common.ReadAllPages(ctx, service.Client, ztGatewayCountEndpoint, &result)
	return result, err
}

func GetAllLicenseInfo(ctx context.Context, service *zscaler.Service) ([]ZtGateway, error) {
	var result []ZtGateway
	err := common.ReadAllPages(ctx, service.Client, ztGatewayLicenseInfoEndpoint, &result)
	return result, err
}

func GetAllLite(ctx context.Context, service *zscaler.Service) ([]ZtGateway, error) {
	var result []ZtGateway
	err := common.ReadAllPages(ctx, service.Client, ztGatewayLiteEndpoint, &result)
	return result, err
}

func GetAllSupportedRegions(ctx context.Context, service *zscaler.Service) ([]ZtGateway, error) {
	var result []ZtGateway
	err := common.ReadAllPages(ctx, service.Client, ztGatewaySupportedRegionsEndpoint, &result)
	return result, err
}

func GetAllTrafficTest(ctx context.Context, service *zscaler.Service) ([]ZtGateway, error) {
	var result []ZtGateway
	err := common.ReadAllPages(ctx, service.Client, ztGatewayTrafficTestEndpoint, &result)
	return result, err
}

func GetConfigCompare(ctx context.Context, service *zscaler.Service, gatewayID int) (*ZtGateway, error) {
	var result ZtGateway
	err := service.Client.Read(ctx, fmt.Sprintf("%s/%d", ztGatewayConfigCompareEndpoint, gatewayID), &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func GetConfigStatus(ctx context.Context, service *zscaler.Service, gatewayID int) (*ZtGateway, error) {
	var result ZtGateway
	err := service.Client.Read(ctx, fmt.Sprintf("%s/%d", ztGatewayConfigStatusEndpoint, gatewayID), &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func GetEvents(ctx context.Context, service *zscaler.Service, gatewayID int) (*ZtGateway, error) {
	var result ZtGateway
	err := service.Client.Read(ctx, fmt.Sprintf("%s/%d", ztGatewayEventsEndpoint, gatewayID), &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func GetStatus(ctx context.Context, service *zscaler.Service, gatewayID int) (*ZtGateway, error) {
	var result ZtGateway
	err := service.Client.Read(ctx, fmt.Sprintf("%s/%d", ztGatewayStatusEndpoint, gatewayID), &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func GetTrafficTestEnv(ctx context.Context, service *zscaler.Service, gatewayID int) (*ZtGateway, error) {
	var result ZtGateway
	err := service.Client.Read(ctx, fmt.Sprintf("%s/%d", ztGatewayTrafficTestEnvEndpoint, gatewayID), &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func GetTrafficTest(ctx context.Context, service *zscaler.Service, testID int) (*ZtGateway, error) {
	var result ZtGateway
	err := service.Client.Read(ctx, fmt.Sprintf("%s/%d", ztGatewayTrafficTestEndpoint, testID), &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func Get(ctx context.Context, service *zscaler.Service, gatewayID int) (*ZtGateway, error) {
	var result ZtGateway
	err := service.Client.Read(ctx, fmt.Sprintf("%s/%d", ztGatewayEndpoint, gatewayID), &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func Create(ctx context.Context, service *zscaler.Service, body *ZtGateway) (*ZtGateway, *http.Response, error) {
	resp, err := service.Client.Create(ctx, ztGatewayEndpoint, *body)
	if err != nil {
		return nil, nil, err
	}
	created, ok := resp.(*ZtGateway)
	if !ok {
		return nil, nil, errors.New("object returned from api was not a gateway pointer")
	}
	return created, nil, nil
}

func CreateTrafficTest(ctx context.Context, service *zscaler.Service, body *ZtGateway) (*ZtGateway, *http.Response, error) {
	resp, err := service.Client.Create(ctx, ztGatewayTrafficTestEndpoint, *body)
	if err != nil {
		return nil, nil, err
	}
	created, ok := resp.(*ZtGateway)
	if !ok {
		return nil, nil, errors.New("object returned from api was not a gateway pointer")
	}
	return created, nil, nil
}

func CreateTrafficTestEnv(ctx context.Context, service *zscaler.Service, gatewayID int, body *ZtGateway) (*ZtGateway, *http.Response, error) {
	resp, err := service.Client.Create(ctx, fmt.Sprintf("%s/%d", ztGatewayTrafficTestEnvEndpoint, gatewayID), *body)
	if err != nil {
		return nil, nil, err
	}
	created, ok := resp.(*ZtGateway)
	if !ok {
		return nil, nil, errors.New("object returned from api was not a gateway pointer")
	}
	return created, nil, nil
}

func CreateTrafficTestRun(ctx context.Context, service *zscaler.Service, gwID int, body *RunRequest) (*ZtGateway, *http.Response, error) {
	resp, err := service.Client.Create(ctx, fmt.Sprintf("%s/%d/trafficTest/run", ztGatewayEndpoint, gwID), *body)
	if err != nil {
		return nil, nil, err
	}
	created, ok := resp.(*ZtGateway)
	if !ok {
		return nil, nil, errors.New("object returned from api was not a gateway pointer")
	}
	return created, nil, nil
}

func UpdateAnalyticMetrics(ctx context.Context, service *zscaler.Service, gatewayID int, body *MetricsRequest) (*ZtGateway, *http.Response, error) {
	resp, err := service.Client.UpdateWithPut(ctx, fmt.Sprintf("%s/%d", ztGatewayAnalyticMetricsEndpoint, gatewayID), *body)
	if err != nil {
		return nil, nil, err
	}
	updated, _ := resp.(*ZtGateway)
	return updated, nil, nil
}

func UpdateTrafficTest(ctx context.Context, service *zscaler.Service, testID int, body *ZtGateway) (*ZtGateway, *http.Response, error) {
	resp, err := service.Client.UpdateWithPut(ctx, fmt.Sprintf("%s/%d", ztGatewayTrafficTestEndpoint, testID), *body)
	if err != nil {
		return nil, nil, err
	}
	updated, _ := resp.(*ZtGateway)
	return updated, nil, nil
}

func Update(ctx context.Context, service *zscaler.Service, gatewayID int, body *ZtGateway) (*ZtGateway, *http.Response, error) {
	resp, err := service.Client.UpdateWithPut(ctx, fmt.Sprintf("%s/%d", ztGatewayEndpoint, gatewayID), *body)
	if err != nil {
		return nil, nil, err
	}
	updated, _ := resp.(*ZtGateway)
	return updated, nil, nil
}

func DeleteTrafficTestEnv(ctx context.Context, service *zscaler.Service, gatewayID int) (*http.Response, error) {
	err := service.Client.Delete(ctx, fmt.Sprintf("%s/%d", ztGatewayTrafficTestEnvEndpoint, gatewayID))
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func DeleteTrafficTest(ctx context.Context, service *zscaler.Service, testID int) (*http.Response, error) {
	err := service.Client.Delete(ctx, fmt.Sprintf("%s/%d", ztGatewayTrafficTestEndpoint, testID))
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func Delete(ctx context.Context, service *zscaler.Service, gatewayID int) (*http.Response, error) {
	err := service.Client.Delete(ctx, fmt.Sprintf("%s/%d", ztGatewayEndpoint, gatewayID))
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func GetZtGatewayByName(ctx context.Context, service *zscaler.Service, name string) (*ZtGateway, error) {
	items, err := GetAll(ctx, service)
	if err != nil {
		return nil, err
	}
	for _, item := range items {
		if strings.EqualFold(item.Name, name) {
			return &item, nil
		}
	}
	return nil, fmt.Errorf("no gateway found with name: %s", name)
}
