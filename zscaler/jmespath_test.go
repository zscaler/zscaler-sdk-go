package zscaler

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testDevice struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	OSType string `json:"osType"`
	Status string `json:"status"`
}

func TestSearchJMESPath_EmptyExpression(t *testing.T) {
	devices := []testDevice{
		{ID: "1", Name: "dev1", OSType: "Windows", Status: "active"},
	}
	result, err := SearchJMESPath(devices, "")
	require.NoError(t, err)
	assert.Equal(t, devices, result)
}

func TestSearchJMESPath_FilterByField(t *testing.T) {
	devices := []testDevice{
		{ID: "1", Name: "laptop-01", OSType: "Windows", Status: "active"},
		{ID: "2", Name: "mac-01", OSType: "macOS", Status: "active"},
		{ID: "3", Name: "laptop-02", OSType: "Windows", Status: "inactive"},
	}

	result, err := SearchJMESPath(devices, "[?osType=='Windows']")
	require.NoError(t, err)

	items, ok := result.([]interface{})
	require.True(t, ok)
	assert.Len(t, items, 2)
}

func TestSearchJMESPath_Projection(t *testing.T) {
	devices := []testDevice{
		{ID: "1", Name: "laptop-01", OSType: "Windows", Status: "active"},
		{ID: "2", Name: "mac-01", OSType: "macOS", Status: "active"},
	}

	result, err := SearchJMESPath(devices, "[*].{name: name, os: osType}")
	require.NoError(t, err)

	items, ok := result.([]interface{})
	require.True(t, ok)
	assert.Len(t, items, 2)

	first := items[0].(map[string]interface{})
	assert.Equal(t, "laptop-01", first["name"])
	assert.Equal(t, "Windows", first["os"])
}

func TestSearchJMESPath_FilterAndProject(t *testing.T) {
	devices := []testDevice{
		{ID: "1", Name: "laptop-01", OSType: "Windows", Status: "active"},
		{ID: "2", Name: "mac-01", OSType: "macOS", Status: "active"},
		{ID: "3", Name: "laptop-02", OSType: "Windows", Status: "inactive"},
	}

	result, err := SearchJMESPath(devices, "[?osType=='Windows'].name")
	require.NoError(t, err)

	names, ok := result.([]interface{})
	require.True(t, ok)
	assert.Equal(t, []interface{}{"laptop-01", "laptop-02"}, names)
}

func TestSearchJMESPath_SingleStruct(t *testing.T) {
	device := testDevice{ID: "1", Name: "laptop-01", OSType: "Windows", Status: "active"}

	result, err := SearchJMESPath(device, "name")
	require.NoError(t, err)
	assert.Equal(t, "laptop-01", result)
}

func TestSearchJMESPath_NestedWrapper(t *testing.T) {
	type wrapper struct {
		TotalCount int          `json:"totalCount"`
		Items      []testDevice `json:"items"`
	}

	data := wrapper{
		TotalCount: 3,
		Items: []testDevice{
			{ID: "1", Name: "a", OSType: "Windows", Status: "active"},
			{ID: "2", Name: "b", OSType: "macOS", Status: "active"},
			{ID: "3", Name: "c", OSType: "Linux", Status: "inactive"},
		},
	}

	result, err := SearchJMESPath(data, "items[?status=='active'].name")
	require.NoError(t, err)

	names, ok := result.([]interface{})
	require.True(t, ok)
	assert.Equal(t, []interface{}{"a", "b"}, names)
}

func TestSearchJMESPath_MapSlice(t *testing.T) {
	data := []map[string]interface{}{
		{"name": "policy-1", "action": "allow"},
		{"name": "policy-2", "action": "block"},
		{"name": "policy-3", "action": "allow"},
	}

	result, err := SearchJMESPath(data, "[?action=='block'].name")
	require.NoError(t, err)

	names, ok := result.([]interface{})
	require.True(t, ok)
	assert.Equal(t, []interface{}{"policy-2"}, names)
}

func TestSearchJMESPath_InvalidExpression(t *testing.T) {
	devices := []testDevice{{ID: "1", Name: "dev1"}}

	_, err := SearchJMESPath(devices, "[?invalid===")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "jmespath: invalid expression")
}

func TestSearchJMESPath_NilData(t *testing.T) {
	result, err := SearchJMESPath(nil, "foo")
	require.NoError(t, err)
	assert.Nil(t, result)
}

func TestSearchJMESPath_EmptySlice(t *testing.T) {
	result, err := SearchJMESPath([]testDevice{}, "[?osType=='Windows']")
	require.NoError(t, err)
	items, ok := result.([]interface{})
	require.True(t, ok)
	assert.Empty(t, items)
}

func TestSearchJMESPath_LengthFunction(t *testing.T) {
	devices := []testDevice{
		{ID: "1", Name: "a"},
		{ID: "2", Name: "b"},
		{ID: "3", Name: "c"},
	}

	result, err := SearchJMESPath(devices, "length(@)")
	require.NoError(t, err)
	assert.Equal(t, float64(3), result)
}

// --- Context-based JMESPath tests ---

func TestContextWithJMESPath_RoundTrip(t *testing.T) {
	ctx := context.Background()
	assert.Equal(t, "", JMESPathFromContext(ctx))

	ctx = ContextWithJMESPath(ctx, "[?enabled==`true`]")
	assert.Equal(t, "[?enabled==`true`]", JMESPathFromContext(ctx))
}

func TestContextWithJMESPath_EmptyExpression(t *testing.T) {
	ctx := ContextWithJMESPath(context.Background(), "")
	assert.Equal(t, "", JMESPathFromContext(ctx))
}

// --- ApplyJMESPathFilter tests ---

func TestApplyJMESPathFilter_FilterByField(t *testing.T) {
	devices := []testDevice{
		{ID: "1", Name: "laptop-01", OSType: "Windows", Status: "active"},
		{ID: "2", Name: "mac-01", OSType: "macOS", Status: "active"},
		{ID: "3", Name: "laptop-02", OSType: "Windows", Status: "inactive"},
	}

	filtered, err := ApplyJMESPathFilter(devices, "[?osType=='Windows']")
	require.NoError(t, err)
	assert.Len(t, filtered, 2)
	assert.Equal(t, "laptop-01", filtered[0].Name)
	assert.Equal(t, "laptop-02", filtered[1].Name)
}

func TestApplyJMESPathFilter_EmptyExpression(t *testing.T) {
	devices := []testDevice{{ID: "1", Name: "dev1"}}

	filtered, err := ApplyJMESPathFilter(devices, "")
	require.NoError(t, err)
	assert.Equal(t, devices, filtered)
}

func TestApplyJMESPathFilter_EmptySlice(t *testing.T) {
	filtered, err := ApplyJMESPathFilter([]testDevice{}, "[?osType=='Windows']")
	require.NoError(t, err)
	assert.Empty(t, filtered)
}

func TestApplyJMESPathFilter_NoMatches(t *testing.T) {
	devices := []testDevice{
		{ID: "1", Name: "dev1", OSType: "Windows"},
	}

	filtered, err := ApplyJMESPathFilter(devices, "[?osType=='Linux']")
	require.NoError(t, err)
	assert.Empty(t, filtered)
}

func TestApplyJMESPathFilter_InvalidExpression(t *testing.T) {
	devices := []testDevice{{ID: "1", Name: "dev1"}}

	_, err := ApplyJMESPathFilter(devices, "[?invalid===")
	assert.Error(t, err)
}

// --- ApplyJMESPathFromContext tests ---

func TestApplyJMESPathFromContext_WithExpression(t *testing.T) {
	devices := []testDevice{
		{ID: "1", Name: "laptop-01", OSType: "Windows", Status: "active"},
		{ID: "2", Name: "mac-01", OSType: "macOS", Status: "active"},
		{ID: "3", Name: "laptop-02", OSType: "Windows", Status: "inactive"},
	}

	ctx := ContextWithJMESPath(context.Background(), "[?status=='active']")
	filtered, err := ApplyJMESPathFromContext(ctx, devices)
	require.NoError(t, err)
	assert.Len(t, filtered, 2)
	assert.Equal(t, "laptop-01", filtered[0].Name)
	assert.Equal(t, "mac-01", filtered[1].Name)
}

func TestApplyJMESPathFromContext_NoExpression(t *testing.T) {
	devices := []testDevice{
		{ID: "1", Name: "dev1"},
		{ID: "2", Name: "dev2"},
	}

	filtered, err := ApplyJMESPathFromContext(context.Background(), devices)
	require.NoError(t, err)
	assert.Equal(t, devices, filtered)
}

func TestApplyJMESPathFromContext_MultipleFilters(t *testing.T) {
	devices := []testDevice{
		{ID: "1", Name: "laptop-01", OSType: "Windows", Status: "active"},
		{ID: "2", Name: "mac-01", OSType: "macOS", Status: "active"},
		{ID: "3", Name: "laptop-02", OSType: "Windows", Status: "inactive"},
	}

	ctx := ContextWithJMESPath(context.Background(), "[?osType=='Windows' && status=='active']")
	filtered, err := ApplyJMESPathFromContext(ctx, devices)
	require.NoError(t, err)
	assert.Len(t, filtered, 1)
	assert.Equal(t, "laptop-01", filtered[0].Name)
}
