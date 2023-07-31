package cbibannercontroller

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/zpa/services/common"
)

const (
	cbiConfig          = "/cbiconfig/cbi/api/customers/"
	cbiBannerEndpoint  = "/banner"
	cbiBannersEndpoint = "/banners"
)

type CBIBannerController struct {
	ID                string `json:"id,omitempty"`
	Name              string `json:"name,omitempty"`
	PrimaryColor      string `json:"primaryColor,omitempty"`
	TextColor         string `json:"textColor,omitempty"`
	NotificationTitle string `json:"notificationTitle,omitempty"`
	NotificationText  string `json:"notificationText,omitempty"`
	Logo              string `json:"logo,omitempty"`
	Banner            bool   `json:"banner,omitempty"`
	IsDefault         bool   `json:"isDefault,omitempty"`
	Persist           bool   `json:"persist,omitempty"`
}

func (service *Service) Get(bannerID string) (*CBIBannerController, *http.Response, error) {
	v := new(CBIBannerController)
	relativeURL := fmt.Sprintf("%s/%s", cbiConfig+service.Client.Config.CustomerID+cbiBannersEndpoint, bannerID)
	resp, err := service.Client.NewRequestDo("GET", relativeURL, common.Filter{MicroTenantID: service.microTenantID}, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func (service *Service) GetByName(bannerName string) (*CBIBannerController, *http.Response, error) {
	list, resp, err := service.GetAll()
	if err != nil {
		return nil, nil, err
	}
	for _, banner := range list {
		if strings.EqualFold(banner.Name, bannerName) {
			return &banner, resp, nil
		}
	}
	return nil, resp, fmt.Errorf("no cloud browser isolation banner named '%s' was found", bannerName)
}

func (service *Service) Create(cbiBanner *CBIBannerController) (*CBIBannerController, *http.Response, error) {
	v := new(CBIBannerController)
	resp, err := service.Client.NewRequestDo("POST", cbiConfig+service.Client.Config.CustomerID+cbiBannerEndpoint, common.Filter{MicroTenantID: service.microTenantID}, cbiBanner, &v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func (service *Service) Update(cbiBannerID string, cbiBanner *CBIBannerController) (*http.Response, error) {
	path := fmt.Sprintf("%v/%v", cbiConfig+service.Client.Config.CustomerID+cbiBannersEndpoint, cbiBannerID)
	resp, err := service.Client.NewRequestDo("PUT", path, common.Filter{MicroTenantID: service.microTenantID}, cbiBanner, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func (service *Service) Delete(cbiBannerID string) (*http.Response, error) {
	path := fmt.Sprintf("%v/%v", cbiConfig+service.Client.Config.CustomerID+cbiBannersEndpoint, cbiBannerID)
	resp, err := service.Client.NewRequestDo("DELETE", path, common.Filter{MicroTenantID: service.microTenantID}, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func (service *Service) GetAll() ([]CBIBannerController, *http.Response, error) {
	relativeURL := cbiConfig + service.Client.Config.CustomerID + cbiBannersEndpoint
	var list []CBIBannerController
	resp, err := service.Client.NewRequestDo("GET", relativeURL, common.Filter{MicroTenantID: service.microTenantID}, nil, &list)
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}
