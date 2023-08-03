package filteringrules

/*
func TestFirewallFilteringRule(t *testing.T) {
	name := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	updateName := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}
	// create ip source group for testing
	sourceIPGroupService := ipsourcegroups.New(client)
	sourceIPGroup, err := sourceIPGroupService.Create(&ipsourcegroups.IPSourceGroups{
		Name:        name,
		Description: name,
		IPAddresses: []string{"192.168.1.1", "192.168.1.2", "192.168.1.3"},
	})
	// Check if the request was successful
	if err != nil {
		t.Errorf("Error creating source ip group for testing server group: %v", err)
	}
	defer func() {
		_, err := sourceIPGroupService.Delete(sourceIPGroup.ID)
		if err != nil {
			t.Errorf("Error deleting source ip group: %v", err)
		}
	}()

	// create ip destination group for testing
	dstIPGroupService := ipdestinationgroups.New(client)
	dstIPGroup, err := dstIPGroupService.Create(&ipdestinationgroups.IPDestinationGroups{
		Name:        name,
		Description: name,
		Type:        "DSTN_FQDN",
		Addresses:   []string{"test1.acme.com", "test2.acme.com", "test3.acme.com"},
	})
	// Check if the request was successful
	if err != nil {
		t.Errorf("Error creating ip destination group for testing server group: %v", err)
	}
	defer func() {
		_, err := dstIPGroupService.Delete(dstIPGroup.ID)
		if err != nil {
			t.Errorf("Error deleting ip destination group: %v", err)
		}
	}()

	// create rule label for testing
	ruleLabelService := rule_labels.New(client)
	ruleLabel, _, err := ruleLabelService.Create(&rule_labels.RuleLabels{
		Name:        name,
		Description: name,
	})
	// Check if the request was successful
	if err != nil {
		t.Errorf("Error creating rule label for testing server group: %v", err)
	}
	defer func() {
		_, err := ruleLabelService.Delete(ruleLabel.ID)
		if err != nil {
			t.Errorf("Error deleting rule label: %v", err)
		}
	}()
	service := New(client)
	rule := FirewallFilteringRules{
		Name:           name,
		Description:    name,
		Order:          6,
		Rank:           7,
		Action:         "ALLOW",
		DestCountries:  []string{"COUNTRY_CA", "COUNTRY_US", "COUNTRY_MX", "COUNTRY_AU", "COUNTRY_GB"},
		NwApplications: []string{"APNS", "GARP", "PERFORCE", "WINDOWS_MARKETPLACE", "DIAMETER"},
		SrcIpGroups: []common.IDNameExtensions{
			{
				ID: sourceIPGroup.ID,
			},
		},
		DestIpGroups: []common.IDNameExtensions{
			{
				ID: dstIPGroup.ID,
			},
		},
		Labels: []common.IDNameExtensions{
			{
				ID: ruleLabel.ID,
			},
		},
	}

	// Test resource creation
	createdResource, err := service.Create(&rule)
	// Check if the request was successful
	if err != nil {
		t.Errorf("Error making POST request: %v", err)
	}

	if createdResource.ID == 0 {
		t.Error("Expected created resource ID to be non-empty, but got ''")
	}
	if createdResource.Name != name {
		t.Errorf("Expected created resource name '%s', but got '%s'", name, createdResource.Name)
	}
	// Test resource retrieval
	retrievedResource, err := service.Get(createdResource.ID)
	if err != nil {
		t.Errorf("Error retrieving resource: %v", err)
	}
	if retrievedResource.ID != createdResource.ID {
		t.Errorf("Expected retrieved resource ID '%d', but got '%d'", createdResource.ID, retrievedResource.ID)
	}
	if retrievedResource.Name != name {
		t.Errorf("Expected retrieved resource name '%s', but got '%s'", name, createdResource.Name)
	}
	// Test resource update
	retrievedResource.Name = updateName
	_, err = service.Update(createdResource.ID, retrievedResource)
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
	if updatedResource.Name != updateName {
		t.Errorf("Expected retrieved updated resource name '%s', but got '%s'", updateName, updatedResource.Name)
	}

	// Test resource retrieval by name
	retrievedResource, err = service.GetByName(updateName)
	if err != nil {
		t.Errorf("Error retrieving resource by name: %v", err)
	}
	if retrievedResource.ID != createdResource.ID {
		t.Errorf("Expected retrieved resource ID '%d', but got '%d'", createdResource.ID, retrievedResource.ID)
	}
	if retrievedResource.Name != updateName {
		t.Errorf("Expected retrieved resource name '%s', but got '%s'", updateName, createdResource.Name)
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
	_, err = service.Delete(createdResource.ID)
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
*/
