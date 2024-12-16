package dlp_incident_receiver_servers

import (
	"context"
	"fmt"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/common"
)

const (
	dlpIncidentReceiverEndpoint = "/zia/api/v1/incidentReceiverServers"
)

type IncidentReceiverServers struct {
	// The unique identifier for the Incident Receiver.
	ID int `json:"id"`

	// The Incident Receiver server name.
	Name string `json:"name,omitempty"`

	// The Incident Receiver server URL.
	URL string `json:"url,omitempty"`

	// The status of the Incident Receiver.
	Status string `json:"status,omitempty"`

	// The Incident Receiver server flag.
	Flags int `json:"flags,omitempty"`
}

func Get(ctx context.Context, service *zscaler.Service, receiverID int) (*IncidentReceiverServers, error) {
	var incidentReceiver IncidentReceiverServers
	err := service.Client.Read(ctx, fmt.Sprintf("%s/%d", dlpIncidentReceiverEndpoint, receiverID), &incidentReceiver)
	if err != nil {
		return nil, err
	}

	service.Client.GetLogger().Printf("[DEBUG]Returning dlp incident receiver from Get: %d", incidentReceiver.ID)
	return &incidentReceiver, nil
}

func GetByName(ctx context.Context, service *zscaler.Service, receiverName string) (*IncidentReceiverServers, error) {
	var incidentReceiver []IncidentReceiverServers
	err := common.ReadAllPages(ctx, service.Client, dlpIncidentReceiverEndpoint, &incidentReceiver)
	if err != nil {
		return nil, err
	}
	for _, receiver := range incidentReceiver {
		if strings.EqualFold(receiver.Name, receiverName) {
			return &receiver, nil
		}
	}
	return nil, fmt.Errorf("no dlp incident receiver found with name: %s", receiverName)
}

func GetAll(ctx context.Context, service *zscaler.Service) ([]IncidentReceiverServers, error) {
	var incidentReceiver []IncidentReceiverServers
	err := common.ReadAllPages(ctx, service.Client, dlpIncidentReceiverEndpoint, &incidentReceiver)
	return incidentReceiver, err
}
