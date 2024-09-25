package dlp_engines

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
	dlpEnginesEndpoint    = "/zia/api/v1/dlpEngines"
	dlpEngineLiteEndpoint = "/zia/api/v1/dlpEngines/lite"
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

func Get(ctx context.Context, service *zscaler.Service, engineID int) (*DLPEngines, error) {
	var dlpEngines DLPEngines
	err := service.Client.Read(ctx, fmt.Sprintf("%s/%d", dlpEnginesEndpoint, engineID), &dlpEngines)
	if err != nil {
		return nil, err
	}

	service.Client.Logger.Printf("[DEBUG]Returning dlp engine from Get: %d", dlpEngines.ID)
	return &dlpEngines, nil
}

func GetByName(ctx context.Context, service *zscaler.Service, engineName string) (*DLPEngines, error) {
	dlpEngines, err := GetAll(ctx, service)
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

func Create(ctx context.Context, service *zscaler.Service, engineID *DLPEngines) (*DLPEngines, *http.Response, error) {
	resp, err := service.Client.Create(ctx, dlpEnginesEndpoint, *engineID)
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

func Update(ctx context.Context, service *zscaler.Service, engineID int, engines *DLPEngines) (*DLPEngines, *http.Response, error) {
	resp, err := service.Client.UpdateWithPut(ctx, fmt.Sprintf("%s/%d", dlpEnginesEndpoint, engineID), *engines)
	if err != nil {
		return nil, nil, err
	}
	updatedDlpEngine, _ := resp.(*DLPEngines)

	service.Client.Logger.Printf("[DEBUG]returning updates dlp engine from update: %d", updatedDlpEngine.ID)
	return updatedDlpEngine, nil, nil
}

func Delete(ctx context.Context, service *zscaler.Service, engineID int) (*http.Response, error) {
	err := service.Client.Delete(ctx, fmt.Sprintf("%s/%d", dlpEnginesEndpoint, engineID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func GetAll(ctx context.Context, service *zscaler.Service) ([]DLPEngines, error) {
	var dlpEngines []DLPEngines
	err := common.ReadAllPages(ctx, service.Client, dlpEnginesEndpoint, &dlpEngines)
	return dlpEngines, err
}

// Functions to for DLP Engine Lite query
func GetEngineLiteID(ctx context.Context, service *zscaler.Service, engineID int) (*DLPEngines, error) {
	dlpEngines, err := GetAllEngineLite(ctx, service)
	if err != nil {
		return nil, err
	}
	for _, engine := range dlpEngines {
		if engine.ID == engineID {
			return &engine, nil
		}
	}
	return nil, fmt.Errorf("no dlp engine found with ID: %d", engineID)
}

func GetByPredefinedEngine(ctx context.Context, service *zscaler.Service, engineName string) (*DLPEngines, error) {
	dlpEngines, err := GetAllEngineLite(ctx, service)
	if err != nil {
		return nil, err
	}
	for _, engine := range dlpEngines {
		if strings.EqualFold(engine.PredefinedEngineName, engineName) {
			return &engine, nil
		}
	}
	return nil, fmt.Errorf("no predefined dlp engine found with name: %s", engineName)
}

func GetAllEngineLite(ctx context.Context, service *zscaler.Service) ([]DLPEngines, error) {
	var engines []DLPEngines
	err := common.ReadAllPages(ctx, service.Client, dlpEngineLiteEndpoint, &engines)
	return engines, err
}
