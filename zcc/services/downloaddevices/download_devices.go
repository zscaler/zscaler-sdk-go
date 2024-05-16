package downloaddevices

import (
	"fmt"
	"io"
	"net/http"
)

const (
	downloadDevicesEndpoint = "/public/v1/downloadDevices"
)

type DownloadDevicesQueryParams struct {
	OSTypes           string `url:"osTypes,omitempty"`
	RegistrationTypes string `url:"registrationTypes,omitempty"`
}

func (service *Service) DownloadDevices(osTypes, registrationTypes string, writer io.Writer) error {
	queryParams := DownloadDevicesQueryParams{
		OSTypes:           osTypes,
		RegistrationTypes: registrationTypes,
	}

	resp, err := service.Client.NewRequestDo("GET", downloadDevicesEndpoint, queryParams, nil, nil)
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
