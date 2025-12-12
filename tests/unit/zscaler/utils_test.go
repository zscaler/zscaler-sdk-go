// Package zscaler provides unit tests for core zscaler SDK utilities
package zscaler

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
)

func TestUtils_Difference(t *testing.T) {
	t.Parallel()

	t.Run("Elements in slice1 not in slice2", func(t *testing.T) {
		slice1 := []string{"a", "b", "c", "d"}
		slice2 := []string{"b", "d"}

		result := zscaler.Difference(slice1, slice2)

		assert.Len(t, result, 2)
		assert.Contains(t, result, "a")
		assert.Contains(t, result, "c")
	})

	t.Run("All elements in both slices", func(t *testing.T) {
		slice1 := []string{"a", "b", "c"}
		slice2 := []string{"a", "b", "c"}

		result := zscaler.Difference(slice1, slice2)

		assert.Empty(t, result)
	})

	t.Run("No common elements", func(t *testing.T) {
		slice1 := []string{"a", "b", "c"}
		slice2 := []string{"x", "y", "z"}

		result := zscaler.Difference(slice1, slice2)

		assert.Len(t, result, 3)
		assert.Contains(t, result, "a")
		assert.Contains(t, result, "b")
		assert.Contains(t, result, "c")
	})

	t.Run("Empty slice1", func(t *testing.T) {
		slice1 := []string{}
		slice2 := []string{"a", "b"}

		result := zscaler.Difference(slice1, slice2)

		assert.Empty(t, result)
	})

	t.Run("Empty slice2", func(t *testing.T) {
		slice1 := []string{"a", "b", "c"}
		slice2 := []string{}

		result := zscaler.Difference(slice1, slice2)

		assert.Len(t, result, 3)
	})

	t.Run("Both slices empty", func(t *testing.T) {
		slice1 := []string{}
		slice2 := []string{}

		result := zscaler.Difference(slice1, slice2)

		assert.Empty(t, result)
	})

	t.Run("Duplicate elements in slice1", func(t *testing.T) {
		slice1 := []string{"a", "a", "b", "b", "c"}
		slice2 := []string{"b"}

		result := zscaler.Difference(slice1, slice2)

		// Should have duplicates of "a" and one "c"
		assert.Len(t, result, 3)
	})

	t.Run("Single element slices", func(t *testing.T) {
		slice1 := []string{"x"}
		slice2 := []string{"y"}

		result := zscaler.Difference(slice1, slice2)

		assert.Len(t, result, 1)
		assert.Equal(t, "x", result[0])
	})

	t.Run("Nil slices", func(t *testing.T) {
		var slice1 []string = nil
		var slice2 []string = nil

		result := zscaler.Difference(slice1, slice2)

		assert.Empty(t, result)
	})
}

