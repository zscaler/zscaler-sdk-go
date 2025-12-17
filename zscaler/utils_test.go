package zscaler

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecodeJSON(t *testing.T) {
	t.Run("Decode valid JSON", func(t *testing.T) {
		type TestStruct struct {
			Name  string `json:"name"`
			Value int    `json:"value"`
		}

		jsonData := []byte(`{"name": "test", "value": 42}`)
		var result TestStruct

		err := decodeJSON(jsonData, &result)

		require.NoError(t, err)
		assert.Equal(t, "test", result.Name)
		assert.Equal(t, 42, result.Value)
	})

	t.Run("Decode JSON array", func(t *testing.T) {
		type Item struct {
			ID string `json:"id"`
		}

		jsonData := []byte(`[{"id": "1"}, {"id": "2"}, {"id": "3"}]`)
		var result []Item

		err := decodeJSON(jsonData, &result)

		require.NoError(t, err)
		assert.Len(t, result, 3)
		assert.Equal(t, "1", result[0].ID)
		assert.Equal(t, "2", result[1].ID)
		assert.Equal(t, "3", result[2].ID)
	})

	t.Run("Decode invalid JSON", func(t *testing.T) {
		type TestStruct struct {
			Name string `json:"name"`
		}

		jsonData := []byte(`{invalid json}`)
		var result TestStruct

		err := decodeJSON(jsonData, &result)

		require.Error(t, err)
	})

	t.Run("Decode empty JSON object", func(t *testing.T) {
		type TestStruct struct {
			Name string `json:"name"`
		}

		jsonData := []byte(`{}`)
		var result TestStruct

		err := decodeJSON(jsonData, &result)

		require.NoError(t, err)
		assert.Empty(t, result.Name)
	})

	t.Run("Decode JSON with nested objects", func(t *testing.T) {
		type Nested struct {
			Inner string `json:"inner"`
		}
		type TestStruct struct {
			Outer  string `json:"outer"`
			Nested Nested `json:"nested"`
		}

		jsonData := []byte(`{"outer": "value", "nested": {"inner": "innerValue"}}`)
		var result TestStruct

		err := decodeJSON(jsonData, &result)

		require.NoError(t, err)
		assert.Equal(t, "value", result.Outer)
		assert.Equal(t, "innerValue", result.Nested.Inner)
	})
}

func TestUnescapeHTML(t *testing.T) {
	t.Run("Unescape name field", func(t *testing.T) {
		entity := map[string]interface{}{
			"name":        "Test &amp; Name",
			"description": "Regular description",
		}

		unescapeHTML(&entity)

		// Note: The function reads from JSON, so we need to verify via marshal/unmarshal
		data, _ := json.Marshal(entity)
		var result map[string]interface{}
		_ = json.Unmarshal(data, &result)

		assert.Equal(t, "Test & Name", result["name"])
	})

	t.Run("Unescape description field", func(t *testing.T) {
		entity := map[string]interface{}{
			"name":        "Regular name",
			"description": "Desc &lt;with&gt; HTML &amp; entities",
		}

		unescapeHTML(&entity)

		data, _ := json.Marshal(entity)
		var result map[string]interface{}
		_ = json.Unmarshal(data, &result)

		assert.Equal(t, "Desc <with> HTML & entities", result["description"])
	})

	t.Run("Nil entity does not panic", func(t *testing.T) {
		// Should not panic
		unescapeHTML(nil)
	})

	t.Run("Double escaped HTML entities", func(t *testing.T) {
		entity := map[string]interface{}{
			"name": "Test &amp;amp; Name",
		}

		unescapeHTML(&entity)

		data, _ := json.Marshal(entity)
		var result map[string]interface{}
		_ = json.Unmarshal(data, &result)

		// Double unescape should handle &amp;amp; -> &amp; -> &
		assert.Equal(t, "Test & Name", result["name"])
	})

	t.Run("Non-map entity with name field", func(t *testing.T) {
		type TestStruct struct {
			Name        string `json:"name"`
			Description string `json:"description"`
		}

		entity := &TestStruct{
			Name:        "Test &amp; Name",
			Description: "Desc &lt;tag&gt;",
		}

		unescapeHTML(entity)

		assert.Equal(t, "Test & Name", entity.Name)
		assert.Equal(t, "Desc <tag>", entity.Description)
	})
}

func TestRemoveOneApiEndpointPrefix(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Remove /zia prefix",
			input:    "/zia/api/v1/adminUsers",
			expected: "/api/v1/adminUsers",
		},
		{
			name:     "Remove /zpa prefix",
			input:    "/zpa/api/v1/policies",
			expected: "/api/v1/policies",
		},
		{
			name:     "Remove /zcc prefix",
			input:    "/zcc/admin/devices",
			expected: "/admin/devices",
		},
		{
			name:     "Remove /zcc/papi prefix",
			input:    "/zcc/papi/public/v1/getDevices",
			expected: "/public/v1/getDevices",
		},
		{
			name:     "Remove /zdx prefix",
			input:    "/zdx/v1/devices",
			expected: "/v1/devices",
		},
		{
			name:     "Remove /ztw/api/v1 prefix",
			input:    "/ztw/api/v1/ecgroup",
			expected: "/ecgroup",
		},
		{
			name:     "No prefix to remove",
			input:    "/api/v1/users",
			expected: "/api/v1/users",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "Just /zia",
			input:    "/zia",
			expected: "",
		},
		{
			name:     "Just /zpa",
			input:    "/zpa",
			expected: "",
		},
		{
			name:     "Case sensitive - ZIA uppercase should not match",
			input:    "/ZIA/api/v1/users",
			expected: "/ZIA/api/v1/users",
		},
		{
			name:     "Partial match - ziatest does match (prefix-based)",
			input:    "/ziatest/api",
			expected: "test/api", // HasPrefix matches /zia in /ziatest
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := removeOneApiEndpointPrefix(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestDifference(t *testing.T) {
	t.Run("Elements in slice1 not in slice2", func(t *testing.T) {
		slice1 := []string{"a", "b", "c", "d"}
		slice2 := []string{"b", "d"}

		result := Difference(slice1, slice2)

		assert.Len(t, result, 2)
		assert.Contains(t, result, "a")
		assert.Contains(t, result, "c")
	})

	t.Run("All elements in both slices", func(t *testing.T) {
		slice1 := []string{"a", "b", "c"}
		slice2 := []string{"a", "b", "c"}

		result := Difference(slice1, slice2)

		assert.Empty(t, result)
	})

	t.Run("No common elements", func(t *testing.T) {
		slice1 := []string{"a", "b", "c"}
		slice2 := []string{"x", "y", "z"}

		result := Difference(slice1, slice2)

		assert.Len(t, result, 3)
	})

	t.Run("Empty slices", func(t *testing.T) {
		result := Difference([]string{}, []string{})
		assert.Empty(t, result)
	})

	t.Run("Nil slices", func(t *testing.T) {
		var slice1 []string = nil
		var slice2 []string = nil

		result := Difference(slice1, slice2)
		assert.Empty(t, result)
	})
}

func TestContainsInt(t *testing.T) {
	t.Run("Contains existing element", func(t *testing.T) {
		codes := []int{200, 201, 204, 429, 500}
		assert.True(t, containsInt(codes, 429))
	})

	t.Run("Does not contain element", func(t *testing.T) {
		codes := []int{200, 201, 204}
		assert.False(t, containsInt(codes, 500))
	})

	t.Run("Empty slice", func(t *testing.T) {
		codes := []int{}
		assert.False(t, containsInt(codes, 200))
	})

	t.Run("Single element slice - found", func(t *testing.T) {
		codes := []int{429}
		assert.True(t, containsInt(codes, 429))
	})

	t.Run("Single element slice - not found", func(t *testing.T) {
		codes := []int{200}
		assert.False(t, containsInt(codes, 429))
	})
}

func TestGetRetryOnStatusCodes(t *testing.T) {
	t.Run("Returns expected status codes", func(t *testing.T) {
		codes := getRetryOnStatusCodes()

		assert.Contains(t, codes, 429) // TooManyRequests
	})

	t.Run("Does not include 4xx errors other than 429", func(t *testing.T) {
		codes := getRetryOnStatusCodes()

		assert.NotContains(t, codes, 400)
		assert.NotContains(t, codes, 401)
		assert.NotContains(t, codes, 403)
		assert.NotContains(t, codes, 404)
	})
}

func TestMin(t *testing.T) {
	t.Run("First value smaller", func(t *testing.T) {
		result := min(5, 10)
		assert.Equal(t, 5, result)
	})

	t.Run("Second value smaller", func(t *testing.T) {
		result := min(10, 5)
		assert.Equal(t, 5, result)
	})

	t.Run("Equal values", func(t *testing.T) {
		result := min(7, 7)
		assert.Equal(t, 7, result)
	})

	t.Run("Negative values", func(t *testing.T) {
		result := min(-5, -10)
		assert.Equal(t, -10, result)
	})

	t.Run("Zero and positive", func(t *testing.T) {
		result := min(0, 5)
		assert.Equal(t, 0, result)
	})
}

