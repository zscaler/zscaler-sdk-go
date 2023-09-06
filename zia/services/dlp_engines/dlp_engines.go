package dlp_engines

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/common"
)

const (
	dlpEnginesEndpoint = "/dlpEngines"
)

type DLPEngines struct {
	// The unique identifier for the DLP engine.
	ID int `json:"id"`

	// The DLP engine name as configured by the admin. This attribute is required in POST and PUT requests for custom DLP engines.
	Name string `json:"name,omitempty"`

	// The DLP engine's description.
	Description string `json:"description,omitempty"`

	// The name of the predefined DLP engine.
	PredefinedEngineName string `json:"predefinedEngineName,omitempty"`

	// The boolean logical operator in which various DLP dictionaries are combined within a DLP engine's expression.
	EngineExpression string `json:"engineExpression,omitempty"`

	// Indicates whether this is a custom DLP engine. If this value is set to true, the engine is custom.
	CustomDlpEngine bool `json:"customDlpEngine,omitempty"`
}

func (service *Service) Get(engineID int) (*DLPEngines, error) {
	var dlpEngines DLPEngines
	err := service.Client.Read(fmt.Sprintf("%s/%d", dlpEnginesEndpoint, engineID), &dlpEngines)
	if err != nil {
		return nil, err
	}

	service.Client.Logger.Printf("[DEBUG]Returning dlp engine from Get: %d", dlpEngines.ID)
	return &dlpEngines, nil
}

func (service *Service) GetByName(engineName string) (*DLPEngines, error) {
	dlpEngines, err := service.GetAll()
	if err != nil {
		return nil, err
	}
	for _, engine := range dlpEngines {
		if strings.EqualFold(engine.Name, engineName) {
			return &engine, nil
		}
	}
	return nil, fmt.Errorf("no dlp engine found with name: %s", engineName)
}

func (service *Service) Create(engineID *DLPEngines) (*DLPEngines, *http.Response, error) {
	resp, err := service.Client.Create(dlpEnginesEndpoint, *engineID)
	if err != nil {
		return nil, nil, err
	}

	createdDlpEngine, ok := resp.(*DLPEngines)
	if !ok {
		return nil, nil, errors.New("object returned from api was not a dlp engine pointer")
	}

	service.Client.Logger.Printf("[DEBUG]returning new dlp engine from create: %d", createdDlpEngine.ID)
	return createdDlpEngine, nil, nil
}

func (service *Service) Update(engineID int, engines *DLPEngines) (*DLPEngines, *http.Response, error) {
	resp, err := service.Client.UpdateWithPut(fmt.Sprintf("%s/%d", dlpEnginesEndpoint, engineID), *engines)
	if err != nil {
		return nil, nil, err
	}
	updatedDlpEngine, _ := resp.(*DLPEngines)

	service.Client.Logger.Printf("[DEBUG]returning updates dlp engine from update: %d", updatedDlpEngine.ID)
	return updatedDlpEngine, nil, nil
}

func (service *Service) Delete(engineID int) (*http.Response, error) {
	err := service.Client.Delete(fmt.Sprintf("%s/%d", dlpEnginesEndpoint, engineID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (service *Service) GetAll() ([]DLPEngines, error) {
	var dlpEngines []DLPEngines
	err := common.ReadAllPages(service.Client, dlpEnginesEndpoint, &dlpEngines)
	for i := range dlpEngines {
		if dlpEngines[i].Name == "" && dlpEngines[i].PredefinedEngineName != "" {
			dlpEngines[i].Name = dlpEngines[i].PredefinedEngineName
		}
	}
	return dlpEngines, err
}
