package ipv6_config

/*

func TestIPv6Config(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	ctx := context.Background()

	// Step 1: Test GetDns64Prefix with no search
	dnsPrefixes, err := GetDns64Prefix(ctx, service)
	if err != nil {
		t.Fatalf("Error fetching DNS64 prefixes: %v", err)
	}
	if len(dnsPrefixes) == 0 {
		t.Fatal("Expected at least one DNS64 prefix, but got none")
	}
	t.Logf("Retrieved %d DNS64 prefixes", len(dnsPrefixes))

	// Step 2: Test GetNat64Prefix with no search
	natPrefixes, err := GetNat64Prefix(ctx, service)
	if err != nil {
		t.Fatalf("Error fetching NAT64 prefixes: %v", err)
	}
	if len(natPrefixes) == 0 {
		t.Fatal("Expected at least one NAT64 prefix, but got none")
	}
	t.Logf("Retrieved %d NAT64 prefixes", len(natPrefixes))

	// Step 3: Use a known name from DNS64 prefix to search
	searchTerm := dnsPrefixes[0].Name
	searchedDnsPrefixes, err := GetDns64Prefix(ctx, service, searchTerm)
	if err != nil {
		t.Fatalf("Error searching DNS64 prefixes with term '%s': %v", searchTerm, err)
	}
	if len(searchedDnsPrefixes) == 0 {
		t.Errorf("Expected at least one DNS64 prefix result for search '%s', but got none", searchTerm)
	}

	// Validate match
	found := false
	for _, prefix := range searchedDnsPrefixes {
		if prefix.Name == searchTerm {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Search term '%s' not found in DNS64 search results", searchTerm)
	}

	// Step 4: Fetch full IPv6 Config
	config, err := GetIPv6Config(ctx, service)
	if err != nil {
		t.Fatalf("Error fetching full IPv6 config: %v", err)
	}
	t.Logf("IPv6 Enabled: %v", config.IpV6Enabled)
	t.Logf("DNS Prefix: %s", config.DnsPrefix)
	t.Logf("NAT Prefixes Count: %d", len(config.NatPrefixes))
}
*/
