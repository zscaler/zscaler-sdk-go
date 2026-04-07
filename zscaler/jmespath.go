package zscaler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/jmespath/go-jmespath"
)

type jmespathContextKey struct{}

// ContextWithJMESPath returns a derived context that carries a JMESPath
// expression. Pagination helpers (ReadAllPages, GetAllPagesGeneric*, etc.)
// automatically apply the expression to the aggregated result set before
// returning, enabling transparent client-side filtering across all services.
//
//	ctx := zscaler.ContextWithJMESPath(ctx, "[?enabled==`true`].{id: id, name: name}")
//	locations, err := location.GetAll(ctx, service, nil)
func ContextWithJMESPath(ctx context.Context, expression string) context.Context {
	return context.WithValue(ctx, jmespathContextKey{}, expression)
}

// JMESPathFromContext extracts the JMESPath expression stored in ctx, or
// returns "" if none was set.
func JMESPathFromContext(ctx context.Context) string {
	if v, ok := ctx.Value(jmespathContextKey{}).(string); ok {
		return v
	}
	return ""
}

// SearchJMESPath applies a JMESPath expression to any SDK result for
// client-side filtering and projection. It works with any []T or struct
// returned by list/get operations across all services (ZIA, ZPA, ZCC,
// ZDX, ZTW, ZID, ZWA).
//
// The data is marshaled to JSON first so that JMESPath operates on the
// camelCase field names defined by each struct's json tags.
//
// Returns the original data unchanged when expression is empty.
func SearchJMESPath(data interface{}, expression string) (interface{}, error) {
	if expression == "" {
		return data, nil
	}

	log.Printf("[DEBUG] jmespath: searching with expression %q", expression)

	b, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("jmespath: failed to marshal data: %w", err)
	}

	var generic interface{}
	if err := json.Unmarshal(b, &generic); err != nil {
		return nil, fmt.Errorf("jmespath: failed to unmarshal data: %w", err)
	}

	result, err := jmespath.Search(expression, generic)
	if err != nil {
		log.Printf("[ERROR] jmespath: expression %q failed: %v", expression, err)
		return nil, fmt.Errorf("jmespath: invalid expression: %w", err)
	}

	return result, nil
}

// ApplyJMESPathFilter applies a JMESPath expression to a typed slice and
// returns the filtered results as the same type. The expression must
// evaluate to an array of objects compatible with T (i.e., filter
// expressions like [?field=='value'], not projections that reshape).
//
// Returns the original slice unchanged when expression is empty.
func ApplyJMESPathFilter[T any](items []T, expression string) ([]T, error) {
	if expression == "" || len(items) == 0 {
		return items, nil
	}

	result, err := SearchJMESPath(items, expression)
	if err != nil {
		return nil, err
	}

	if result == nil {
		return []T{}, nil
	}

	b, err := json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("jmespath: failed to marshal filtered result: %w", err)
	}

	var filtered []T
	if err := json.Unmarshal(b, &filtered); err != nil {
		return nil, fmt.Errorf("jmespath: filtered result is not compatible with target type: %w", err)
	}

	return filtered, nil
}

// ApplyJMESPathFromContext is used by pagination engines to apply a JMESPath
// expression carried in the context. If no expression is present, the list
// is returned unchanged. All applications are logged for troubleshooting.
func ApplyJMESPathFromContext[T any](ctx context.Context, items []T) ([]T, error) {
	expr := JMESPathFromContext(ctx)
	if expr == "" {
		return items, nil
	}

	log.Printf("[DEBUG] jmespath: applying expression %q to %d items", expr, len(items))

	filtered, err := ApplyJMESPathFilter(items, expr)
	if err != nil {
		log.Printf("[ERROR] jmespath: expression %q failed: %v", expr, err)
		return nil, err
	}

	log.Printf("[DEBUG] jmespath: expression %q reduced %d items to %d", expr, len(items), len(filtered))
	return filtered, nil
}
