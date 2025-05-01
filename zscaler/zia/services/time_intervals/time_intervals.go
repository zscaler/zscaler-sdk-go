package time_intervals

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
	timeIntervalEndpoint = "/zia/api/v1/timeIntervals"
)

type TimeInterval struct {
	ID int `json:"id,omitempty"`

	// Name to identify the time interval
	Name string `json:"name,omitempty"`

	// The time interval start time.
	StartTime int `json:"startTime"`

	// The time interval end time.
	EndTime int `json:"endTime,omitempty"`

	// Values supported: `EVERYDAY`, `SUN`, `MON`, `TUE`, `WED`, `THU`, `FRI`, `SAT`
	DaysOfWeek []string `json:"daysOfWeek,omitempty"`
}

func Get(ctx context.Context, service *zscaler.Service, intervalID int) (*TimeInterval, error) {
	var interval TimeInterval
	err := service.Client.Read(ctx, fmt.Sprintf("%s/%d", timeIntervalEndpoint, intervalID), &interval)
	if err != nil {
		return nil, err
	}

	service.Client.GetLogger().Printf("[DEBUG] Returning time interval from Get: %d", interval.ID)
	return &interval, nil
}

func GetTimeIntervalByName(ctx context.Context, service *zscaler.Service, timeIntervalName string) (*TimeInterval, error) {
	var timeInterval []TimeInterval
	err := common.ReadAllPages(ctx, service.Client, timeIntervalEndpoint, &timeInterval)
	if err != nil {
		return nil, err
	}
	for _, timeInterval := range timeInterval {
		if strings.EqualFold(timeInterval.Name, timeIntervalName) {
			return &timeInterval, nil
		}
	}
	return nil, fmt.Errorf("no time interval found with name: %s", timeIntervalName)
}

func Create(ctx context.Context, service *zscaler.Service, intervalID *TimeInterval) (*TimeInterval, *http.Response, error) {
	resp, err := service.Client.Create(ctx, timeIntervalEndpoint, *intervalID)
	if err != nil {
		return nil, nil, err
	}

	createdInterval, ok := resp.(*TimeInterval)
	if !ok {
		return nil, nil, errors.New("object returned from api was not a interval pointer")
	}

	service.Client.GetLogger().Printf("[DEBUG]returning new time interval from create: %d", createdInterval.ID)
	return createdInterval, nil, nil
}

func Update(ctx context.Context, service *zscaler.Service, intervalID int, interval *TimeInterval) (*TimeInterval, *http.Response, error) {
	resp, err := service.Client.UpdateWithPut(ctx, fmt.Sprintf("%s/%d", timeIntervalEndpoint, intervalID), *interval)
	if err != nil {
		return nil, nil, err
	}
	updatedInterval, _ := resp.(*TimeInterval)

	service.Client.GetLogger().Printf("[DEBUG]returning updates time interval  from update: %d", updatedInterval.ID)
	return updatedInterval, nil, nil
}

func Delete(ctx context.Context, service *zscaler.Service, intervalID int) (*http.Response, error) {
	err := service.Client.Delete(ctx, fmt.Sprintf("%s/%d", timeIntervalEndpoint, intervalID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func GetAll(ctx context.Context, service *zscaler.Service) ([]TimeInterval, error) {
	var timeInterval []TimeInterval
	err := common.ReadAllPages(ctx, service.Client, timeIntervalEndpoint, &timeInterval)
	return timeInterval, err
}
