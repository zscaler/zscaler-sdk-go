package region

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v2/tests"
)

func TestRegionSearch(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	service := New(client)

	// Define test cases with different prefix options
	testCases := []struct {
		name   string
		prefix string
	}{
		{
			name:   "TestSearchByContinentCode",
			prefix: "NA",
		},
		{
			name:   "TestSearchByCountryCode",
			prefix: "US",
		},
		{
			name:   "TestSearchByStateName",
			prefix: "California",
		},
		{
			name:   "TestSearchByCityName",
			prefix: "San Jose",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			regions, err := service.GetDatacenterRegion(tc.prefix)
			if err != nil {
				t.Errorf("Error in %v: %v", tc.name, err)
				return
			}
			if len(regions) == 0 {
				t.Errorf("No regions found for prefix %v in %v", tc.prefix, tc.name)
				return
			}
		})
	}
}
