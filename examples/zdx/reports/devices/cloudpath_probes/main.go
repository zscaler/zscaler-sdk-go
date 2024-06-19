package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/zscaler/zscaler-sdk-go/v2/zdx"
	"github.com/zscaler/zscaler-sdk-go/v2/zdx/services"
	"github.com/zscaler/zscaler-sdk-go/v2/zdx/services/common"
	"github.com/zscaler/zscaler-sdk-go/v2/zdx/services/reports/devices"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	// Check for API key and secret in environment variables
	apiKey := os.Getenv("ZDX_API_KEY_ID")
	apiSecret := os.Getenv("ZDX_API_SECRET")

	if apiKey == "" || apiSecret == "" {
		log.Fatalf("[ERROR] API key and secret must be set in environment variables (ZDX_API_KEY_ID, ZDX_API_SECRET)\n")
	}

	// Create configuration and client
	cfg, err := zdx.NewConfig(apiKey, apiSecret, "userAgent")
	if err != nil {
		log.Fatalf("[ERROR] creating client failed: %v\n", err)
	}
	cli := zdx.NewClient(cfg)
	service := services.New(cli)

	// Prompt the user to choose a CloudPath Probe report
	fmt.Println("Choose the CloudPath Probe Report:")
	fmt.Println("a. All active Cloud Path probes on a device")
	fmt.Println("b. Web probe's Page Fetch Time (PFT) for an application")
	fmt.Println("c. Cloud Path hop data for an application")
	fmt.Print("Enter choice (a/b/c): ")
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	// Prompt for common inputs
	fmt.Print("Enter device ID: ")
	deviceIDInput, _ := reader.ReadString('\n')
	deviceIDInput = strings.TrimSpace(deviceIDInput)
	deviceID, err := strconv.Atoi(deviceIDInput)
	if err != nil {
		log.Fatalf("[ERROR] Invalid device ID: %v\n", err)
	}

	fmt.Print("Enter app ID: ")
	appIDInput, _ := reader.ReadString('\n')
	appIDInput = strings.TrimSpace(appIDInput)
	appID, err := strconv.Atoi(appIDInput)
	if err != nil {
		log.Fatalf("[ERROR] Invalid app ID: %v\n", err)
	}

	// Define a time filter for the last 2 hours
	now := time.Now()
	from := now.Add(-2 * time.Hour).Unix()
	to := now.Unix()
	filters := common.GetFromToFilters{
		From: int(from),
		To:   int(to),
	}

	switch choice {
	case "a":
		// Get all Cloud Path probes
		probes, _, err := devices.GetAllCloudPathProbes(service, deviceID, appID, filters)
		if err != nil {
			log.Fatalf("Error getting cloud path probes: %v", err)
		}
		displayCloudPathProbes(probes)
	case "b":
		// Prompt for probe ID
		fmt.Print("Enter probe ID: ")
		probeIDInput, _ := reader.ReadString('\n')
		probeIDInput = strings.TrimSpace(probeIDInput)
		probeID, err := strconv.Atoi(probeIDInput)
		if err != nil {
			log.Fatalf("[ERROR] Invalid probe ID: %v\n", err)
		}

		// Get device app Cloud Path probe
		networkStats, _, err := devices.GetDeviceAppCloudPathProbe(service, deviceID, appID, probeID, filters)
		if err != nil {
			log.Fatalf("Error getting device app cloud path probe: %v", err)
		}
		displayNetworkStats(networkStats)
	case "c":
		// Prompt for probe ID
		fmt.Print("Enter probe ID: ")
		probeIDInput, _ := reader.ReadString('\n')
		probeIDInput = strings.TrimSpace(probeIDInput)
		probeID, err := strconv.Atoi(probeIDInput)
		if err != nil {
			log.Fatalf("[ERROR] Invalid probe ID: %v\n", err)
		}

		// Get Cloud Path app device
		cloudPathProbes, _, err := devices.GetCloudPathAppDevice(service, deviceID, appID, probeID, filters)
		if err != nil {
			log.Fatalf("Error getting cloud path app device: %v", err)
		}
		displayCloudPathAppDevice(cloudPathProbes)
	default:
		log.Fatalf("[ERROR] Invalid choice: %s\n", choice)
	}
}

func displayCloudPathProbes(probes []devices.DeviceCloudPathProbe) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Probe ID", "Name", "NumProbes", "LegSRC", "LegDst", "Latency"})

	for _, probe := range probes {
		for _, latency := range probe.AverageLatency {
			table.Append([]string{
				strconv.Itoa(probe.ID),
				probe.Name,
				strconv.Itoa(probe.NumProbes),
				latency.LegSRC,
				latency.LegDst,
				fmt.Sprintf("%f", latency.Latency),
			})
		}
	}
	table.Render()
}

func displayNetworkStats(stats []devices.NetworkStats) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"LegSRC", "LegDst", "Metric", "Unit", "Timestamp", "Value"})

	for _, stat := range stats {
		for _, metric := range stat.Stats {
			for _, dp := range metric.DataPoints {
				table.Append([]string{
					stat.LegSRC,
					stat.LegDst,
					metric.Metric,
					metric.Unit,
					time.Unix(int64(dp.TimeStamp), 0).Format("2006-01-02 15:04:05"),
					fmt.Sprintf("%f", dp.Value),
				})
			}
		}
	}
	table.Render()
}

func displayCloudPathAppDevice(probes []devices.CloudPathProbe) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Timestamp", "SRC", "DST", "NumHops", "Latency", "Loss", "NumUnrespHops", "TunnelType", "HopIP", "HopGWMac", "HopGWMacVendor", "PktSent", "PktRcvd", "LatencyMin", "LatencyMax", "LatencyAvg", "LatencyDiff"})

	for _, probe := range probes {
		for _, path := range probe.CloudPath {
			for _, hop := range path.Hops {
				table.Append([]string{
					time.Unix(int64(probe.TimeStamp), 0).Format("2006-01-02 15:04:05"),
					path.SRC,
					path.DST,
					strconv.Itoa(path.NumHops),
					fmt.Sprintf("%f", path.Latency),
					fmt.Sprintf("%f", path.Loss),
					strconv.Itoa(path.NumUnrespHops),
					strconv.Itoa(path.TunnelType),
					hop.IP,
					hop.GWMac,
					hop.GWMacVendor,
					strconv.Itoa(hop.PktSent),
					strconv.Itoa(hop.PktRcvd),
					strconv.Itoa(hop.LatencyMin),
					strconv.Itoa(hop.LatencyMax),
					strconv.Itoa(hop.LatencyAvg),
					strconv.Itoa(hop.LatencyDiff),
				})
			}
		}
	}
	table.Render()
}
