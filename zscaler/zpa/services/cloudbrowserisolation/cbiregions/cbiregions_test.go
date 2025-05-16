package cbiregions

/*
func TestGetAllRegions(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	// service, err := tests.NewZPAClient()
	// if err != nil {
	// 	t.Fatalf("Error creating client: %v", err)
	// }

	// 1. First GetAll regions and ensure a response is returned.
	regions, resp, err := GetAll(context.Background(), service)
	if err != nil || resp.StatusCode >= 400 || len(regions) == 0 {
		t.Fatalf("Failed to fetch regions: %v", err)
	}

	// To ensure that subsequent tests don't fail due to an empty regions list
	if len(regions) == 0 {
		t.Fatal("No regions returned. Can't proceed with further tests.")
		return
	}

	// 3. Test the GetByName method by querying the Name of any of the returned regions from GetAll.
	firstRegionName := regions[0].Name
	singleRegionByName, resp, err := GetByName(context.Background(), service, firstRegionName)
	if err != nil || resp.StatusCode >= 400 || singleRegionByName == nil {
		t.Errorf("Failed to fetch region by Name %s: %v", firstRegionName, err)
	} else if singleRegionByName.Name != firstRegionName {
		t.Errorf("Mismatch in region Name. Expected %s, got %s", firstRegionName, singleRegionByName.Name)
	}
}

func TestCaseSensitivityOfGetByName(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	// service, err := tests.NewZPAClient()
	// if err != nil {
	// 	t.Fatalf("Error creating client: %v", err)
	// }

	requiredNames := []string{"Frankfurt", "Ireland", "Washington", "Singapore"}

	for _, knownName := range requiredNames {
		variations := []string{
			strings.ToUpper(knownName),
			strings.ToLower(knownName),
			cases.Title(language.English).String(knownName),
		}

		found := false

		for _, variation := range variations {
			t.Run(fmt.Sprintf("GetByName case sensitivity test for %s", variation), func(t *testing.T) {
				t.Logf("Attempting to retrieve region with name variation: %s", variation)
				region, _, err := GetByName(context.Background(), service, variation)

				if err != nil {
					if strings.Contains(err.Error(), "no region named") {
						t.Logf("Region with name variation '%s' not found: %v", variation, err)
						return
					}
					t.Errorf("Error getting region with name variation '%s': %v", variation, err)
				} else {
					found = true
					// Check if the region's actual name matches the known name
					if region.Name != knownName {
						t.Errorf("Expected region name to be '%s' for variation '%s', but got '%s'", knownName, variation, region.Name)
					}
				}
			})

			if found {
				break
			}
		}

		if !found {
			t.Logf("Region '%s' was not found in any variation, but moving on to next region", knownName)
		}
	}
}

func TestGetByNameNonExistentResource(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	// service, err := tests.NewZPAClient()
	// if err != nil {
	// 	t.Fatalf("Error creating client: %v", err)
	// }

	_, _, err = GetByName(context.Background(), service, "non_existent_name")
	if err == nil {
		t.Error("Expected error retrieving resource by non-existent name, but got nil")
	}
}
*/
