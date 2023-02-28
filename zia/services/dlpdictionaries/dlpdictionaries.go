package dlpdictionaries

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/zia/services/common"
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
	IDMProfileMatchAccuracy []IDMProfileMatchAccuracy `json:"idmProfileMatchAccuracy"`
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

type ValidateDLPPattern struct {
	Status        string `json:"status,omitempty"`
	ErrPosition   int    `json:"errPosition,omitempty"`
	ErrMsg        string `json:"errMsg,omitempty"`
	ErrParameter  string `json:"errParameter,omitempty"`
	ErrSuggestion string `json:"errSuggestion,omitempty"`
	IDList        []int  `json:"idList,omitempty"`
}

func (service *Service) Get(dlpDictionariesID int) (*DlpDictionary, error) {
	var dlpDictionary DlpDictionary
	err := service.Client.Read(fmt.Sprintf("%s/%d", dlpDictionariesEndpoint, dlpDictionariesID), &dlpDictionary)
	if err != nil {
		return nil, err
	}

	service.Client.Logger.Printf("[DEBUG]Returning dictionary from Get: %d", dlpDictionary.ID)
	return &dlpDictionary, nil
}

func (service *Service) GetByName(dictionaryName string) (*DlpDictionary, error) {
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

func (service *Service) Create(dlpDictionariesID *DlpDictionary) (*DlpDictionary, *http.Response, error) {
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

func (service *Service) Update(dlpDictionariesID int, dlpDictionaries *DlpDictionary) (*DlpDictionary, *http.Response, error) {
	resp, err := service.Client.UpdateWithPut(fmt.Sprintf("%s/%d", dlpDictionariesEndpoint, dlpDictionariesID), *dlpDictionaries)
	if err != nil {
		return nil, nil, err
	}
	updatedDlpDictionary, _ := resp.(*DlpDictionary)

	service.Client.Logger.Printf("[DEBUG]returning updates custom dlp dictionary that uses patterns and phrases from ppdate: %d", updatedDlpDictionary.ID)
	return updatedDlpDictionary, nil, nil
}

func (service *Service) DeleteDlpDictionary(dlpDictionariesID int) (*http.Response, error) {
	err := service.Client.Delete(fmt.Sprintf("%s/%d", dlpDictionariesEndpoint, dlpDictionariesID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (service *Service) GetAll() ([]DlpDictionary, error) {
	var dictionaries []DlpDictionary
	err := common.ReadAllPages(service.Client, dlpDictionariesEndpoint, &dictionaries)
	return dictionaries, err
}

func (service *Service) ValidateDLPPattern(validatePattern *ValidateDLPPattern) (*ValidateDLPPattern, error) {
	resp, err := service.Client.Create(validateDLPPatternEndpoint, validatePattern)
	if err != nil {
		return nil, err
	}

	createdDLPPattern, ok := resp.(*ValidateDLPPattern)
	if !ok {
		return nil, errors.New("object returned from api was not dlp pattern pointer")
	}

	service.Client.Logger.Printf("[DEBUG]returning new dlp pattern from create: %d", createdDLPPattern)
	return createdDLPPattern, nil
}
