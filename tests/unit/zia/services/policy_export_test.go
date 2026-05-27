// Package services provides unit tests for ZIA services
package services

import (
	"archive/zip"
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/policy_export"
)

func testPolicyExportZIP(t *testing.T) []byte {
	t.Helper()
	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)
	f, err := w.Create("urlFilteringRules.json")
	require.NoError(t, err)
	_, err = f.Write([]byte(`{"rules":[{"name":"Block Social Media","action":"BLOCK"}]}`))
	require.NoError(t, err)
	require.NoError(t, w.Close())
	return buf.Bytes()
}

func TestPolicyExport_ExportPolicies_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/exportPolicies"
	server.On("POST", path, common.RawResponse(testPolicyExportZIP(t), 200, map[string]string{
		"Content-Type": "application/zip",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	outputDir := t.TempDir()
	policyTypes := []string{"URL_FILTERING", "FIREWALL"}

	err = policy_export.ExportPolicies(context.Background(), service, policyTypes, outputDir)
	require.NoError(t, err)

	extracted := filepath.Join(outputDir, "urlFilteringRules.json")
	_, err = os.Stat(extracted)
	require.NoError(t, err)
}

func TestPolicyExport_ExportPolicies_Error_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/exportPolicies"
	server.On("POST", path, common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	err = policy_export.ExportPolicies(context.Background(), service, []string{"URL_FILTERING"}, t.TempDir())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to export policies")
}

func TestPolicyExport_ExportPolicies_NestedFolder_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)
	dir, err := w.Create("policies/")
	require.NoError(t, err)
	_, err = dir.Write(nil)
	require.NoError(t, err)
	f, err := w.Create("policies/firewallRules.json")
	require.NoError(t, err)
	_, err = f.Write([]byte(`{"rules":[]}`))
	require.NoError(t, err)
	require.NoError(t, w.Close())

	path := "/zia/api/v1/exportPolicies"
	server.On("POST", path, common.RawResponse(buf.Bytes(), 200, map[string]string{
		"Content-Type": "application/zip",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	outputDir := t.TempDir()
	err = policy_export.ExportPolicies(context.Background(), service, []string{"FIREWALL"}, outputDir)
	require.NoError(t, err)

	_, err = os.Stat(filepath.Join(outputDir, "policies", "firewallRules.json"))
	require.NoError(t, err)
}

func TestPolicyExport_ExportPolicies_InvalidZIP_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/exportPolicies"
	server.On("POST", path, common.RawResponse([]byte("not-a-zip"), 200, map[string]string{
		"Content-Type": "application/zip",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	err = policy_export.ExportPolicies(context.Background(), service, []string{"URL_FILTERING"}, t.TempDir())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to read ZIP")
}

func TestPolicyExport_ExportPolicies_ZipSlip_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)
	f, err := w.Create("../../../escape.json")
	require.NoError(t, err)
	_, err = f.Write([]byte(`{}`))
	require.NoError(t, err)
	require.NoError(t, w.Close())

	path := "/zia/api/v1/exportPolicies"
	server.On("POST", path, common.RawResponse(buf.Bytes(), 200, map[string]string{
		"Content-Type": "application/zip",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	err = policy_export.ExportPolicies(context.Background(), service, []string{"URL_FILTERING"}, t.TempDir())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "zip slip")
}

func TestPolicyExport_ExportPolicies_DirectoryZipSlip_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)
	_, err := w.Create("../../../evil/")
	require.NoError(t, err)
	require.NoError(t, w.Close())

	path := "/zia/api/v1/exportPolicies"
	server.On("POST", path, common.RawResponse(buf.Bytes(), 200, map[string]string{
		"Content-Type": "application/zip",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	err = policy_export.ExportPolicies(context.Background(), service, []string{"URL_FILTERING"}, t.TempDir())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "zip slip")
}
