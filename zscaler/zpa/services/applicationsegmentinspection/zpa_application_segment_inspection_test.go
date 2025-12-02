package applicationsegmentinspection

import (
	"context"
	"testing"
	"time"

	"github.com/zscaler/zscaler-sdk-go/v3/tests"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/bacertificate"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/segmentgroup"
)

func TestAppSegmentInspectionInspection(t *testing.T) {
	name := tests.GetTestName("tests-inspapp")
	updateName := tests.GetTestName("tests-inspapp")
	segmentGroupName := tests.GetTestName("tests-inspapp")

	client, err := tests.NewVCRTestClient(t, "applicationsegmentinspection", "zpa")
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	defer client.Stop()
	service := client.Service

	// service, err := tests.NewZPAClient()
	// if err != nil {
	// 	t.Fatalf("Error creating client: %v", err)
	// }

	appGroup := segmentgroup.SegmentGroup{
		Name:        segmentGroupName,
		Description: segmentGroupName,
		Enabled:     true,
	}
	createdAppGroup, _, err := segmentgroup.Create(context.Background(), service, &appGroup)
	if err != nil {
		t.Errorf("Error creating application segment group: %v", err)
		return
	}
	defer func() {
		time.Sleep(time.Second * 2) // Sleep for 2 seconds before deletion
		_, _, getErr := segmentgroup.Get(context.Background(), service, createdAppGroup.ID)
		if getErr != nil {
			t.Logf("Resource might have already been deleted: %v", getErr)
		} else {
			_, err := segmentgroup.Delete(context.Background(), service, createdAppGroup.ID)
			if err != nil {
				t.Errorf("Error deleting application segment group: %v", err)
			}
		}
	}()

	certificateList, _, err := bacertificate.GetAll(context.Background(), service)
	if err != nil {
		t.Errorf("Error getting saml attributes: %v", err)
		return
	}
	if len(certificateList) == 0 {
		t.Error("Expected retrieved saml attributes to be non-empty, but got empty slice")
	}

	appSegment := AppSegmentInspection{
		Name:             name,
		Description:      name,
		Enabled:          true,
		SegmentGroupID:   createdAppGroup.ID,
		SegmentGroupName: createdAppGroup.Name,
		IsCnameEnabled:   true,
		BypassType:       "NEVER",
		ICMPAccessType:   "PING_TRACEROUTING",
		HealthReporting:  "ON_ACCESS",
		HealthCheckType:  "DEFAULT",
		TCPKeepAlive:     "1",
		DomainNames:      []string{"server1.bd-hashicorp.com"},
		TCPAppPortRange: []common.NetworkPorts{
			{
				From: "8443",
				To:   "8443",
			},
		},
		CommonAppsDto: CommonAppsDto{
			AppsConfig: []AppsConfig{
				{
					Name:                name,
					Description:         name,
					Enabled:             true,
					AppTypes:            []string{"INSPECT"},
					ApplicationPort:     "8443",
					ApplicationProtocol: "HTTPS",
					Domain:              "server1.bd-hashicorp.com",
					CertificateID:       certificateList[0].ID,
				},
			},
		},
	}

	// Test resource creation
	createdResource, _, err := Create(context.Background(), service, appSegment)
	// Check if the request was successful
	if err != nil {
		t.Errorf("Error making POST request: %v", err)
	}

	if createdResource.ID == "" {
		t.Error("Expected created resource ID to be non-empty, but got ''")
	}
	if createdResource.Name != name {
		t.Errorf("Expected created resource name '%s', but got '%s'", name, createdResource.Name)
	}

	// Test resource retrieval
	retrievedResource, _, err := Get(context.Background(), service, createdResource.ID)
	if err != nil {
		t.Errorf("Error retrieving resource: %v", err)
	}
	if retrievedResource.ID != createdResource.ID {
		t.Errorf("Expected retrieved resource ID '%s', but got '%s'", createdResource.ID, retrievedResource.ID)
	}
	if retrievedResource.Name != name {
		t.Errorf("Expected retrieved resource name '%s', but got '%s'", name, createdResource.Name)
	}
	retrievedResource.Name = updateName

	_, err = Update(context.Background(), service, createdResource.ID, retrievedResource)
	if err != nil {
		t.Errorf("Error updating resource: %v", err)
	}
	updatedResource, _, err := Get(context.Background(), service, createdResource.ID)
	if err != nil {
		t.Errorf("Error retrieving resource: %v", err)
	}
	if updatedResource.ID != createdResource.ID {
		t.Errorf("Expected retrieved updated resource ID '%s', but got '%s'", createdResource.ID, updatedResource.ID)
	}
	if updatedResource.Name != updateName {
		t.Errorf("Expected retrieved updated resource name '%s', but got '%s'", updateName, updatedResource.Name)
	}

	// Test resource retrieval by name
	retrievedResource, _, err = GetByName(context.Background(), service, updateName)
	if err != nil {
		t.Errorf("Error retrieving resource by name: %v", err)
	}
	if retrievedResource.ID != createdResource.ID {
		t.Errorf("Expected retrieved resource ID '%s', but got '%s'", createdResource.ID, retrievedResource.ID)
	}
	if retrievedResource.Name != updateName {
		t.Errorf("Expected retrieved resource name '%s', but got '%s'", updateName, createdResource.Name)
	}
	// Test resource removal
	_, err = Delete(context.Background(), service, createdResource.ID)
	if err != nil {
		t.Errorf("Error deleting resource: %v", err)
		return
	}

	// Test resource retrieval after deletion
	_, _, err = Get(context.Background(), service, createdResource.ID)
	if err == nil {
		t.Errorf("Expected error retrieving deleted resource, but got nil")
	}
}

func TestRetrieveNonExistentResource(t *testing.T) {
	client, err := tests.NewVCRTestClient(t, "applicationsegmentinspection", "zpa")
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	defer client.Stop()
	service := client.Service

	// service, err := tests.NewZPAClient()
	// if err != nil {
	// 	t.Fatalf("Error creating client: %v", err)
	// }

	_, _, err = Get(context.Background(), service, "non_existent_id")
	if err == nil {
		t.Error("Expected error retrieving non-existent resource, but got nil")
	}
}

func TestDeleteNonExistentResource(t *testing.T) {
	client, err := tests.NewVCRTestClient(t, "applicationsegmentinspection", "zpa")
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	defer client.Stop()
	service := client.Service

	// service, err := tests.NewZPAClient()
	// if err != nil {
	// 	t.Fatalf("Error creating client: %v", err)
	// }

	_, err = Delete(context.Background(), service, "non_existent_id")
	if err == nil {
		t.Error("Expected error deleting non-existent resource, but got nil")
	}
}

func TestUpdateNonExistentResource(t *testing.T) {
	client, err := tests.NewVCRTestClient(t, "applicationsegmentinspection", "zpa")
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	defer client.Stop()
	service := client.Service

	// service, err := tests.NewZPAClient()
	// if err != nil {
	// 	t.Fatalf("Error creating client: %v", err)
	// }

	_, err = Update(context.Background(), service, "non_existent_id", &AppSegmentInspection{})
	if err == nil {
		t.Error("Expected error updating non-existent resource, but got nil")
	}
}

func TestGetByNameNonExistentResource(t *testing.T) {
	client, err := tests.NewVCRTestClient(t, "applicationsegmentinspection", "zpa")
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	defer client.Stop()
	service := client.Service

	// service, err := tests.NewZPAClient()
	// if err != nil {
	// 	t.Fatalf("Error creating client: %v", err)
	// }

	_, _, err = GetByName(context.Background(), service, "non_existent_name")
	if err == nil {
		t.Error("Expected error retrieving resource by non-existent name, but got nil")
	}
}
