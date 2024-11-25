package dlpdictionaries

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler"
	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler/zia/services/common"
)

const (
	dlpDictionariesEndpoint          = "/zia/api/v1/dlpDictionaries"
	dlpPredefinedIdentifiersEndpoint = "/zia/api/v1/predefinedIdentifiers"
)

type DlpDictionary struct {
	// Unique identifier for the DLP dictionary
	ID int `json:"id"`

	// The DLP dictionary's name
	Name string `json:"name,omitempty"`

	// The description of the DLP dictionary
	Description string `json:"description,omitempty"`

	// The DLP confidence threshold
	ConfidenceThreshold string `json:"confidenceThreshold,omitempty"`

	// The DLP custom phrase match type
	CustomPhraseMatchType string `json:"customPhraseMatchType,omitempty"`

	// Indicates whether the name is localized or not. This is always set to True for predefined DLP dictionaries.
	NameL10nTag bool `json:"nameL10nTag"`

	// This value is set to true for custom DLP dictionaries.
	Custom bool `json:"custom,omitempty"`

	// DLP threshold type
	ThresholdType string `json:"thresholdType,omitempty"`

	// The DLP dictionary type
	DictionaryType string `json:"dictionaryType,omitempty"`

	// The DLP dictionary proximity length.
	Proximity int `json:"proximity,omitempty"`

	// List containing the phrases used within a custom DLP dictionary. This attribute is not applicable to predefined DLP dictionaries.
	Phrases []Phrases `json:"phrases"`

	// List containing the patterns used within a custom DLP dictionary. This attribute is not applicable to predefined DLP dictionaries
	Patterns []Patterns `json:"patterns"`

	// Exact Data Match (EDM) related information for custom DLP dictionaries.
	EDMMatchDetails []EDMMatchDetails `json:"exactDataMatchDetails"`

	// List of Indexed Document Match (IDM) profiles and their corresponding match accuracy for custom DLP dictionaries.
	IDMProfileMatchAccuracy []IDMProfileMatchAccuracy `json:"idmProfileMatchAccuracyDetails"`

	// Indicates whether to exclude documents that are a 100% match to already-indexed documents from triggering an Indexed Document Match (IDM) Dictionary.
	IgnoreExactMatchIdmDict bool `json:"ignoreExactMatchIdmDict,omitempty"`

	// A true value denotes that the specified Bank Identification Number (BIN) values are included in the Credit Cards dictionary. A false value denotes that the specified BIN values are excluded from the Credit Cards dictionary.
	// Note: This field is applicable only to the predefined Credit Cards dictionary and its clones.
	IncludeBinNumbers bool `json:"includeBinNumbers,omitempty"`

	// The list of Bank Identification Number (BIN) values that are included or excluded from the Credit Cards dictionary. BIN values can be specified only for Diners Club, Mastercard, RuPay, and Visa cards. Up to 512 BIN values can be configured in a dictionary.
	// Note: This field is applicable only to the predefined Credit Cards dictionary and its clones.
	BinNumbers []int `json:"binNumbers,omitempty"`

	// ID of the predefined dictionary (original source dictionary) that is used for cloning. This field is applicable only to cloned dictionaries. Only a limited set of identification-based predefined dictionaries (e.g., Credit Cards, Social Security Numbers, National Identification Numbers, etc.) can be cloned. Up to 4 clones can be created from a predefined dictionary.
	DictTemplateId int `json:"dictTemplateId,omitempty"`

	// This field is set to true if the dictionary is cloned from a predefined dictionary. Otherwise, it is set to false.
	PredefinedClone bool `json:"predefinedClone,omitempty"`

	// This field specifies whether duplicate matches of a phrase from a dictionary must be counted individually toward the match count or ignored, thereby maintaining a single count for multiple occurrences.
	PredefinedCountActionType string `json:"predefinedCountActionType,omitempty"`

	// This value is set to true if proximity length and high confidence phrases are enabled for the DLP dictionary.
	ProximityLengthEnabled bool `json:"proximityLengthEnabled,omitempty"`

	// A Boolean constant that indicates if proximity length is enabled or disabled for a custom DLP dictionary. A true value indicates that proximity length is enabled, whereas a false value indicates that it is disabled.
	ProximityEnabledForCustomDictionary bool `json:"proximityEnabledForCustomDictionary,omitempty"`

	// A Boolean constant that indicates that the cloning option is supported for the DLP dictionary using the true value. This field is applicable only to predefined DLP dictionaries.
	DictionaryCloningEnabled bool `json:"dictionaryCloningEnabled"`

	// A Boolean constant that indicates that custom phrases are supported for the DLP dictionary using the true value. This field is applicable only to predefined DLP dictionaries with a high confidence score threshold.
	CustomPhraseSupported bool `json:"customPhraseSupported,omitempty"`

	// A true value indicates that the DLP dictionary is of hierarchical type that includes sub-dictionaries. A false value indicates that the dictionary is not hierarchical.
	HierarchicalDictionary bool `json:"hierarchicalDictionary,omitempty"`

	// The list of identifiers selected within a DLP dictionary of hierarchical type. Each identifier represents a sub-dictionary that consists of specific patterns. To retrieve the list of identifiers that are available for selection within a specific hierarchical dictionary, send a GET request to /dlpDictionaries/{dictId}/predefinedIdentifiers.
	HierarchicalIdentifiers []string `json:"hierarchicalIdentifiers,omitempty"`

	PredefinedPhrases []string `json:"predefinedPhrases,omitempty"`

	ThresholdAllowed bool `json:"thresholdAllowed,omitempty"`

	ConfidenceLevelForPredefinedDict string `json:"confidenceLevelForPredefinedDict"`
}

type Phrases struct {
	// The action applied to a DLP dictionary using phrases
	Action string `json:"action"`

	// DLP dictionary phrase
	Phrase string `json:"phrase"`
}

type Patterns struct {
	// The action applied to a DLP dictionary using patterns
	Action string `json:"action"`

	// DLP dictionary pattern
	Pattern string `json:"pattern"`
}

type EDMMatchDetails struct {
	// The unique identifier for the EDM mapping.
	DictionaryEdmMappingID int `json:"dictionaryEdmMappingId,omitempty"`

	// The unique identifier for the EDM template (or schema).
	SchemaID int `json:"schemaId,omitempty"`

	// The EDM template's primary field.
	PrimaryField int `json:"primaryField,omitempty"`

	// The EDM template's secondary fields.
	SecondaryFields []int `json:"secondaryFields,omitempty"`

	// The EDM secondary field to match on.
	SecondaryFieldMatchOn string `json:"secondaryFieldMatchOn,omitempty"`
}

type IDMProfileMatchAccuracy struct {
	// The IDM template reference.
	AdpIdmProfile *common.IDNameExtensions `json:"adpIdmProfile,omitempty"`

	// The IDM template match accuracy.
	MatchAccuracy string `json:"matchAccuracy,omitempty"`
}

func Get(ctx context.Context, service *zscaler.Service, dlpDictionariesID int) (*DlpDictionary, error) {
	var dlpDictionary DlpDictionary
	err := service.Client.Read(ctx, fmt.Sprintf("%s/%d", dlpDictionariesEndpoint, dlpDictionariesID), &dlpDictionary)
	if err != nil {
		return nil, err
	}

	service.Client.GetLogger().Printf("[DEBUG]Returning dictionary from Get: %d", dlpDictionary.ID)
	return &dlpDictionary, nil
}

func GetByName(ctx context.Context, service *zscaler.Service, dictionaryName string) (*DlpDictionary, error) {
	var dictionaries []DlpDictionary
	err := common.ReadAllPages(ctx, service.Client, dlpDictionariesEndpoint, &dictionaries)
	if err != nil {
		return nil, err
	}
	for _, dictionary := range dictionaries {
		if strings.EqualFold(dictionary.Name, dictionaryName) {
			return &dictionary, nil
		}
	}
	return nil, fmt.Errorf("no dictionary found with name: %s", dictionaryName)
}

// func GetPredefinedIdentifiers(ctx context.Context, service *zscaler.Service, dictionaryName string) ([]string, error) {
// 	// Use the GetByName function to retrieve the dictionary by name
// 	dictionary, err := GetByName(ctx,service, dictionaryName)
// 	if err != nil {
// 		return nil, fmt.Errorf("error retrieving dictionary by name: %v", err)
// 	}

// 	// If dictionary is found, get the predefined identifiers using the dictionary ID
// 	predefinedIdentifiersEndpoint := fmt.Sprintf(dlpDictionariesEndpoint+"/%d/predefinedIdentifiers", dictionary.ID)
// 	var predefinedIdentifiers []string
// 	err = service.Client.Read(ctx, predefinedIdentifiersEndpoint, &predefinedIdentifiers)
// 	if err != nil {
// 		return nil, fmt.Errorf("error retrieving predefined identifiers: %v", err)
// 	}

// 	return predefinedIdentifiers, nil
// }

func GetPredefinedIdentifiers(ctx context.Context, service *zscaler.Service, dictionaryName string) ([]string, int, error) {
	dictionary, err := GetByName(ctx, service, dictionaryName)
	if err != nil {
		return nil, 0, err
	}

	var predefinedIdentifiers []string
	endpoint := fmt.Sprintf("%s/%d/predefinedIdentifiers", dlpDictionariesEndpoint, dictionary.ID)
	err = service.Client.Read(ctx, endpoint, &predefinedIdentifiers)
	if err != nil {
		return nil, 0, err
	}

	return predefinedIdentifiers, dictionary.ID, nil
}

func Create(ctx context.Context, service *zscaler.Service, dlpDictionariesID *DlpDictionary) (*DlpDictionary, *http.Response, error) {
	resp, err := service.Client.Create(ctx, dlpDictionariesEndpoint, *dlpDictionariesID)
	if err != nil {
		return nil, nil, err
	}

	createdDlpDictionary, ok := resp.(*DlpDictionary)
	if !ok {
		return nil, nil, errors.New("object returned from api was not a dlp dictionary pointer")
	}

	service.Client.GetLogger().Printf("[DEBUG]returning new custom dlp dictionary that uses patterns and phrases from create: %d", createdDlpDictionary.ID)
	return createdDlpDictionary, nil, nil
}

func Update(ctx context.Context, service *zscaler.Service, dlpDictionariesID int, dlpDictionaries *DlpDictionary) (*DlpDictionary, *http.Response, error) {
	resp, err := service.Client.UpdateWithPut(ctx, fmt.Sprintf("%s/%d", dlpDictionariesEndpoint, dlpDictionariesID), *dlpDictionaries)
	if err != nil {
		return nil, nil, err
	}
	updatedDlpDictionary, _ := resp.(*DlpDictionary)

	service.Client.GetLogger().Printf("[DEBUG]returning updates custom dlp dictionary that uses patterns and phrases from ppdate: %d", updatedDlpDictionary.ID)
	return updatedDlpDictionary, nil, nil
}

func DeleteDlpDictionary(ctx context.Context, service *zscaler.Service, dlpDictionariesID int) (*http.Response, error) {
	err := service.Client.Delete(ctx, fmt.Sprintf("%s/%d", dlpDictionariesEndpoint, dlpDictionariesID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func GetAll(ctx context.Context, service *zscaler.Service) ([]DlpDictionary, error) {
	var dictionaries []DlpDictionary
	err := common.ReadAllPages(ctx, service.Client, dlpDictionariesEndpoint, &dictionaries)
	return dictionaries, err
}
