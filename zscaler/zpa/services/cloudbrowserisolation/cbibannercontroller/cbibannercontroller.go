package cbibannercontroller

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
)

const (
	cbiConfig          = "/zpa/cbiconfig/cbi/api/customers/"
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

func Get(ctx context.Context, service *zscaler.Service, bannerID string) (*CBIBannerController, *http.Response, error) {
	v := new(CBIBannerController)
	relativeURL := fmt.Sprintf("%s/%s", cbiConfig+service.Client.GetCustomerID()+cbiBannersEndpoint, bannerID)
	resp, err := service.Client.NewRequestDo(ctx, "GET", relativeURL, nil, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func GetByNameOrID(ctx context.Context, service *zscaler.Service, identifier string) (*CBIBannerController, *http.Response, error) {
	// Retrieve all banners
	list, resp, err := GetAll(ctx, service)
	if err != nil {
		return nil, nil, err
	}

	// Try to find by ID
	for _, banner := range list {
		if banner.ID == identifier {
			return Get(ctx, service, banner.ID)
		}
	}

	// Try to find by name
	for _, banner := range list {
		if strings.EqualFold(banner.Name, identifier) {
			return Get(ctx, service, banner.ID)
		}
	}

	return nil, resp, fmt.Errorf("no isolation banner named or with ID '%s' was found", identifier)
}

func Create(ctx context.Context, service *zscaler.Service, cbiBanner *CBIBannerController) (*CBIBannerController, *http.Response, error) {
	v := new(CBIBannerController)
	resp, err := service.Client.NewRequestDo(ctx, "POST", cbiConfig+service.Client.GetCustomerID()+cbiBannerEndpoint, nil, cbiBanner, &v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func Update(ctx context.Context, service *zscaler.Service, cbiBannerID string, cbiBanner *CBIBannerController) (*http.Response, error) {
	path := fmt.Sprintf("%v/%v", cbiConfig+service.Client.GetCustomerID()+cbiBannersEndpoint, cbiBannerID)
	resp, err := service.Client.NewRequestDo(ctx, "PUT", path, nil, cbiBanner, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func Delete(ctx context.Context, service *zscaler.Service, cbiBannerID string) (*http.Response, error) {
	path := fmt.Sprintf("%v/%v", cbiConfig+service.Client.GetCustomerID()+cbiBannersEndpoint, cbiBannerID)
	resp, err := service.Client.NewRequestDo(ctx, "DELETE", path, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func GetAll(ctx context.Context, service *zscaler.Service) ([]CBIBannerController, *http.Response, error) {
	relativeURL := cbiConfig + service.Client.GetCustomerID() + cbiBannersEndpoint
	var list []CBIBannerController
	resp, err := service.Client.NewRequestDo(ctx, "GET", relativeURL, nil, nil, &list)
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}
