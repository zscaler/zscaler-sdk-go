package unit

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v2/tests"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/servergroup"
)

func TestService_Get(t *testing.T) {
	client, mux, server := tests.NewZpaClientMock()
	defer server.Close()
	mux.HandleFunc("/mgmtconfig/v1/admin/customers/customerid/serverGroup/groupID", func(w http.ResponseWriter, r *http.Request) {
		// Write a mock response
		response := `{
			"id": "groupID",
			"enabled": true,
			"name": "Test Group",
			"description": "Test description",
			"ipAnchored": true,
			"configSpace": "testConfigSpace",
			"dynamicDiscovery": true,
			"creationTime": "2023-06-12T10:00:00Z",
			"modifiedBy": "John Doe",
			"modifiedTime": "2023-06-12T10:00:00Z",
			"appConnectorGroups": [
				{
					"cityCountry": "City, Country",
					"countryCode": "CC",
					"creationTime": "2023-06-12T10:00:00Z",
					"description": "Test group",
					"dnsQueryType": "A",
					"enabled": true,
					"geoLocationId": "locationID",
					"id": "groupID",
					"latitude": "0.0",
					"location": "Test location",
					"longitude": "0.0",
					"modifiedBy": "John Doe",
					"modifiedTime": "2023-06-12T10:00:00Z",
					"name": "Test Group",
					"siemAppConnectorGroup": false,
					"upgradeDay": "Saturday",
					"upgradeTimeInSecs": "7200",
					"versionProfileId": "profileID",
					"serverGroups": [
						{
							"configSpace": "testConfigSpace",
							"creationTime": "2023-06-12T10:00:00Z",
							"description": "Test group",
							"enabled": true,
							"id": "groupID",
							"dynamicDiscovery": true,
							"modifiedBy": "John Doe",
							"modifiedTime": "2023-06-12T10:00:00Z",
							"name": "Test Group"
						}
					],
					"connectors": [
						{
							"applicationStartTime": "2023-06-12T10:00:00Z",
							"appConnectorGroupId": "groupID",
							"appConnectorGroupName": "Test Group",
							"controlChannelStatus": "Active",
							"creationTime": "2023-06-12T10:00:00Z",
							"ctrlBrokerName": "Test Broker",
							"currentVersion": "1.0",
							"description": "Test connector",
							"enabled": true,
							"expectedUpgradeTime": "2023-06-12T10:00:00Z",
							"expectedVersion": "2.0",
							"fingerprint": "Test fingerprint",
							"id": "connectorID",
							"ipAcl": [
								"0.0.0.0/0"
							],
							"issuedCertId": "certID",
							"lastBrokerConnectTime": "2023-06-12T10:00:00Z",
							"lastBrokerDisconnectTime": "2023-06-12T10:00:00Z",
							"lastUpgradeTime": "2023-06-12T10:00:00Z",
							"latitude": 0.0,
							"location": "Test location",
							"longitude": 0.0,
							"modifiedBy": "John Doe",
							"modifiedTime": "2023-06-12T10:00:00Z",
							"name": "Test Connector",
							"platform": "Test Platform",
							"previousVersion": "1.0",
							"privateIp": "10.0.0.1",
							"publicIp": "192.168.0.1",
							"signingCert": {},
							"upgradeAttempt": "Test attempt",
							"upgradeStatus": "Success"
						}
					]
				}
			],
			"servers": [
				{
					"address": "192.168.0.1",
					"appServerGroupIds": ["groupID"],
					"configSpace": "testConfigSpace",
					"creationTime": "2023-06-12T10:00:00Z",
					"description": "Test server",
					"enabled": true,
					"id": "serverID",
					"modifiedBy": "John Doe",
					"modifiedTime": "2023-06-12T10:00:00Z",
					"name": "Test Server"
				}
			],
			"applications": [
				{
					"id": "appID",
					"name": "Test Application"
				}
			]
		}`
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	})

	// Create a service with the client
	service := &servergroup.Service{Client: client}

	// Call the Get method
	groupID := "groupID"
	serverGroup, resp, err := service.Get(groupID)
	// Check the error
	if err != nil {
		t.Errorf("Error calling Get: %s", err)
	}

	// Check the response
	expectedResponse := &servergroup.ServerGroup{
		ID:               "groupID",
		Enabled:          true,
		Name:             "Test Group",
		Description:      "Test description",
		IpAnchored:       true,
		ConfigSpace:      "testConfigSpace",
		DynamicDiscovery: true,
		CreationTime:     "2023-06-12T10:00:00Z",
		ModifiedBy:       "John Doe",
		ModifiedTime:     "2023-06-12T10:00:00Z",
		AppConnectorGroups: []servergroup.AppConnectorGroups{
			{
				Citycountry:           "City, Country",
				CountryCode:           "CC",
				CreationTime:          "2023-06-12T10:00:00Z",
				Description:           "Test group",
				DnsqueryType:          "A",
				Enabled:               true,
				GeolocationID:         "locationID",
				ID:                    "groupID",
				Latitude:              "0.0",
				Location:              "Test location",
				Longitude:             "0.0",
				ModifiedBy:            "John Doe",
				ModifiedTime:          "2023-06-12T10:00:00Z",
				Name:                  "Test Group",
				SiemAppconnectorGroup: false,
				UpgradeDay:            "Saturday",
				UpgradeTimeinSecs:     "7200",
				VersionProfileID:      "profileID",
				AppServerGroups: []servergroup.AppServerGroups{
					{
						ConfigSpace:      "testConfigSpace",
						CreationTime:     "2023-06-12T10:00:00Z",
						Description:      "Test group",
						Enabled:          true,
						ID:               "groupID",
						DynamicDiscovery: true,
						ModifiedBy:       "John Doe",
						ModifiedTime:     "2023-06-12T10:00:00Z",
						Name:             "Test Group",
					},
				},
				Connectors: []servergroup.Connectors{
					{
						ApplicationStartTime:     "2023-06-12T10:00:00Z",
						AppConnectorGroupID:      "groupID",
						AppConnectorGroupName:    "Test Group",
						ControlChannelStatus:     "Active",
						CreationTime:             "2023-06-12T10:00:00Z",
						CtrlBrokerName:           "Test Broker",
						CurrentVersion:           "1.0",
						Description:              "Test connector",
						Enabled:                  true,
						ExpectedUpgradeTime:      "2023-06-12T10:00:00Z",
						ExpectedVersion:          "2.0",
						Fingerprint:              "Test fingerprint",
						ID:                       "connectorID",
						IPACL:                    []string{"0.0.0.0/0"},
						IssuedCertID:             "certID",
						LastBrokerConnecttime:    "2023-06-12T10:00:00Z",
						LastBrokerDisconnectTime: "2023-06-12T10:00:00Z",
						LastUpgradeTime:          "2023-06-12T10:00:00Z",
						Latitude:                 0.0,
						Location:                 "Test location",
						Longitude:                0.0,
						ModifiedBy:               "John Doe",
						ModifiedTime:             "2023-06-12T10:00:00Z",
						Name:                     "Test Connector",
						Platform:                 "Test Platform",
						PreviousVersion:          "1.0",
						PrivateIP:                "10.0.0.1",
						PublicIP:                 "192.168.0.1",
						SigningCert:              map[string]interface{}{},
						UpgradeAttempt:           "Test attempt",
						UpgradeStatus:            "Success",
					},
				},
			},
		},
		Servers: []servergroup.ApplicationServer{
			{
				Address:           "192.168.0.1",
				AppServerGroupIds: []string{"groupID"},
				ConfigSpace:       "testConfigSpace",
				CreationTime:      "2023-06-12T10:00:00Z",
				Description:       "Test server",
				Enabled:           true,
				ID:                "serverID",
				ModifiedBy:        "John Doe",
				ModifiedTime:      "2023-06-12T10:00:00Z",
				Name:              "Test Server",
			},
		},
		Applications: []servergroup.Applications{
			{
				ID:   "appID",
				Name: "Test Application",
			},
		},
	}

	if !reflect.DeepEqual(serverGroup, expectedResponse) {
		t.Errorf("Expected response:\n%+v\n\nGot:\n%+v", expectedResponse, serverGroup)
	}

	// Check the response status code
	expectedStatusCode := http.StatusOK
	if resp.StatusCode != expectedStatusCode {
		t.Errorf("Expected status code: %d, Got: %d", expectedStatusCode, resp.StatusCode)
	}
}

/*
	func TestServerGroup_GetByName(t *testing.T) {
		client, mux, server := tests.NewZpaClientMock()
		defer server.Close()
		mux.HandleFunc("/mgmtconfig/v1/admin/customers/customerid/serverGroup", func(w http.ResponseWriter, r *http.Request) {
			// Get the query parameter "name" from the request
			query := r.URL.Query()
			name := query.Get("search")

			// Check if the name matches the expected value
			if name == "Group1" {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{
					"list": [
						{
							"appConnectorGroups": [],
							"applications": [],
							"configSpace": "",
							"creationTime": "",
							"description": "",
							"enabled": true,
							"id": "123",
							"modifiedBy": "",
							"modifiedTime": "",
							"name": "Group1",
							"servers": []
						}
					],
					"totalPages": 1
				}`))
			} else {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(`{"error": "Server group not found"}`))
			}
		})
		service := &servergroup.Service{
			Client: client,
		}

		// Make the GetByName request
		group, _, err := service.GetByName("Group1")
		// Check if the request was successful
		if err != nil {
			t.Errorf("Error making GetByName request: %v", err)
		}

		// Check if the group ID and name match the expected values
		if group.ID != "123" {
			t.Errorf("Expected group ID '123', but got '%s'", group.ID)
		}
		if group.Name != "Group1" {
			t.Errorf("Expected group name 'Group1', but got '%s'", group.Name)
		}
	}
*/
func TestServerGroup_Create(t *testing.T) {
	client, mux, server := tests.NewZpaClientMock()
	defer server.Close()
	mux.HandleFunc("/mgmtconfig/v1/admin/customers/customerid/serverGroup", func(w http.ResponseWriter, r *http.Request) {
		// Write a JSON response
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"appConnectorGroups": [],
			"applications": [],
			"configSpace": "",
			"creationTime": "",
			"description": "",
			"enabled": true,
			"id": "123",
			"modifiedBy": "",
			"modifiedTime": "",
			"name": "Group1",
			"servers": []
		}`))
	})

	service := &servergroup.Service{
		Client: client,
	}

	// Create a sample group
	group := &servergroup.ServerGroup{
		Name:    "Group1",
		Enabled: true,
	}

	// Make the Create request
	createdGroup, _, err := service.Create(group)
	// Check if the request was successful
	if err != nil {
		t.Errorf("Error making Create request: %v", err)
	}

	// Check if the created group ID and name match the expected values
	if createdGroup.ID != "123" {
		t.Errorf("Expected created group ID '123', but got '%s'", createdGroup.ID)
	}
	if createdGroup.Name != "Group1" {
		t.Errorf("Expected created group name 'Group1', but got '%s'", createdGroup.Name)
	}
}

func TestServerGroup_Update(t *testing.T) {
	client, mux, server := tests.NewZpaClientMock()
	defer server.Close()
	mux.HandleFunc("/mgmtconfig/v1/admin/customers/customerid/serverGroup/123", func(w http.ResponseWriter, r *http.Request) {
		// Write a JSON response
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"appConnectorGroups": [],
			"applications": [],
			"configSpace": "",
			"creationTime": "",
			"description": "Updated description",
			"enabled": true,
			"id": "123",
			"modifiedBy": "",
			"modifiedTime": "",
			"name": "Group1",
			"servers": []
		}`))
	})

	service := &servergroup.Service{
		Client: client,
	}

	// Create a sample group with updated description
	group := &servergroup.ServerGroup{
		ID:          "123",
		Name:        "Group1",
		Description: "Updated description",
		Enabled:     true,
	}

	// Make the Update request
	_, err := service.Update("123", group)
	// Check if the request was successful
	if err != nil {
		t.Errorf("Error making Update request: %v", err)
	}
}

func TestServerGroup_Delete(t *testing.T) {
	client, mux, server := tests.NewZpaClientMock()
	defer server.Close()
	mux.HandleFunc("/mgmtconfig/v1/admin/customers/customerid/serverGroup/123", func(w http.ResponseWriter, r *http.Request) {
		// Write a JSON response
		w.WriteHeader(http.StatusNoContent)
	})

	service := &servergroup.Service{
		Client: client,
	}

	// Make the Delete request
	_, err := service.Delete("123")
	// Check if the request was successful
	if err != nil {
		t.Errorf("Error making Delete request: %v", err)
	}
}

func TestServerGroup_GetAll(t *testing.T) {
	client, mux, server := tests.NewZpaClientMock()
	defer server.Close()
	mux.HandleFunc("/mgmtconfig/v1/admin/customers/customerid/serverGroup", func(w http.ResponseWriter, r *http.Request) {
		// Write a JSON response with an array of server groups
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"list":[
				{
					"appConnectorGroups": [],
					"applications": [],
					"configSpace": "",
					"creationTime": "",
					"description": "",
					"enabled": true,
					"id": "123",
					"modifiedBy": "",
					"modifiedTime": "",
					"name": "Group1",
					"servers": []
				},
				{
					"appConnectorGroups": [],
					"applications": [],
					"configSpace": "",
					"creationTime": "",
					"description": "",
					"enabled": true,
					"id": "456",
					"modifiedBy": "",
					"modifiedTime": "",
					"name": "Group2",
					"servers": []
				}
			],
			"totalPages": 1
		}`))
	})

	service := &servergroup.Service{
		Client: client,
	}

	// Make the GetAll request
	groups, _, err := service.GetAll()
	// Check if the request was successful
	if err != nil {
		t.Errorf("Error making GetAll request: %v", err)
	}

	// Check the number of returned groups
	if len(groups) != 2 {
		t.Errorf("Expected 2 groups, but got %d", len(groups))
	}

	// Check the ID and name of each group
	if groups[0].ID != "123" || groups[0].Name != "Group1" {
		t.Errorf("Expected group ID '123' and name 'Group1', but got ID '%s' and name '%s'", groups[0].ID, groups[0].Name)
	}
	if groups[1].ID != "456" || groups[1].Name != "Group2" {
		t.Errorf("Expected group ID '456' and name 'Group2', but got ID '%s' and name '%s'", groups[1].ID, groups[1].Name)
	}
}
