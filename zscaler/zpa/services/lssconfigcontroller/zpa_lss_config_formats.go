package lssconfigcontroller

import (
	"context"
	"fmt"
	"net/http"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
)

type LSSFormats struct {
	Csv  string `json:"csv"`
	Tsv  string `json:"tsv"`
	Json string `json:"json"`
}

func GetFormats(ctx context.Context, service *zscaler.Service, logType string) (*LSSFormats, *http.Response, error) {
	v := new(LSSFormats)
	relativeURL := fmt.Sprintf("%slssConfig/logType/formats", mgmtConfigTypesAndFormats)
	resp, err := service.Client.NewRequestDo(ctx, "GET", relativeURL, struct {
		LogType string `url:"logType"`
	}{
		LogType: logType,
	}, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}
