package analysis

import (
	"context"
	"fmt"
	"net/http"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
)

const (
	analysisEndpoint = "/zdx/v1/analysis"
)

type AnalysisRequest struct {
	DeviceID int `json:"device_id"`
	AppID    int `json:"app_id"`
	T0       int `json:"t0"`
	T1       int `json:"t1"`
}

type AnalysisResult struct {
	ErrMsg string `json:"err_msg"`
	Result Result `json:"result"`
}

type Result struct {
	Issue      string `json:"issue"`
	Confidence int    `json:"confidence"`
	Message    string `json:"message"`
	Times      []int  `json:"times"`
}

func GetAnalysis(ctx context.Context, service *zscaler.Service, analysisID string) (*AnalysisResult, *http.Response, error) {
	var response AnalysisResult
	path := fmt.Sprintf("%s/%s", analysisEndpoint, analysisID)
	resp, err := service.Client.NewRequestDo(ctx, "GET", path, nil, nil, &response)
	if err != nil {
		return nil, nil, err
	}
	return &response, resp, nil
}

func CreateAnalysis(ctx context.Context, service *zscaler.Service, request AnalysisRequest) (*http.Response, error) {
	path := analysisEndpoint
	resp, err := service.Client.NewRequestDo(ctx, "POST", path, nil, request, nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func DeleteAnalysis(ctx context.Context, service *zscaler.Service, analysisID string) (*http.Response, error) {
	path := fmt.Sprintf("%s/%s", analysisEndpoint, analysisID)
	resp, err := service.Client.NewRequestDo(ctx, "DELETE", path, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
