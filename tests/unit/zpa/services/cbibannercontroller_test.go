// Package unit provides unit tests for ZPA CBI Banner Controller service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// CBIBannerController represents the CBI banner for testing
type CBIBannerController struct {
	ID                string `json:"id,omitempty"`
	Name              string `json:"name,omitempty"`
	PrimaryColor      string `json:"primaryColor,omitempty"`
	TextColor         string `json:"textColor,omitempty"`
	NotificationTitle string `json:"notificationTitle,omitempty"`
	NotificationText  string `json:"notificationText,omitempty"`
	Logo              string `json:"logo,omitempty"`
	Banner            bool   `json:"banner,omitempty"`
	IsDefault         bool   `json:"isDefault,omitempty"`
	Persist           bool   `json:"persist,omitempty"`
}

// TestCBIBanner_Structure tests the struct definitions
func TestCBIBanner_Structure(t *testing.T) {
	t.Parallel()

	t.Run("CBIBannerController JSON marshaling", func(t *testing.T) {
		banner := CBIBannerController{
			ID:                "banner-123",
			Name:              "Corporate Banner",
			PrimaryColor:      "#0066CC",
			TextColor:         "#FFFFFF",
			NotificationTitle: "Secure Browsing",
			NotificationText:  "You are now browsing securely through Cloud Browser Isolation",
			Logo:              "base64encodedlogo==",
			Banner:            true,
			IsDefault:         false,
			Persist:           true,
		}

		data, err := json.Marshal(banner)
		require.NoError(t, err)

		var unmarshaled CBIBannerController
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, banner.ID, unmarshaled.ID)
		assert.Equal(t, banner.Name, unmarshaled.Name)
		assert.Equal(t, banner.PrimaryColor, unmarshaled.PrimaryColor)
		assert.True(t, unmarshaled.Banner)
		assert.True(t, unmarshaled.Persist)
	})

	t.Run("CBIBannerController from API response", func(t *testing.T) {
		apiResponse := `{
			"id": "banner-456",
			"name": "Default Isolation Banner",
			"primaryColor": "#003366",
			"textColor": "#FFFFFF",
			"notificationTitle": "Browser Isolation Active",
			"notificationText": "This session is protected by browser isolation",
			"logo": "aW1hZ2VkYXRh",
			"banner": true,
			"isDefault": true,
			"persist": false
		}`

		var banner CBIBannerController
		err := json.Unmarshal([]byte(apiResponse), &banner)
		require.NoError(t, err)

		assert.Equal(t, "banner-456", banner.ID)
		assert.Equal(t, "#003366", banner.PrimaryColor)
		assert.True(t, banner.IsDefault)
		assert.True(t, banner.Banner)
	})
}

// TestCBIBanner_ResponseParsing tests parsing of API responses
func TestCBIBanner_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse banner list response", func(t *testing.T) {
		response := `[
			{"id": "1", "name": "Banner 1", "isDefault": true, "banner": true},
			{"id": "2", "name": "Banner 2", "isDefault": false, "banner": true},
			{"id": "3", "name": "Banner 3", "isDefault": false, "banner": false}
		]`

		var banners []CBIBannerController
		err := json.Unmarshal([]byte(response), &banners)
		require.NoError(t, err)

		assert.Len(t, banners, 3)
		assert.True(t, banners[0].IsDefault)
	})
}

// TestCBIBanner_MockServerOperations tests CRUD operations
func TestCBIBanner_MockServerOperations(t *testing.T) {
	t.Run("GET banner by ID", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Contains(t, r.URL.Path, "/banners/")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{"id": "banner-123", "name": "Mock Banner", "banner": true}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/cbi/api/customers/123/banners/banner-123")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("GET all banners", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `[{"id": "1", "name": "Banner A"}, {"id": "2", "name": "Banner B"}]`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/cbi/api/customers/123/banners")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("POST create banner", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{"id": "new-banner", "name": "New Banner"}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Post(server.URL+"/cbi/api/customers/123/banner", "application/json", nil)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("PUT update banner", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		req, _ := http.NewRequest("PUT", server.URL+"/banners/banner-123", nil)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("DELETE banner", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "DELETE", r.Method)
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		req, _ := http.NewRequest("DELETE", server.URL+"/banners/banner-123", nil)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})
}

// TestCBIBanner_SpecialCases tests edge cases
func TestCBIBanner_SpecialCases(t *testing.T) {
	t.Parallel()

	t.Run("Banner color schemes", func(t *testing.T) {
		colors := []struct {
			primary string
			text    string
		}{
			{"#0066CC", "#FFFFFF"},
			{"#FF0000", "#000000"},
			{"#00FF00", "#333333"},
		}

		for _, c := range colors {
			banner := CBIBannerController{
				ID:           "banner-color",
				Name:         "Color Test",
				PrimaryColor: c.primary,
				TextColor:    c.text,
			}

			data, err := json.Marshal(banner)
			require.NoError(t, err)

			assert.Contains(t, string(data), c.primary)
			assert.Contains(t, string(data), c.text)
		}
	})

	t.Run("Default banner", func(t *testing.T) {
		banner := CBIBannerController{
			ID:        "banner-default",
			Name:      "Default Banner",
			IsDefault: true,
		}

		data, err := json.Marshal(banner)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"isDefault":true`)
	})
}

