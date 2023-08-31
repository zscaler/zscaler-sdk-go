package user_authentication_settings

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/zscaler/zscaler-sdk-go/tests"
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

	client, err := tests.NewZiaClient()
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}
	service := New(client)
	resources, err := service.Get()
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
		_, err := service.Client.Create(fmt.Sprintf("%s?action=REMOVE_FROM_LIST", exemptedUrlsEndpoint), ExemptedUrls{urlsToRemove})
		if err != nil {
			log.Printf("Error removing exempted URL during cleanup: %v", err)
		}
	}
}

func TestUserAuthenticationSettings(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	service := &Service{Client: client}

	// Create 3 random exempted URLs
	initialUrls, err := service.Get()
	if err != nil {
		t.Fatalf("Error fetching initial exempted URLs: %v", err)
	}

	newUrls := generateRandomUrls(3)
	allUrls := append(initialUrls.URLs, newUrls...)

	// Update with the new exempted URLs
	_, err = service.Update(ExemptedUrls{URLs: allUrls})
	if err != nil {
		t.Fatalf("Error updating exempted URLs: %v", err)
	}

	// Fetch and validate the exempted URLs after updating
	updatedUrls, err := service.Get()
	if err != nil {
		t.Fatalf("Error fetching updated exempted URLs: %v", err)
	}
	for _, url := range newUrls {
		if !contains(updatedUrls.URLs, url) {
			t.Errorf("URL %v was not updated properly", url)
		}
	}

	// Clean up by removing the URLs we added for testing
	_, err = service.Update(ExemptedUrls{URLs: initialUrls.URLs})
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
