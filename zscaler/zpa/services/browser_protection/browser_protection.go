package browser_protection

import (
	"context"
	"fmt"
	"net/http"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
)

const (
	mgmtConfig                = "/zpa/mgmtconfig/v1/admin/customers/"
	browserProtectionEndpoint = "/browserProtectionProfile"
)

type BrowserProtection struct {
	CreationTime      string   `json:"creationTime,omitempty"`
	Criteria          Criteria `json:"criteria,omitempty"`
	CriteriaFlagsMask string   `json:"criteriaFlagsMask,omitempty"`
	DefaultCSP        bool     `json:"defaultCSP,omitempty"`
	Description       string   `json:"description,omitempty"`
	ID                string   `json:"id,omitempty"`
	ModifiedBy        string   `json:"modifiedBy,omitempty"`
	ModifiedTime      string   `json:"modifiedTime,omitempty"`
	Name              string   `json:"name,omitempty"`
}

type Criteria struct {
	FingerPrintCriteria FingerPrintCriteria `json:"fingerPrintCriteria,omitempty"`
}

type FingerPrintCriteria struct {
	Browser            BrowserCriteria  `json:"browser,omitempty"`
	CollectLocation    bool             `json:"collect_location,omitempty"`
	FingerprintTimeout string           `json:"fingerprint_timeout,omitempty"`
	Location           LocationCriteria `json:"location,omitempty"`
	System             SystemCriteria   `json:"system,omitempty"`
}

type BrowserCriteria struct {
	BrowserEng     bool `json:"browser_eng,omitempty"`
	BrowserEngVer  bool `json:"browser_eng_ver,omitempty"`
	BrowserName    bool `json:"browser_name,omitempty"`
	BrowserVersion bool `json:"browser_version,omitempty"`
	Canvas         bool `json:"canvas,omitempty"`
	FlashVer       bool `json:"flash_ver,omitempty"`
	FpUsrAgentStr  bool `json:"fp_usr_agent_str,omitempty"`
	IsCookie       bool `json:"is_cookie,omitempty"`
	IsLocalStorage bool `json:"is_local_storage,omitempty"`
	IsSessStorage  bool `json:"is_sess_storage,omitempty"`
	Ja3            bool `json:"ja3,omitempty"`
	Mime           bool `json:"mime,omitempty"`
	Plugin         bool `json:"plugin,omitempty"`
	SilverlightVer bool `json:"silverlight_ver,omitempty"`
}

type LocationCriteria struct {
	Lat bool `json:"lat,omitempty"`
	Lon bool `json:"lon,omitempty"`
}

type SystemCriteria struct {
	AvailScreenResolution bool `json:"avail_screen_resolution,omitempty"`
	CPUArch               bool `json:"cpu_arch,omitempty"`
	CurrScreenResolution  bool `json:"curr_screen_resolution,omitempty"`
	Font                  bool `json:"font,omitempty"`
	JavaVer               bool `json:"java_ver,omitempty"`
	MobileDevType         bool `json:"mobile_dev_type,omitempty"`
	MonitorMobile         bool `json:"monitor_mobile,omitempty"`
	OSName                bool `json:"os_name,omitempty"`
	OSVersion             bool `json:"os_version,omitempty"`
	SysLang               bool `json:"sys_lang,omitempty"`
	Tz                    bool `json:"tz,omitempty"`
	UsrLang               bool `json:"usr_lang,omitempty"`
}

func GetActiveBrowserProtectionProfile(ctx context.Context, service *zscaler.Service) ([]BrowserProtection, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + "/activeBrowserProtectionProfile"
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[BrowserProtection](ctx, service.Client, relativeURL, common.Filter{})
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}

func GetBrowserProtectionProfile(ctx context.Context, service *zscaler.Service) ([]BrowserProtection, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + "/browserProtectionProfile"
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[BrowserProtection](ctx, service.Client, relativeURL, common.Filter{})
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}

func UpdateBrowserProtectionProfile(ctx context.Context, service *zscaler.Service, profileID string) (*http.Response, error) {
	path := fmt.Sprintf("%s%s%s/setActive/%s", mgmtConfig, service.Client.GetCustomerID(), browserProtectionEndpoint, profileID)
	resp, err := service.Client.NewRequestDo(ctx, "PUT", path, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
