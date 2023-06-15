package tests

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/zpa"
)

func NewZpaClient() (*zpa.Client, error) {
	zpa_client_id := os.Getenv("ZPA_CLIENT_ID")
	zpa_client_secret := os.Getenv("ZPA_CLIENT_SECRET")
	zpa_customer_id := os.Getenv("ZPA_CUSTOMER_ID")
	zpa_cloud := os.Getenv("ZPA_CLOUD")
	config, err := zpa.NewConfig(zpa_client_id, zpa_client_secret, zpa_customer_id, zpa_cloud, "testing")
	if err != nil {
		log.Printf("[ERROR] creating config failed: %v\n", err)
		return nil, err
	}
	zpaClient := zpa.NewClient(config)
	return zpaClient, nil
}

func NewZpaClientMock() (*zpa.Client, *http.ServeMux, *httptest.Server) {
	mux := http.NewServeMux()

	// Create a request handler for the exact endpoint
	mux.HandleFunc("/signin", func(w http.ResponseWriter, r *http.Request) {
		// Write a JSON response
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"token_type": "bearer", "access_token": "Group 1"}`))
	})

	// Create a test server using the ServeMux
	server := httptest.NewServer(mux)

	serverURL, _ := url.Parse(server.URL)
	// Create a client and set the base URL to the mock server URL
	client := &zpa.Client{
		Config: &zpa.Config{
			ClientID:     "clientid",
			ClientSecret: "clientsecret",
			CustomerID:   "customerid",
			// Logger:       logger.NewNopLogger(),
			BaseURL: serverURL,
		},
	}
	return client, mux, server
}

// ParseJSONRequest parses the JSON request body from the given HTTP request.
func ParseJSONRequest(t *testing.T, r *http.Request, v interface{}) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.Unmarshal(body, v)
}

// WriteJSONResponse writes the JSON response with the given status code and data to the HTTP response writer.
func WriteJSONResponse(t *testing.T, w http.ResponseWriter, statusCode int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if data != nil {
		encoder := json.NewEncoder(w)
		return encoder.Encode(data)
	}

	return nil
}
