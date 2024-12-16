package proxy_gateways

import (
	"context"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/tests"
)

func TestProxyGateway(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	// Retrieve all proxy gateways
	proxyGateways, err := GetAll(context.Background(), service)
	if err != nil {
		t.Errorf("Error getting proxy gateways: %v", err)
		return
	}
	if len(proxyGateways) == 0 {
		t.Log("No proxy gateways found")
		return
	}

	// Use the first proxy gateway for testing
	firstGW := proxyGateways[0]

	// Retrieve lite proxy gateways
	t.Run("GetLite", func(t *testing.T) {
		proxyGatewaysLite, err := GetLite(context.Background(), service)
		if err != nil {
			t.Errorf("Error getting lite proxy gateways: %v", err)
			return
		}
		if len(proxyGatewaysLite) == 0 {
			t.Log("No lite proxy gateways found")
			return
		}
	})

	// Test GetProxyGWByName
	t.Run("GetProxyGWByName", func(t *testing.T) {
		gwByName, err := GetByName(context.Background(), service, firstGW.Name)
		if err != nil {
			t.Errorf("Error getting Proxy gateway by name: %v", err)
			return
		}
		if gwByName == nil || gwByName.Name != firstGW.Name {
			t.Errorf("Proxy gateway name does not match: expected %s, got %s", firstGW.Name, gwByName.Name)
		}
	})
}
