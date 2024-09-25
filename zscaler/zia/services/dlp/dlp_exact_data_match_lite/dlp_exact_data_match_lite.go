package dlp_exact_data_match_lite

import (
	"context"
	"fmt"
	"net/url"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/common"
)

const (
	dlpEDMELiteEndpoint = "/zia/api/v1/dlpExactDataMatchSchemas/lite"
)

// Gets a list of active EDM templates (or EDM profiles) and their criteria, only.
type DLPEDMLite struct {
	// The identifier (1-65519) for the EDM schema (i.e., EDM template) that is unique within the organization.
	Schema SchemaIDNameExtension `json:"schema,omitempty"`

	// Indicates the status of a specified EDM schema (i.e., EDM template). If this value is set to true, the schema is active and can be used by DLP dictionaries.
	TokenList []TokenList `json:"tokenList,omitempty"`
}

type SchemaIDNameExtension struct {
	// Identifier that uniquely identifies an entity
	ID int `json:"ID,omitempty"`

	// The configured name of the entity
	Name string `json:"name,omitempty"`

	// An external identifier used for an entity that is managed outside of ZIA.
	// Examples include zpaServerGroup and zpaAppSegments.
	// This field is not applicable to ZIA-managed entities.
	ExternalID string `json:"externalId,omitempty"`

	Extensions map[string]interface{} `json:"extensions,omitempty"`
}

type TokenList struct {
	// The token (i.e., criteria) name. This attribute is required by PUT and POST requests.
	Name string `json:"name,omitempty"`

	// The token (i.e., criteria) name. This attribute is required by PUT and POST requests.
	Type string `json:"type,omitempty"`

	// Indicates whether the token is a primary key.
	PrimaryKey bool `json:"primaryKey,omitempty"`

	// The column position for the token in the original CSV file uploaded to the Index Tool, starting from 1. This attribue required by PUT and POST requests.
	OriginalColumn int `json:"originalColumn,omitempty"`

	// The column position for the token in the hashed file, starting from 1.
	HashfileColumnOrder int `json:"hashfileColumnOrder,omitempty"`

	// The length of the column bitmap in the hashed file.
	ColLengthBitmap int `json:"colLengthBitmap,omitempty"`
}

func GetBySchemaName(ctx context.Context, service *zscaler.Service, schemaName string, activeOnly, fetchTokens bool) ([]DLPEDMLite, error) {
	queryParameters := url.Values{}
	queryParameters.Set("schemaName", schemaName)
	if activeOnly {
		queryParameters.Set("activeOnly", "true")
	}
	if fetchTokens {
		queryParameters.Set("fetchTokens", "true")
	}

	endpoint := fmt.Sprintf("%s?%s", dlpEDMELiteEndpoint, queryParameters.Encode())
	var edmData []DLPEDMLite
	err := common.ReadAllPages(ctx, service.Client, endpoint, &edmData)
	if err != nil {
		return nil, err
	}
	return edmData, nil
}

func GetAllEDMSchema(ctx context.Context, service *zscaler.Service, activeOnly, fetchTokens bool) ([]DLPEDMLite, error) {
	queryParameters := url.Values{}
	if activeOnly {
		queryParameters.Set("activeOnly", "true")
	}
	if fetchTokens {
		queryParameters.Set("fetchTokens", "true")
	}

	endpoint := dlpEDMELiteEndpoint
	if len(queryParameters) > 0 {
		endpoint += "?" + queryParameters.Encode()
	}

	var edmData []DLPEDMLite
	err := common.ReadAllPages(ctx, service.Client, endpoint, &edmData)
	return edmData, err
}
