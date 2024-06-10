package scimattributeheader

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/common"
)

const (
	mgmtConfig       = "/mgmtconfig/v1/admin/customers/"
	userConfig       = "/userconfig/v1/customers"
	idpId            = "/idp"
	scimAttrEndpoint = "/scimattribute"
)

type ScimAttributeHeader struct {
	CanonicalValues []string `json:"canonicalValues,omitempty"`
	CaseSensitive   bool     `json:"caseSensitive,omitempty"`
	CreationTime    string   `json:"creationTime,omitempty,"`
	DataType        string   `json:"dataType,omitempty"`
	Description     string   `json:"description,omitempty"`
	ID              string   `json:"id,omitempty"`
	IdpID           string   `json:"idpId,omitempty"`
	ModifiedBy      string   `json:"modifiedBy,omitempty"`
	ModifiedTime    string   `json:"modifiedTime,omitempty"`
	MultiValued     bool     `json:"multivalued,omitempty"`
	Mutability      string   `json:"mutability,omitempty"`
	Name            string   `json:"name,omitempty"`
	Required        bool     `json:"required,omitempty"`
	Returned        string   `json:"returned,omitempty"`
	SchemaURI       string   `json:"schemaURI,omitempty"`
	Uniqueness      bool     `json:"uniqueness,omitempty"`
}

func Get(service *services.Service, idpId, scimAttrHeaderID string) (*ScimAttributeHeader, *http.Response, error) {
	v := new(ScimAttributeHeader)
	relativeURL := fmt.Sprintf("%s/idp/%s/scimattribute/%s", mgmtConfig+service.Client.Config.CustomerID, idpId, scimAttrHeaderID)
	resp, err := service.Client.NewRequestDo("GET", relativeURL, nil, nil, v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

// SearchValues searchs by features and fields for the API.
func SearchValues(service *services.Service, idpId, ScimAttrHeaderID, searchQuery string) ([]string, error) {
	searchQuery = strings.Split(searchQuery, "@")[0]
	relativeURL := fmt.Sprintf("%s/%s/scimattribute/idpId/%s/attributeId/%s", userConfig, service.Client.Config.CustomerID, idpId, ScimAttrHeaderID)
	l, _, err := common.GetAllPagesGeneric[string](service.Client, relativeURL, searchQuery)
	return l, err
}

func GetValues(service *services.Service, idpId, ScimAttrHeaderID string) ([]string, error) {
	relativeURL := fmt.Sprintf("%s/%s/scimattribute/idpId/%s/attributeId/%s", userConfig, service.Client.Config.CustomerID, idpId, ScimAttrHeaderID)
	l, _, err := common.GetAllPagesGeneric[string](service.Client, relativeURL, "")
	return l, err
}

func GetByName(service *services.Service, scimAttributeName, IdpId string) (*ScimAttributeHeader, *http.Response, error) {
	relativeURL := fmt.Sprintf("%s/%s%s", mgmtConfig+service.Client.Config.CustomerID+idpId, IdpId, scimAttrEndpoint)
	list, resp, err := common.GetAllPagesGeneric[ScimAttributeHeader](service.Client, relativeURL, "")
	if err != nil {
		return nil, nil, err
	}
	for _, scimAttribute := range list {
		if strings.EqualFold(scimAttribute.Name, scimAttributeName) {
			return &scimAttribute, resp, nil
		}
	}
	return nil, resp, fmt.Errorf("no scim named '%s' was found", scimAttributeName)
}

func GetAllByIdpId(service *services.Service, IdpId string) ([]ScimAttributeHeader, *http.Response, error) {
	relativeURL := fmt.Sprintf("%s/%s%s", mgmtConfig+service.Client.Config.CustomerID+idpId, IdpId, scimAttrEndpoint)
	list, resp, err := common.GetAllPagesGeneric[ScimAttributeHeader](service.Client, relativeURL, "")
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}
