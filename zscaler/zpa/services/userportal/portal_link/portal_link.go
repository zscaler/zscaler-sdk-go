package portal_link

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/userportal/portal_controller"
)

const (
	mgmtConfig             = "/zpa/mgmtconfig/v1/admin/customers/"
	mgmtConfigV2           = "/zpa/mgmtconfig/v2/admin/customers/"
	userPortalLinkEndpoint = "/userPortalLink"
)

type UserPortalLink struct {
	ApplicationID   string                                   `json:"applicationId,omitempty"`
	CreationTime    string                                   `json:"creationTime,omitempty"`
	Description     string                                   `json:"description,omitempty"`
	Enabled         bool                                     `json:"enabled,omitempty"`
	IconText        string                                   `json:"iconText,omitempty"`
	ID              string                                   `json:"id,omitempty"`
	Link            string                                   `json:"link,omitempty"`
	LinkPath        string                                   `json:"linkPath,omitempty"`
	ModifiedBy      string                                   `json:"modifiedBy,omitempty"`
	ModifiedTime    string                                   `json:"modifiedTime,omitempty"`
	Name            string                                   `json:"name,omitempty"`
	Protocol        string                                   `json:"protocol,omitempty"`
	MicrotenantID   string                                   `json:"microtenantId,omitempty"`
	MicrotenantName string                                   `json:"microtenantName,omitempty"`
	NameWithoutTrim string                                   `json:"nameWithoutTrim,omitempty"`
	UserPortalID    string                                   `json:"userPortalId,omitempty"`
	UserPortals     []portal_controller.UserPortalController `json:"userPortals"`
}

type PortalLinkBulkRequest struct {
	UserPortalLinks []UserPortalLink                         `json:"userPortalLinks"`
	UserPortals     []portal_controller.UserPortalController `json:"userPortals"`
}

func Get(ctx context.Context, service *zscaler.Service, portalID string) (*UserPortalLink, *http.Response, error) {
	v := new(UserPortalLink)
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.GetCustomerID()+userPortalLinkEndpoint, portalID)
	resp, err := service.Client.NewRequestDo(ctx, "GET", relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()}, nil, v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func GetByName(ctx context.Context, service *zscaler.Service, portalName string) (*UserPortalLink, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + userPortalLinkEndpoint
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[UserPortalLink](ctx, service.Client, relativeURL, common.Filter{Search: portalName, MicroTenantID: service.MicroTenantID()})
	if err != nil {
		return nil, nil, err
	}
	for _, app := range list {
		if strings.EqualFold(app.Name, portalName) {
			return &app, resp, nil
		}
	}
	return nil, resp, fmt.Errorf("no user portal link named '%s' was found", portalName)
}

func GetUserPortalLinks(ctx context.Context, service *zscaler.Service, portalID string) (*UserPortalLink, *http.Response, error) {
	v := new(UserPortalLink)
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.GetCustomerID()+userPortalLinkEndpoint+"/userPortal", portalID)
	resp, err := service.Client.NewRequestDo(ctx, "GET", relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()}, nil, v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func Create(ctx context.Context, service *zscaler.Service, controllerGroup UserPortalLink) (*UserPortalLink, *http.Response, error) {
	v := new(UserPortalLink)
	resp, err := service.Client.NewRequestDo(ctx, "POST", mgmtConfig+service.Client.GetCustomerID()+userPortalLinkEndpoint, common.Filter{MicroTenantID: service.MicroTenantID()}, controllerGroup, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func CreatePortalLinkBulk(ctx context.Context, service *zscaler.Service, portalBulk []UserPortalLink) ([]UserPortalLink, *http.Response, error) {
	var responseBody PortalLinkBulkRequest

	requestBody := PortalLinkBulkRequest{
		UserPortalLinks: portalBulk,
	}

	// Optionally extract userPortalId values and include them in the outer object
	var userPortals []portal_controller.UserPortalController
	for _, link := range portalBulk {
		if link.UserPortalID != "" {
			userPortals = append(userPortals, portal_controller.UserPortalController{ID: link.UserPortalID})
		}
	}
	if len(userPortals) > 0 {
		requestBody.UserPortals = userPortals
	}

	relativeURL := mgmtConfigV2 + service.Client.GetCustomerID() + userPortalLinkEndpoint + "/bulk"
	resp, err := service.Client.NewRequestDo(
		ctx,
		"POST",
		relativeURL,
		common.Filter{MicroTenantID: service.MicroTenantID()},
		requestBody,
		&responseBody,
	)
	if err != nil {
		return nil, nil, err
	}

	return responseBody.UserPortalLinks, resp, nil
}

func Update(ctx context.Context, service *zscaler.Service, portalID string, controllerGroup *UserPortalLink) (*http.Response, error) {
	path := fmt.Sprintf("%v/%v", mgmtConfig+service.Client.GetCustomerID()+userPortalLinkEndpoint, portalID)
	resp, err := service.Client.NewRequestDo(ctx, "PUT", path, common.Filter{MicroTenantID: service.MicroTenantID()}, controllerGroup, nil)
	if err != nil {
		return nil, err
	}

	return resp, err
}

func Delete(ctx context.Context, service *zscaler.Service, portalID string) (*http.Response, error) {
	path := fmt.Sprintf("%v/%v", mgmtConfig+service.Client.GetCustomerID()+userPortalLinkEndpoint, portalID)
	resp, err := service.Client.NewRequestDo(ctx, "DELETE", path, common.Filter{MicroTenantID: service.MicroTenantID()}, nil, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func GetAll(ctx context.Context, service *zscaler.Service) ([]UserPortalLink, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + userPortalLinkEndpoint
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[UserPortalLink](ctx, service.Client, relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()})
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}
