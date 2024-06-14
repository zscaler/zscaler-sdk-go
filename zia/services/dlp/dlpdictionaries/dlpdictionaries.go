package dlpdictionaries

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v2/zia/services"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/common"
)

const (
	dlpDictionariesEndpoint    = "/dlpDictionaries"
	validateDLPPatternEndpoint = "/dlpDictionaries/validateDlpPattern"
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
	Custom bool `json:"custom"`

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

	// This value is set to true if proximity length and high confidence phrases are enabled for the DLP dictionary.
	ProximityLengthEnabled bool `json:"proximityLengthEnabled,omitempty"`
}

type Phrases struct {
	// The action applied to a DLP dictionary using phrases
	Action string `json:"action,omitempty"`

	// DLP dictionary phrase
	Phrase string `json:"phrase,omitempty"`
}

type Patterns struct {
	// The action applied to a DLP dictionary using patterns
	Action string `json:"action,omitempty"`

	// DLP dictionary pattern
	Pattern string `json:"pattern,omitempty"`
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

func Get(service *services.Service, dlpDictionariesID int) (*DlpDictionary, error) {
	var dlpDictionary DlpDictionary
	err := service.Client.Read(fmt.Sprintf("%s/%d", dlpDictionariesEndpoint, dlpDictionariesID), &dlpDictionary)
	if err != nil {
		return nil, err
	}

	service.Client.Logger.Printf("[DEBUG]Returning dictionary from Get: %d", dlpDictionary.ID)
	return &dlpDictionary, nil
}

func GetByName(service *services.Service, dictionaryName string) (*DlpDictionary, error) {
	var dictionaries []DlpDictionary
	err := common.ReadAllPages(service.Client, dlpDictionariesEndpoint, &dictionaries)
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

func Create(service *services.Service, dlpDictionariesID *DlpDictionary) (*DlpDictionary, *http.Response, error) {
	resp, err := service.Client.Create(dlpDictionariesEndpoint, *dlpDictionariesID)
	if err != nil {
		return nil, nil, err
	}

	createdDlpDictionary, ok := resp.(*DlpDictionary)
	if !ok {
		return nil, nil, errors.New("object returned from api was not a dlp dictionary pointer")
	}

	service.Client.Logger.Printf("[DEBUG]returning new custom dlp dictionary that uses patterns and phrases from create: %d", createdDlpDictionary.ID)
	return createdDlpDictionary, nil, nil
}

func Update(service *services.Service, dlpDictionariesID int, dlpDictionaries *DlpDictionary) (*DlpDictionary, *http.Response, error) {
	resp, err := service.Client.UpdateWithPut(fmt.Sprintf("%s/%d", dlpDictionariesEndpoint, dlpDictionariesID), *dlpDictionaries)
	if err != nil {
		return nil, nil, err
	}
	updatedDlpDictionary, _ := resp.(*DlpDictionary)

	service.Client.Logger.Printf("[DEBUG]returning updates custom dlp dictionary that uses patterns and phrases from ppdate: %d", updatedDlpDictionary.ID)
	return updatedDlpDictionary, nil, nil
}

func DeleteDlpDictionary(service *services.Service, dlpDictionariesID int) (*http.Response, error) {
	err := service.Client.Delete(fmt.Sprintf("%s/%d", dlpDictionariesEndpoint, dlpDictionariesID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func GetAll(service *services.Service) ([]DlpDictionary, error) {
	var dictionaries []DlpDictionary
	err := common.ReadAllPages(service.Client, dlpDictionariesEndpoint, &dictionaries)
	return dictionaries, err
}
