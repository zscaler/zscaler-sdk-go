package bandwidth_classes

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/common"
)

const (
	bandwidthClassEndpoint = "/zia/api/v1/bandwidthClasses"
)

type BandwidthClasses struct {
	ID                       int      `json:"id,omitempty"`
	IsNameL10nTag            bool     `json:"isNameL10nTag,omitempty"`
	Name                     string   `json:"name,omitempty"`
	GetfileSize              string   `json:"getfileSize,omitempty"`
	FileSize                 string   `json:"fileSize,omitempty"`
	Type                     string   `json:"type,omitempty"`
	WebApplications          []string `json:"webApplications,omitempty"`
	Urls                     []string `json:"urls,omitempty"`
	ApplicationServiceGroups []string `json:"applicationServiceGroups,omitempty"`
	NetworkApplications      []string `json:"networkApplications,omitempty"`
	NetworkServices          []string `json:"networkServices,omitempty"`
	UrlCategories            []string `json:"urlCategories,omitempty"`
	Applications             []string `json:"applications,omitempty"`
}

func Get(ctx context.Context, service *zscaler.Service, classID int) (*BandwidthClasses, error) {
	var class BandwidthClasses
	err := service.Client.Read(ctx, fmt.Sprintf("%s/%d", bandwidthClassEndpoint, classID), &class)
	if err != nil {
		return nil, err
	}

	service.Client.GetLogger().Printf("[DEBUG] Returning bandwidth class from Get: %d", class.ID)
	return &class, nil
}

func GetByName(ctx context.Context, service *zscaler.Service, className string) (*BandwidthClasses, error) {
	var bdwClasses []BandwidthClasses
	err := common.ReadAllPages(ctx, service.Client, bandwidthClassEndpoint, &bdwClasses)
	if err != nil {
		return nil, err
	}
	for _, bdwClass := range bdwClasses {
		if strings.EqualFold(bdwClass.Name, className) {
			return &bdwClass, nil
		}
	}
	return nil, fmt.Errorf("no bandwidth classes found with name: %s", className)
}

func Create(ctx context.Context, service *zscaler.Service, classID *BandwidthClasses) (*BandwidthClasses, *http.Response, error) {
	resp, err := service.Client.Create(ctx, bandwidthClassEndpoint, *classID)
	if err != nil {
		return nil, nil, err
	}

	createdClass, ok := resp.(*BandwidthClasses)
	if !ok {
		return nil, nil, errors.New("object returned from api was not a bandwidth classes pointer")
	}

	service.Client.GetLogger().Printf("[DEBUG]returning new bandwidth classes from create: %d", createdClass.ID)
	return createdClass, nil, nil
}

func Update(ctx context.Context, service *zscaler.Service, classID int, classes *BandwidthClasses) (*BandwidthClasses, *http.Response, error) {
	resp, err := service.Client.UpdateWithPut(ctx, fmt.Sprintf("%s/%d", bandwidthClassEndpoint, classID), *classes)
	if err != nil {
		return nil, nil, err
	}
	updatedClass, _ := resp.(*BandwidthClasses)

	service.Client.GetLogger().Printf("[DEBUG]returning updates bandwidth classes from update: %d", updatedClass.ID)
	return updatedClass, nil, nil
}

func Delete(ctx context.Context, service *zscaler.Service, classID int) (*http.Response, error) {
	err := service.Client.Delete(ctx, fmt.Sprintf("%s/%d", bandwidthClassEndpoint, classID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func GetAllLite(ctx context.Context, service *zscaler.Service) ([]BandwidthClasses, error) {
	var classes []BandwidthClasses
	err := common.ReadAllPages(ctx, service.Client, bandwidthClassEndpoint+"/lite", &classes)
	return classes, err
}

func GetAll(ctx context.Context, service *zscaler.Service) ([]BandwidthClasses, error) {
	var classes []BandwidthClasses
	err := common.ReadAllPages(ctx, service.Client, bandwidthClassEndpoint, &classes)
	return classes, err
}
