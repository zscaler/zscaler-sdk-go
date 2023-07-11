package integration

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/zscaler/zscaler-sdk-go/tests"
	"github.com/zscaler/zscaler-sdk-go/zia/services/trafficforwarding/staticips"
	"github.com/zscaler/zscaler-sdk-go/zia/services/trafficforwarding/vpncredentials"
)

func TestTrafficForwardingVPNCreds(t *testing.T) {
	ipAddress, _ := acctest.RandIpAddress("1.1.1.1/24")
	comment := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	updateComment := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}
	// static ip for vpn credentials testing
	staticipsService := staticips.New(client)
	// Test resource creation
	staticIP, _, err := staticipsService.Create(&staticips.StaticIP{
		IpAddress: ipAddress,
		Comment:   comment,
	})
	if err != nil {
		t.Errorf("creating static ip failed: %v", err)
		return
	}
	defer func() {
		_, err := staticipsService.Delete(staticIP.ID)
		if err != nil {
			t.Errorf("deleting static ip failed: %v", err)
		}
	}()

	service := vpncredentials.New(client)

	cred := vpncredentials.VPNCredentials{
		Type:         "IP",
		IPAddress:    ipAddress,
		Comments:     comment,
		PreSharedKey: "newPassword123!",
	}

	// Test resource creation
	createdResource, _, err := service.Create(&cred)

	// Check if the request was successful
	if err != nil {
		t.Errorf("Error making POST request: %v", err)
	}

	if createdResource.ID == 0 {
		t.Error("Expected created resource ID to be non-empty, but got ''")
	}
	if createdResource.Comments != comment {
		t.Errorf("Expected created resource comment '%s', but got '%s'", comment, createdResource.Comments)
	}
	// Test resource retrieval
	retrievedResource, err := service.Get(createdResource.ID)
	if err != nil {
		t.Errorf("Error retrieving resource: %v", err)
	}
	if retrievedResource.ID != createdResource.ID {
		t.Errorf("Expected retrieved resource ID '%d', but got '%d'", createdResource.ID, retrievedResource.ID)
	}
	if retrievedResource.Comments != comment {
		t.Errorf("Expected retrieved resource comment '%s', but got '%s'", comment, createdResource.Comments)
	}
	// Test resource update
	retrievedResource.Comments = updateComment
	_, _, err = service.Update(createdResource.ID, retrievedResource)
	if err != nil {
		t.Errorf("Error updating resource: %v", err)
	}
	updatedResource, err := service.Get(createdResource.ID)
	if err != nil {
		t.Errorf("Error retrieving resource: %v", err)
	}
	if updatedResource.ID != createdResource.ID {
		t.Errorf("Expected retrieved updated resource ID '%d', but got '%d'", createdResource.ID, updatedResource.ID)
	}
	if updatedResource.Comments != updateComment {
		t.Errorf("Expected retrieved updated resource comment '%s', but got '%s'", updateComment, updatedResource.Comments)
	}

	// Test resource retrieval by name
	retrievedResource, err = service.GetVPNByType("IP")
	if err != nil {
		t.Errorf("Error retrieving resource by name: %v", err)
	}
	if retrievedResource.ID != createdResource.ID {
		t.Errorf("Expected retrieved resource ID '%d', but got '%d'", createdResource.ID, retrievedResource.ID)
	}
	if retrievedResource.Comments != updateComment {
		t.Errorf("Expected retrieved resource comment '%s', but got '%s'", updateComment, retrievedResource.Comments)
	}
	// Test resources retrieval
	resources, err := service.GetAll()
	if err != nil {
		t.Errorf("Error retrieving resources: %v", err)
	}
	if len(resources) == 0 {
		t.Error("Expected retrieved resources to be non-empty, but got empty slice")
	}
	// check if the created resource is in the list
	found := false
	for _, resource := range resources {
		if resource.ID == createdResource.ID {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected retrieved resources to contain created resource '%d', but it didn't", createdResource.ID)
	}
	// Test resource removal
	err = service.Delete(createdResource.ID)
	if err != nil {
		t.Errorf("Error deleting resource: %v", err)
		return
	}

	// Test resource retrieval after deletion
	_, err = service.Get(createdResource.ID)
	if err == nil {
		t.Errorf("Expected error retrieving deleted resource, but got nil")
	}

}
