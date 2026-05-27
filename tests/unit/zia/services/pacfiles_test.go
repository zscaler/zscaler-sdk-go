package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/pacfiles"
)

const pacFilesPath = "/zia/api/v1/pacFiles"

const samplePACContent = `function FindProxyForURL(url, host) { return "DIRECT"; }`

func samplePACFile() pacfiles.PACFileConfig {
	return pacfiles.PACFileConfig{
		Name:                  "tests-sample.pac",
		Description:           "Test PAC file description",
		Domain:                "bd-hashicorp.com",
		PACContent:            samplePACContent,
		PACCommitMessage:      "Initial version via unit test",
		PACVersionStatus:      "DEPLOYED",
		PACVerificationStatus: "VERIFY_NOERR",
	}
}

func TestPacFiles_GetPacFiles_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", pacFilesPath, common.SuccessResponse([]pacfiles.PACFileConfig{
		samplePACFile(),
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := pacfiles.GetPacFiles(context.Background(), service, "")

	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "tests-sample.pac", result[0].Name)
}

func TestPacFiles_GetPacFileByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	pacName := "tests-sample.pac"
	server.On("GET", pacFilesPath, common.SuccessResponse([]pacfiles.PACFileConfig{
		func() pacfiles.PACFileConfig {
			p := samplePACFile()
			p.ID = 100
			return p
		}(),
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := pacfiles.GetPacFileByName(context.Background(), service, pacName)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, pacName, result.Name)
}

func TestPacFiles_CreatePacFile_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	created := samplePACFile()
	created.ID = 99999
	created.PACVersion = 1

	server.On("POST", pacFilesPath, common.SuccessResponse(created))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	input := samplePACFile()
	result, err := pacfiles.CreatePacFile(context.Background(), service, &input)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 99999, result.ID)
}

func TestPacFiles_GetPacFileVersion_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", pacFilesPath+"/100/version", common.SuccessResponse([]pacfiles.PACFileConfig{
		{ID: 100, Name: "tests-sample.pac", PACVersion: 1, PACVersionStatus: "DEPLOYED"},
		{ID: 100, Name: "tests-sample.pac", PACVersion: 2, PACVersionStatus: "STAGE"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := pacfiles.GetPacFileVersion(context.Background(), service, 100, "")

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestPacFiles_GetPacVersionID_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", pacFilesPath+"/100/version/1", common.SuccessResponse(pacfiles.PACFileConfig{
		ID:               100,
		Name:             "tests-sample.pac",
		PACVersion:       1,
		PACVersionStatus: "DEPLOYED",
		PACContent:       samplePACContent,
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := pacfiles.GetPacVersionID(context.Background(), service, 100, 1, "")

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 1, result.PACVersion)
}

func TestPacFiles_UpdatePacFile_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	updated := samplePACFile()
	updated.ID = 100
	updated.PACVersion = 1
	updated.PACVersionStatus = "STAGE"

	server.On("PUT", pacFilesPath+"/100/version/1/action/STAGE", common.SuccessResponse(updated))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	input := samplePACFile()
	result, err := pacfiles.UpdatePacFile(context.Background(), service, 100, 1, "STAGE", &input, nil)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "STAGE", result.PACVersionStatus)
}

func TestPacFiles_DeletePacFile_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("DELETE", pacFilesPath+"/100", common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, err = pacfiles.DeletePacFile(context.Background(), service, 100)

	require.NoError(t, err)
}

func TestPacFiles_ValidatePacFile_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("POST", pacFilesPath+"/validate", common.SuccessResponse(pacfiles.PacResult{
		Success:      true,
		WarningCount: 0,
		ErrorCount:   0,
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := pacfiles.ValidatePacFile(context.Background(), service, samplePACContent)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.True(t, result.Success)
}

func TestPacFiles_CreateClonedPacFileVersion_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	cloned := samplePACFile()
	cloned.ID = 100
	cloned.PACVersion = 2
	cloned.Name = "tests-sample-cloned.pac"

	server.On("POST", pacFilesPath+"/100/version/1", common.SuccessResponse(cloned))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	input := samplePACFile()
	input.Name = "tests-sample-cloned.pac"
	result, err := pacfiles.CreateClonedPacFileVersion(context.Background(), service, 100, 1, nil, &input)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 2, result.PACVersion)
}

func TestPacFiles_Structure(t *testing.T) {
	t.Parallel()

	t.Run("PACFileConfig JSON marshaling", func(t *testing.T) {
		pac := samplePACFile()
		pac.ID = 100

		data, err := json.Marshal(pac)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"pacVersionStatus":"DEPLOYED"`)
		assert.Contains(t, string(data), `"pacVerificationStatus":"VERIFY_NOERR"`)
	})
}
