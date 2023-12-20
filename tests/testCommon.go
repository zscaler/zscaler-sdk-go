package tests

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/zscaler/zscaler-sdk-go/v2/logger"
	"github.com/zscaler/zscaler-sdk-go/v2/zcon"
	"github.com/zscaler/zscaler-sdk-go/v2/zdx"
	"github.com/zscaler/zscaler-sdk-go/v2/zia"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa"
)

const (
	charSetAlphaUpper  = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	charSetAlphaLower  = "abcdefghijklmnopqrstuvwxyz"
	charSetNumeric     = "0123456789"
	charSetSpecialChar = "!@#$%^&*"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func TestPassword(length int) string {
	if length < 8 {
		length = 8
	} else if length > 100 {
		length = 100
	}

	result := make([]byte, length)
	result[0] = charSetAlphaLower[rand.Intn(len(charSetAlphaLower))]
	result[1] = charSetAlphaUpper[rand.Intn(len(charSetAlphaUpper))]
	result[2] = charSetNumeric[rand.Intn(len(charSetNumeric))]
	result[3] = charSetSpecialChar[rand.Intn(len(charSetSpecialChar))]

	charSetAll := charSetAlphaLower + charSetAlphaUpper + charSetNumeric + charSetSpecialChar
	for i := 4; i < length; i++ {
		result[i] = charSetAll[rand.Intn(len(charSetAll))]
	}
	// Shuffle the result to avoid predictable patterns (lower, upper, numeric, special)
	rand.Shuffle(len(result), func(i, j int) {
		result[i], result[j] = result[j], result[i]
	})

	return string(result)
}

func NewZpaClient() (*zpa.Client, error) {
	zpa_client_id := os.Getenv("ZPA_CLIENT_ID")
	zpa_client_secret := os.Getenv("ZPA_CLIENT_SECRET")
	zpa_customer_id := os.Getenv("ZPA_CUSTOMER_ID")
	zpa_cloud := os.Getenv("ZPA_CLOUD")
	config, err := zpa.NewConfig(zpa_client_id, zpa_client_secret, zpa_customer_id, zpa_cloud, "zscaler-sdk-go")
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
			Logger:       logger.NewNopLogger(),
			BaseURL:      serverURL,
		},
	}
	return client, mux, server
}

func NewZiaClient() (*zia.Client, error) {
	username := os.Getenv("ZIA_USERNAME")
	password := os.Getenv("ZIA_PASSWORD")
	apiKey := os.Getenv("ZIA_API_KEY")
	ziaCloud := os.Getenv("ZIA_CLOUD")

	cli, err := zia.NewClient(username, password, apiKey, ziaCloud, "zscaler-sdk-go")
	if err != nil {
		log.Printf("[ERROR] creating client failed: %v\n", err)
		return nil, err
	}
	return cli, nil
}

func NewZConClient() (*zcon.Client, error) {
	username := os.Getenv("ZCON_USERNAME")
	password := os.Getenv("ZCON_PASSWORD")
	apiKey := os.Getenv("ZCON_API_KEY")
	zconCloud := os.Getenv("ZCON_CLOUD")

	cli, err := zcon.NewClient(username, password, apiKey, zconCloud, "zscaler-sdk-go")
	if err != nil {
		log.Printf("[ERROR] creating client failed: %v\n", err)
		return nil, err
	}
	return cli, nil
}

func NewZdxClient() (*zdx.Client, error) {
	apiKeyID := os.Getenv("ZDX_API_KEY_ID")
	apiSecret := os.Getenv("ZDX_API_SECRET")

	config, err := zdx.NewConfig(apiKeyID, apiSecret, "zscaler-sdk-go")
	if err != nil {
		log.Printf("[ERROR] creating config failed: %v\n", err)
		return nil, err
	}
	zdxClient := zdx.NewClient(config)
	return zdxClient, nil
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
