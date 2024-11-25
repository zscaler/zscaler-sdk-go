package user_authentication_settings

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/tests"
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup() {
	cleanResources()
}

func teardown() {
	cleanResources()
}

func shouldClean() bool {
	val, present := os.LookupEnv("ZSCALER_SDK_TEST_SWEEP")
	return !present || (present && (val == "" || val == "true"))
}

func cleanResources() {
	if !shouldClean() {
		return
	}

	// Define the context here
	ctx := context.Background()

	service, err := tests.NewOneAPIClient()
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	resources, err := Get(ctx, service) // Use ctx in the Get call
	if err != nil {
		log.Printf("Error retrieving exempted URLs during cleanup: %v", err)
		return
	}

	// Identify and remove URLs starting with "site"
	var urlsToRemove []string
	for _, url := range resources.URLs {
		if strings.HasPrefix(url, "site") {
			urlsToRemove = append(urlsToRemove, url)
		}
	}
	if len(urlsToRemove) > 0 {
		// Use ctx in the Create call
		_, err := service.Client.Create(ctx, fmt.Sprintf("%s?action=REMOVE_FROM_LIST", exemptedUrlsEndpoint), ExemptedUrls{urlsToRemove})
		if err != nil {
			log.Printf("Error removing exempted URL during cleanup: %v", err)
		}
	}
}

func TestUserAuthenticationSettings(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	// Create 3 random exempted URLs
	initialUrls, err := Get(context.Background(), service)
	if err != nil {
		t.Fatalf("Error fetching initial exempted URLs: %v", err)
	}

	newUrls := generateRandomUrls(3)
	allUrls := append(initialUrls.URLs, newUrls...)

	// Update with the new exempted URLs
	_, err = Update(context.Background(), service, ExemptedUrls{URLs: allUrls})
	if err != nil {
		t.Fatalf("Error updating exempted URLs: %v", err)
	}

	// Fetch and validate the exempted URLs after updating
	updatedUrls, err := Get(context.Background(), service)
	if err != nil {
		t.Fatalf("Error fetching updated exempted URLs: %v", err)
	}
	for _, url := range newUrls {
		if !contains(updatedUrls.URLs, url) {
			t.Errorf("URL %v was not updated properly", url)
		}
	}

	// Clean up by removing the URLs we added for testing
	_, err = Update(context.Background(), service, ExemptedUrls{URLs: initialUrls.URLs})
	if err != nil {
		t.Fatalf("Error cleaning up exempted URLs: %v", err)
	}
}

func generateRandomUrls(count int) []string {
	var urls []string
	domains := []string{".example.com", ".test.com"}

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < count; i++ {
		url := fmt.Sprintf("%s%d%s", "site", rand.Intn(1000), domains[rand.Intn(len(domains))])
		urls = append(urls, url)
	}
	return urls
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
