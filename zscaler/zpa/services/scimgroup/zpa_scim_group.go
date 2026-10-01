package scimgroup

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
)

const (
	userConfig        = "/zpa/userconfig/v1/customers/"
	scimGroupEndpoint = "/scimgroup"
	idpIdPath         = "/idpId"
)

type ScimGroup struct {
	CreationTime int64  `json:"creationTime,omitempty"`
	ID           int64  `json:"id,omitempty"`
	IdpGroupID   string `json:"idpGroupId,omitempty"`
	IdpID        int64  `json:"idpId,omitempty"`
	IdpName      string `json:"idpName,omitempty"`
	ModifiedTime int64  `json:"modifiedTime,omitempty"`
	Name         string `json:"name,omitempty"`
	InternalID   string `json:"internalId,omitempty"`
	IamIdpID     string `json:"iamIdpId,omitempty"`
	IamIdpName   string `json:"iamIdpName,omitempty"`
}

// IamIdpCriteria narrows a SCIM group name lookup to the groups that belong to a
// specific ZIdentity (IAM) IdP. Several IAM IdPs can map to the same ZPA IdP, so
// a group name is only unique within an IAM IdP. Set at most one field.
type IamIdpCriteria struct {
	IamIdpID   string
	IamIdpName string
}

func Get(ctx context.Context, service *zscaler.Service, scimGroupID string) (*ScimGroup, *http.Response, error) {
	v := new(ScimGroup)
	relativeURL := fmt.Sprintf("%s/%s", userConfig+service.Client.GetCustomerID()+scimGroupEndpoint, scimGroupID)
	resp, err := service.Client.NewRequestDo(ctx, "GET", relativeURL, nil, nil, v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func GetByName(ctx context.Context, service *zscaler.Service, scimName, idpId string) (*ScimGroup, *http.Response, error) {
	// Construct the API endpoint URL with query parameters
	relativeURL := fmt.Sprintf("%s/%s", userConfig+service.Client.GetCustomerID()+scimGroupEndpoint+idpIdPath, idpId)
	// Fetch the pages
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[ScimGroup](ctx, service.Client, relativeURL, common.Filter{
		Search:    scimName,
		SortBy:    string(service.SortBy),
		SortOrder: string(service.SortOrder),
	})
	if err != nil {
		return nil, resp, err
	}

	// Look for the group with the specified name
	for _, scim := range list {
		if strings.EqualFold(scim.Name, scimName) {
			return &scim, resp, nil
		}
	}

	return nil, resp, fmt.Errorf("no SCIM group named '%s' was found", scimName)
}

// GetByNameAndIamIdp returns the SCIM group named scimName under the ZPA IdP
// idpId that belongs to the IAM IdP described by criteria. Unlike GetByName,
// which returns the first name match, it fails when more than one group
// matches, listing the candidates, so a name shared across IAM IdPs can never
// resolve to the wrong group. With an empty criteria it behaves like GetByName.
func GetByNameAndIamIdp(ctx context.Context, service *zscaler.Service, scimName, idpId string, criteria IamIdpCriteria) (*ScimGroup, *http.Response, error) {
	if criteria.IamIdpID == "" && criteria.IamIdpName == "" {
		return GetByName(ctx, service, scimName, idpId)
	}
	if criteria.IamIdpID != "" && criteria.IamIdpName != "" {
		return nil, nil, fmt.Errorf("only one of IamIdpID or IamIdpName can be set")
	}

	relativeURL := fmt.Sprintf("%s/%s", userConfig+service.Client.GetCustomerID()+scimGroupEndpoint+idpIdPath, idpId)
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[ScimGroup](ctx, service.Client, relativeURL, common.Filter{
		Search:    scimName,
		SortBy:    string(service.SortBy),
		SortOrder: string(service.SortOrder),
	})
	if err != nil {
		return nil, resp, err
	}

	var nameMatches, matches []ScimGroup
	for _, scim := range list {
		if !scimGroupNameMatches(scim, scimName, criteria) {
			continue
		}
		nameMatches = append(nameMatches, scim)
		if scimGroupInIamIdp(scim, criteria) {
			matches = append(matches, scim)
		}
	}

	switch len(matches) {
	case 1:
		return &matches[0], resp, nil
	case 0:
		if len(nameMatches) == 0 {
			return nil, resp, fmt.Errorf("no SCIM group named '%s' was found", scimName)
		}
		return nil, resp, fmt.Errorf("no SCIM group named '%s' was found for %s; groups with that name: %s",
			scimName, describeIamIdpCriteria(criteria), describeScimGroups(nameMatches))
	default:
		return nil, resp, fmt.Errorf("multiple SCIM groups named '%s' were found for %s: %s",
			scimName, describeIamIdpCriteria(criteria), describeScimGroups(matches))
	}
}

// scimGroupNameMatches compares the group name case-insensitively. Before the
// API returned iamIdpId/iamIdpName it decorated names as "<name>(<iamIdpName>)",
// so that suffix is stripped before comparing.
func scimGroupNameMatches(scim ScimGroup, scimName string, criteria IamIdpCriteria) bool {
	if strings.EqualFold(scim.Name, scimName) {
		return true
	}
	suffixName := scim.IamIdpName
	if suffixName == "" {
		suffixName = criteria.IamIdpName
	}
	if suffixName == "" {
		return false
	}
	return strings.EqualFold(stripIamIdpSuffix(scim.Name, suffixName), scimName)
}

// scimGroupInIamIdp reports whether the group belongs to the IAM IdP in
// criteria. A decorated name is accepted as evidence of the IAM IdP when the API
// does not return iamIdpName.
func scimGroupInIamIdp(scim ScimGroup, criteria IamIdpCriteria) bool {
	if criteria.IamIdpID != "" {
		return scim.IamIdpID == criteria.IamIdpID
	}
	if scim.IamIdpName != "" {
		return strings.EqualFold(scim.IamIdpName, criteria.IamIdpName)
	}
	return stripIamIdpSuffix(scim.Name, criteria.IamIdpName) != scim.Name
}

func stripIamIdpSuffix(name, iamIdpName string) string {
	suffix := "(" + iamIdpName + ")"
	if len(name) > len(suffix) && strings.EqualFold(name[len(name)-len(suffix):], suffix) {
		return strings.TrimSpace(name[:len(name)-len(suffix)])
	}
	return name
}

func describeIamIdpCriteria(criteria IamIdpCriteria) string {
	if criteria.IamIdpID != "" {
		return fmt.Sprintf("IAM IdP ID '%s'", criteria.IamIdpID)
	}
	return fmt.Sprintf("IAM IdP name '%s'", criteria.IamIdpName)
}

func describeScimGroups(groups []ScimGroup) string {
	parts := make([]string, 0, len(groups))
	for _, g := range groups {
		parts = append(parts, fmt.Sprintf("(id=%d, name=%q, iamIdpId=%q, iamIdpName=%q)", g.ID, g.Name, g.IamIdpID, g.IamIdpName))
	}
	return strings.Join(parts, ", ")
}

func GetAllByIdpId(ctx context.Context, service *zscaler.Service, idpId string) ([]ScimGroup, *http.Response, error) {
	relativeURL := fmt.Sprintf("%s/%s", userConfig+service.Client.GetCustomerID()+scimGroupEndpoint+idpIdPath, idpId)
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[ScimGroup](ctx, service.Client, relativeURL, common.Filter{
		SortBy:    string(service.SortBy),
		SortOrder: string(service.SortOrder),
	})
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}
