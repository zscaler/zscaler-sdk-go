package timewindow

import (
	"fmt"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v1/zia/services/common"
)

const (
	timeWindowEndpoint = "/timeWindows"
)

type TimeWindow struct {
	ID        int      `json:"id"`
	Name      string   `json:"name,omitempty"`
	StartTime int32    `json:"startTime,omitempty"`
	EndTime   int32    `json:"endTime,omitempty"`
	DayOfWeek []string `json:"dayOfWeek,omitempty"`
}

func (service *Service) GetTimeWindow(timeWindowID int) (*TimeWindow, error) {
	var timeWindow TimeWindow
	err := service.Client.Read(fmt.Sprintf("%s/%d", timeWindowEndpoint, timeWindowID), &timeWindow)
	if err != nil {
		return nil, err
	}

	service.Client.Logger.Printf("[DEBUG]Returning time window from Get: %d", timeWindow.ID)
	return &timeWindow, nil
}

func (service *Service) GetTimeWindowByName(timeWindowName string) (*TimeWindow, error) {
	var timeWindow []TimeWindow
	err := common.ReadAllPages(service.Client, timeWindowEndpoint, &timeWindow)
	if err != nil {
		return nil, err
	}
	for _, timeWindow := range timeWindow {
		if strings.EqualFold(timeWindow.Name, timeWindowName) {
			return &timeWindow, nil
		}
	}
	return nil, fmt.Errorf("no time window found with name: %s", timeWindowName)
}

func (service *Service) GetAll() ([]TimeWindow, error) {
	var timeWindow []TimeWindow
	err := common.ReadAllPages(service.Client, timeWindowEndpoint, &timeWindow)
	return timeWindow, err
}
