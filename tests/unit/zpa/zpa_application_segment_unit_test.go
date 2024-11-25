package unit

/*
import (
	"context"
	"net/http"
	"testing"

	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/tests"
	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler/zpa/services"
	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler/zpa/services/applicationsegment"
)

func TestApplicationSegment_GetByName(t *testing.T) {
	client, mux, server := tests.NewOneAPIClientMock()
	defer server.Close()

	service := services.New(client)
	mux.HandleFunc("/mgmtconfig/v1/admin/customers/customerid/application", func(w http.ResponseWriter, r *http.Request) {
		// Write a JSON response with an array of application segments
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"list":[
			{
				"id": "123",
				"domainNames": ["example.com"],
				"name": "App1",
				"description": "Application 1",
				"enabled": true,
				"passiveHealthEnabled": true,
				"doubleEncrypt": false,
				"configSpace": "",
				"applications": "",
				"bypassType": "",
				"healthCheckType": "",
				"isCnameEnabled": false,
				"ipAnchored": false,
				"healthReporting": "",
				"selectConnectorCloseToApp": true,
				"icmpAccessType": "",
				"appRecommendationId": "",
				"segmentGroupId": "",
				"segmentGroupName": "",
				"creationTime": "",
				"modifiedBy": "",
				"modifiedTime": "",
				"tcpKeepAlive": "",
				"isIncompleteDRConfig": false,
				"useInDrMode": false,
				"tcpPortRanges": [],
				"udpPortRanges": [],
				"tcpPortRange": [],
				"udpPortRange": [],
				"serverGroups": [],
				"defaultIdleTimeout": "",
				"defaultMaxAge": "",
				"commonAppsDto": {},
				"clientlessApps": []
			},
			{
				"id": "456",
				"domainNames": ["example.org"],
				"name": "App2",
				"description": "Application 2",
				"enabled": true,
				"passiveHealthEnabled": false,
				"doubleEncrypt": true,
				"configSpace": "",
				"applications": "",
				"bypassType": "",
				"healthCheckType": "",
				"isCnameEnabled": true,
				"ipAnchored": true,
				"healthReporting": "",
				"selectConnectorCloseToApp": false,
				"icmpAccessType": "",
				"appRecommendationId": "",
				"segmentGroupId": "",
				"segmentGroupName": "",
				"creationTime": "",
				"modifiedBy": "",
				"modifiedTime": "",
				"tcpKeepAlive": "",
				"isIncompleteDRConfig": false,
				"useInDrMode": false,
				"tcpPortRanges": [],
				"udpPortRanges": [],
				"tcpPortRange": [],
				"udpPortRange": [],
				"serverGroups": [],
				"defaultIdleTimeout": "",
				"defaultMaxAge": "",
				"commonAppsDto": {},
				"clientlessApps": []
			}
		],
		"totalPages": 1
		}`))
	})

	// Make the GetByName request
	appSegment, _, err := applicationsegment.GetByName(context.Background(), service, "App1")
	if err != nil {
		t.Errorf("GetByName returned an error: %v", err)
	}

	// Check the ID and name of the application segment
	if appSegment.ID != "123" || appSegment.Name != "App1" {
		t.Errorf("Expected application segment ID '123' and name 'App1', but got ID '%s' and name '%s'", appSegment.ID, appSegment.Name)
	}
}

func TestApplicationSegment_Create(t *testing.T) {
	client, mux, server := tests.NewOneAPIClientMock()
	defer server.Close()

	service := services.New(client)

	mux.HandleFunc("/mgmtconfig/v1/admin/customers/customerid/application", func(w http.ResponseWriter, r *http.Request) {
		// Check the HTTP method
		if r.Method != http.MethodPost {
			t.Errorf("Expected HTTP method 'POST', but got '%s'", r.Method)
		}

		// Check the request body
		reqBody := &applicationsegment.ApplicationSegmentResource{}
		tests.ParseJSONRequest(t, r, reqBody)

		// Create a new application segment with the provided data
		respBody := &applicationsegment.ApplicationSegmentResource{
			ID:          "789",
			Name:        reqBody.Name,
			Description: reqBody.Description,
			// Include other fields as needed
		}

		// Write a JSON response with the created application segment
		tests.WriteJSONResponse(t, w, http.StatusOK, respBody)
	})

	// Create a new application segment
	appSegment := applicationsegment.ApplicationSegmentResource{
		Name:        "NewApp",
		Description: "New application segment",
		// Include other fields as needed
	}

	// Make the Create request
	createdAppSegment, _, err := applicationsegment.Create(context.Background(), service, appSegment)
	if err != nil {
		t.Errorf("Create returned an error: %v", err)
	}

	// Check the ID and name of the created application segment
	if createdAppSegment.ID != "789" || createdAppSegment.Name != "NewApp" {
		t.Errorf("Expected created application segment ID '789' and name 'NewApp', but got ID '%s' and name '%s'", createdAppSegment.ID, createdAppSegment.Name)
	}
}

func TestApplicationSegment_Update(t *testing.T) {
	client, mux, server := tests.NewOneAPIClientMock()
	defer server.Close()

	service := services.New(client)

	mux.HandleFunc("/mgmtconfig/v1/admin/customers/customerid/application/", func(w http.ResponseWriter, r *http.Request) {
		// Check the HTTP method
		if r.Method != http.MethodPut {
			t.Errorf("Expected HTTP method 'PUT', but got '%s'", r.Method)
		}

		// Check the request body
		reqBody := &applicationsegment.ApplicationSegmentResource{}
		tests.ParseJSONRequest(t, r, reqBody)

		// Update the application segment with the provided data
		respBody := &applicationsegment.ApplicationSegmentResource{
			ID:          reqBody.ID,
			Name:        reqBody.Name,
			Description: reqBody.Description,
			// Include other fields as needed
		}

		// Write a JSON response with the updated application segment

		tests.WriteJSONResponse(t, w, http.StatusOK, respBody)
	})

	// Update an existing application segment
	appSegment := applicationsegment.ApplicationSegmentResource{
		ID:          "123",
		Name:        "UpdatedApp",
		Description: "Updated application segment",
		// Include other fields as needed
	}

	// Make the Update request
	_, err := applicationsegment.Update(context.Background(), service, "123", appSegment)
	if err != nil {
		t.Errorf("Update returned an error: %v", err)
	}
}

func TestApplicationSegment_Delete(t *testing.T) {
	client, mux, server := tests.NewOneAPIClientMock()
	defer server.Close()

	service := services.New(client)

	mux.HandleFunc("/mgmtconfig/v1/admin/customers/customerid/application/", func(w http.ResponseWriter, r *http.Request) {
		// Check the HTTP method
		if r.Method != http.MethodDelete {
			t.Errorf("Expected HTTP method 'DELETE', but got '%s'", r.Method)
		}

		// Write a JSON response with a success message
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "Application segment deleted successfully"}`))
	})

	// Delete an existing application segment
	_, err := applicationsegment.Delete(context.Background(), service, "123")
	if err != nil {
		t.Errorf("Delete returned an error: %v", err)
	}
}

func TestApplicationSegment_GetAll(t *testing.T) {
	client, mux, server := tests.NewOneAPIClientMock()
	defer server.Close()

	service := services.New(client)

	mux.HandleFunc("/mgmtconfig/v1/admin/customers/customerid/application", func(w http.ResponseWriter, r *http.Request) {
		// Write a JSON response with an array of application segments
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"list": [
			{
				"id": "123",
				"domainNames": ["example.com"],
				"name": "App1",
				"description": "Application 1",
				"enabled": true,
				"passiveHealthEnabled": true,
				"doubleEncrypt": false,
				"configSpace": "",
				"applications": "",
				"bypassType": "",
				"healthCheckType": "",
				"isCnameEnabled": false,
				"ipAnchored": false,
				"healthReporting": "",
				"selectConnectorCloseToApp": true,
				"icmpAccessType": "",
				"appRecommendationId": "",
				"segmentGroupId": "",
				"segmentGroupName": "",
				"creationTime": "",
				"modifiedBy": "",
				"modifiedTime": "",
				"tcpKeepAlive": "",
				"isIncompleteDRConfig": false,
				"useInDrMode": false,
				"tcpPortRanges": [],
				"udpPortRanges": [],
				"tcpPortRange": [],
				"udpPortRange": [],
				"serverGroups": [],
				"defaultIdleTimeout": "",
				"defaultMaxAge": "",
				"commonAppsDto": {},
				"clientlessApps": []
			},
			{
				"id": "456",
				"domainNames": ["example.org"],
				"name": "App2",
				"description": "Application 2",
				"enabled": true,
				"passiveHealthEnabled": false,
				"doubleEncrypt": true,
				"configSpace": "",
				"applications": "",
				"bypassType": "",
				"healthCheckType": "",
				"isCnameEnabled": true,
				"ipAnchored": true,
				"healthReporting": "",
				"selectConnectorCloseToApp": false,
				"icmpAccessType": "",
				"appRecommendationId": "",
				"segmentGroupId": "",
				"segmentGroupName": "",
				"creationTime": "",
				"modifiedBy": "",
				"modifiedTime": "",
				"tcpKeepAlive": "",
				"isIncompleteDRConfig": false,
				"useInDrMode": false,
				"tcpPortRanges": [],
				"udpPortRanges": [],
				"tcpPortRange": [],
				"udpPortRange": [],
				"serverGroups": [],
				"defaultIdleTimeout": "",
				"defaultMaxAge": "",
				"commonAppsDto": {},
				"clientlessApps": []
			}
		],
		"totalPages": 1
		}`))
	})

	// Make the GetAll request
	appSegments, _, err := applicationsegment.GetAll(context.Background(), service)
	if err != nil {
		t.Errorf("GetAll returned an error: %v", err)
	}

	// Check the number of application segments returned
	expectedCount := 2
	if len(appSegments) != expectedCount {
		t.Errorf("Expected %d application segments, but got %d", expectedCount, len(appSegments))
	}
}
*/
