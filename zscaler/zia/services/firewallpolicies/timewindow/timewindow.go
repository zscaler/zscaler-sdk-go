package timewindow

import (
	"context"
	"fmt"
	"strings"

	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler"
	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler/zia/services/common"
)

const (
	timeWindowEndpoint = "/zia/api/v1/timeWindows"
)

type TimeWindow struct {
	ID        int      `json:"id"`
	Name      string   `json:"name,omitempty"`
	StartTime int32    `json:"startTime,omitempty"`
	EndTime   int32    `json:"endTime,omitempty"`
	DayOfWeek []string `json:"dayOfWeek,omitempty"`
}

func GetTimeWindowByName(ctx context.Context, service *zscaler.Service, timeWindowName string) (*TimeWindow, error) {
	var timeWindow []TimeWindow
	err := common.ReadAllPages(ctx, service.Client, timeWindowEndpoint, &timeWindow)
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

func GetAll(ctx context.Context, service *zscaler.Service) ([]TimeWindow, error) {
	var timeWindow []TimeWindow
	err := common.ReadAllPages(ctx, service.Client, timeWindowEndpoint, &timeWindow)
	return timeWindow, err
}
