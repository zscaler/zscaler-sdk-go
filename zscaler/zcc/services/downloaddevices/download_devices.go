package downloaddevices

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
)

const (
	downloadDevicesEndpoint = "/zcc/papi/public/v1/downloadDevices"
)

type DownloadDevicesQueryParams struct {
	OSTypes           string `url:"osTypes,omitempty"`
	RegistrationTypes string `url:"registrationTypes,omitempty"`
}

func DownloadDevices(ctx context.Context, service *zscaler.Service, osTypes, registrationTypes string, writer io.Writer) error {
	queryParams := DownloadDevicesQueryParams{
		OSTypes:           osTypes,
		RegistrationTypes: registrationTypes,
	}

	resp, err := service.Client.NewZccRequestDo(ctx, "GET", downloadDevicesEndpoint, queryParams, nil, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download devices: %s", resp.Status)
	}

	_, err = io.Copy(writer, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write response to writer: %v", err)
	}

	return nil
}
