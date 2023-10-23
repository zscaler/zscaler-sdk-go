package appconnectorcontroller

/*
import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/zscaler/zscaler-sdk-go/v2/tests"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/appconnectorgroup"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/enrollmentcert"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/provisioningkey"
)

const (
	connGrpAssociationType = "CONNECTOR_GRP"
)

func startDockerContainer(imageName, provisionKey, containerName string) (string, error) {
	ctx := context.Background()

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return "", fmt.Errorf("Error initializing Docker client: %v", err)
	}

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: imageName,
		Env:   []string{"ZPA_PROVISION_KEY=" + provisionKey},
	}, &container.HostConfig{
		CapAdd:        []string{"cap_net_admin", "cap_net_bind_service", "cap_net_raw", "cap_sys_nice", "cap_sys_time"},
		RestartPolicy: container.RestartPolicy{Name: "always"},
	}, nil, nil, containerName)
	if err != nil {
		return "", fmt.Errorf("Error creating container: %v", err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return "", fmt.Errorf("Error starting container: %v", err)
	}

	return resp.ID, nil
}

func pullDockerImage(imageName string) error {
	ctx := context.Background()

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return fmt.Errorf("Error initializing Docker client: %v", err)
	}

	reader, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		return fmt.Errorf("Error pulling Docker image: %v", err)
	}
	defer reader.Close()

	// This will block until the image is downloaded
	io.Copy(ioutil.Discard, reader)

	return nil
}

func generateContainerName(prefix string) string {
	randStr := acctest.RandStringFromCharSet(6, acctest.CharSetAlpha)
	return prefix + "-" + randStr
}

func TestAppConnectorController(t *testing.T) {
	// Generate names for resources
	name := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	appConnGroupName := "tests-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	imageName := "zscaler/zpa-connector:latest.arm64"

	// Ensure the Docker image is available
	if err := pullDockerImage(imageName); err != nil {
		t.Fatalf("Error pulling Docker image: %v", err)
		return
	}

	client, err := tests.NewZpaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
		return
	}

	// Create application connector group for testing
	appConnectorGroupService := appconnectorgroup.New(client)
	appGroup := appconnectorgroup.AppConnectorGroup{
		Name:                     appConnGroupName,
		Description:              appConnGroupName,
		Enabled:                  true,
		CityCountry:              "San Jose, US",
		Latitude:                 "37.3382082",
		Longitude:                "-121.8863286",
		Location:                 "San Jose, CA, USA",
		UpgradeDay:               "SUNDAY",
		UpgradeTimeInSecs:        "66600",
		OverrideVersionProfile:   true,
		VersionProfileName:       "Default",
		VersionProfileID:         "0",
		DNSQueryType:             "IPV4_IPV6",
		PRAEnabled:               false,
		WAFDisabled:              true,
		TCPQuickAckApp:           true,
		TCPQuickAckAssistant:     true,
		TCPQuickAckReadAssistant: true,
	}

	createdAppConnGroup, _, err := appConnectorGroupService.Create(appGroup)
	if err != nil || createdAppConnGroup == nil || createdAppConnGroup.ID == "" {
		t.Fatalf("Error creating application connector group or ID is empty")
		return
	}
	defer cleanupAppConnGroup(t, appConnectorGroupService, createdAppConnGroup.ID)

	// Get enrollment cert for testing
	enrollmentCertService := enrollmentcert.New(client)
	enrollmentCert, _, err := enrollmentCertService.GetByName("Connector")
	if err != nil {
		t.Fatalf("Error getting enrollment cert: %v", err)
		return
	}

	// Create provisioning key
	provisioningKeyService := provisioningkey.New(client)
	provisioningKey := provisioningkey.ProvisioningKey{
		AssociationType:       connGrpAssociationType,
		Name:                  name,
		AppConnectorGroupID:   createdAppConnGroup.ID,
		AppConnectorGroupName: createdAppConnGroup.Name,
		EnrollmentCertID:      enrollmentCert.ID,
		ZcomponentID:          createdAppConnGroup.ID,
		MaxUsage:              "10",
	}

	createdResource, _, err := provisioningKeyService.Create(connGrpAssociationType, &provisioningKey)
	if err != nil || createdResource == nil || createdResource.ID == "" {
		t.Fatalf("Error making POST request or created resource is nil/empty: %v", err)
		return
	}

	// Start two Docker containers using the provisioningKey from the response
	provisionKey := createdResource.ProvisioningKey
	// Create the first container of the pair
	container1Name := generateContainerName("zpa-connector")
	container1ID, err := startDockerContainer(imageName, provisionKey, container1Name)
	if err != nil {
		t.Fatalf("Error starting first container: %v", err)
		return
	}
	defer stopAndRemoveContainer(container1ID) // Cleanup after test

	// Create the second container of the pair
	container2Name := generateContainerName("zpa-connector")
	container2ID, err := startDockerContainer(imageName, provisionKey, container2Name)
	if err != nil {
		t.Fatalf("Error starting second container: %v", err)
		return
	}
	defer stopAndRemoveContainer(container2ID) // Cleanup after test

	// Continue with rest of your tests...
}

// Helper function to cleanup App Connector Group
func cleanupAppConnGroup(t *testing.T, service *appconnectorgroup.Service, id string) {
	existingGroup, _, errCheck := service.Get(id)
	if errCheck == nil && existingGroup != nil {
		_, errDelete := service.Delete(id)
		if errDelete != nil {
			t.Errorf("Error deleting application connector group: %v", errDelete)
		}
	}
}

func stopAndRemoveContainer(containerID string) {
	ctx := context.Background()

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		fmt.Printf("Error initializing Docker client: %v", err)
		return
	}

	stopOptions := container.StopOptions{}
	if err := cli.ContainerStop(ctx, containerID, stopOptions); err != nil {
		fmt.Printf("Error stopping container %s: %v", containerID, err)
		return
	}

	if err := cli.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{}); err != nil {
		fmt.Printf("Error removing container %s: %v", containerID, err)
	}
}
*/
