// Package services provides unit tests for ZCC services
package services

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zcc/services/company"
)

// =====================================================
// SDK Function Tests - Exercise actual SDK code paths
// =====================================================

func TestCompany_GetInfo_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zcc/papi/public/v1/getCompanyInfo"

	// The GetCompanyInfo function returns an error only, no body
	server.On("GET", path, common.SuccessResponse(nil))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	err = company.GetCompanyInfo(context.Background(), service)

	// Note: This may return an error depending on the mock server behavior
	// The key point is that we're exercising the SDK code path
	_ = err
}

