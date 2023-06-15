package cbibannercontroller

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/zpa/services/common"
)

const (
	cbiConfig         = "/cbiconfig/cbi/api/customers/"
	cbiBannerEndpoint = "/banners"
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
}

func (service *Service) Get(bannerID string) (*CBIBannerController, *http.Response, error) {
	v := new(CBIBannerController)
	relativeURL := fmt.Sprintf("%s/%s", cbiConfig+service.Client.Config.CustomerID+cbiBannerEndpoint, bannerID)
	resp, err := service.Client.NewRequestDo("GET", relativeURL, nil, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func (service *Service) GetByName(bannerName string) (*CBIBannerController, *http.Response, error) {
	relativeURL := cbiConfig + service.Client.Config.CustomerID + cbiBannerEndpoint
	list, resp, err := common.GetAllPagesGeneric[CBIBannerController](service.Client, relativeURL, "")
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

func (service *Service) Create(server CBIBannerController) (*CBIBannerController, *http.Response, error) {
	v := new(CBIBannerController)
	resp, err := service.Client.NewRequestDo("POST", cbiConfig+service.Client.Config.CustomerID+cbiBannerEndpoint, nil, server, &v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func (service *Service) Update(id string, cbiBanner CBIBannerController) (*http.Response, error) {
	path := fmt.Sprintf("%s/%s", cbiConfig+service.Client.Config.CustomerID+cbiBannerEndpoint, id)
	resp, err := service.Client.NewRequestDo("PUT", path, nil, cbiBanner, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func (service *Service) Delete(id string) (*http.Response, error) {
	path := fmt.Sprintf("%s/%s", cbiConfig+service.Client.Config.CustomerID+cbiBannerEndpoint, id)
	resp, err := service.Client.NewRequestDo("DELETE", path, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func (service *Service) GetAll() ([]CBIBannerController, *http.Response, error) {
	relativeURL := cbiConfig + service.Client.Config.CustomerID + cbiBannerEndpoint
	list, resp, err := common.GetAllPagesGeneric[CBIBannerController](service.Client, relativeURL, "")
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}
