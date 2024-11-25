package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler/zdx"
	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler/zdx/services"
	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler/zdx/services/inventory"
	"github.com/olekukonko/tablewriter"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	// Check for API key and secret in environment variables
	apiKey := os.Getenv("ZDX_API_KEY_ID")
	apiSecret := os.Getenv("ZDX_API_SECRET")

	if apiKey == "" || apiSecret == "" {
		log.Fatalf("[ERROR] API key and secret must be set in environment variables (ZDX_API_KEY_ID, ZDX_API_SECRET)\n")
	}

	// Prompt for Zscaler location IDs (comma-separated list)
	fmt.Print("Enter Zscaler location IDs (comma-separated, optional): ")
	locInput, _ := reader.ReadString('\n')
	locInput = strings.TrimSpace(locInput)
	locIDs := parseCommaSeparatedInts(locInput)

	// Prompt for department IDs (comma-separated list)
	fmt.Print("Enter department IDs (comma-separated, optional): ")
	deptInput, _ := reader.ReadString('\n')
	deptInput = strings.TrimSpace(deptInput)
	deptIDs := parseCommaSeparatedInts(deptInput)

	// Prompt for geolocation IDs (comma-separated list)
	fmt.Print("Enter geolocation IDs (comma-separated, optional): ")
	geoInput, _ := reader.ReadString('\n')
	geoInput = strings.TrimSpace(geoInput)
	geoIDs := parseCommaSeparatedStrings(geoInput)

	// Prompt for user IDs (comma-separated list)
	fmt.Print("Enter user IDs (comma-separated, optional): ")
	userInput, _ := reader.ReadString('\n')
	userInput = strings.TrimSpace(userInput)
	userIDs := parseCommaSeparatedInts(userInput)

	// Prompt for device IDs (comma-separated list)
	fmt.Print("Enter device IDs (comma-separated, optional): ")
	deviceInput, _ := reader.ReadString('\n')
	deviceInput = strings.TrimSpace(deviceInput)
	deviceIDs := parseCommaSeparatedInts(deviceInput)

	// Prompt for software key (optional)
	fmt.Print("Enter software key (optional): ")
	softwareKeyInput, _ := reader.ReadString('\n')
	softwareKeyInput = strings.TrimSpace(softwareKeyInput)

	// Create configuration and client
	cfg, err := zdx.NewConfig(apiKey, apiSecret, "userAgent")
	if err != nil {
		log.Fatalf("[ERROR] creating client failed: %v\n", err)
	}
	cli := zdx.NewClient(cfg)
	softwareService := services.New(cli)

	// Define filters
	filters := inventory.GetSoftwareFilters{
		Loc:         locIDs,
		Dept:        deptIDs,
		Geo:         geoIDs,
		UserIDs:     userIDs,
		DeviceIDs:   deviceIDs,
		SoftwareKey: softwareKeyInput,
	}

	if softwareKeyInput != "" {
		// Get software key details
		softwareList, nextOffset, resp, err := inventory.GetSoftwareKey(softwareService, softwareKeyInput, filters)
		if err != nil {
			log.Fatalf("[ERROR] getting software key details failed: %v\n", err)
		}

		if resp.StatusCode != http.StatusOK {
			log.Fatalf("[ERROR] Expected status code 200, got %d", resp.StatusCode)
		}

		// Display the data in a formatted table
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Software Key", "Software Name", "Vendor", "Software Version", "OS", "Hostname", "Username", "Install Date"})
		table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
		table.SetCenterSeparator("|")

		for _, software := range softwareList {
			installDate := formatEpochToDate(software.InstallDate)
			table.Append([]string{
				software.SoftwareKey,
				software.SoftwareName,
				software.Vendor,
				software.SoftwareVersion,
				software.OS,
				software.Hostname,
				software.Username,
				installDate,
			})
		}

		table.Render()
		log.Printf("Next Offset: %s", nextOffset)

	} else {
		// Get software inventory
		softwareList, nextOffset, resp, err := inventory.GetSoftware(softwareService, filters)
		if err != nil {
			log.Fatalf("[ERROR] getting software inventory failed: %v\n", err)
		}

		if resp.StatusCode != http.StatusOK {
			log.Fatalf("[ERROR] Expected status code 200, got %d", resp.StatusCode)
		}

		// Display the data in a formatted table
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Software Key", "Software Name", "Vendor", "Software Group", "Install Type", "User Total", "Device Total"})
		table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
		table.SetCenterSeparator("|")

		for _, software := range softwareList {
			table.Append([]string{
				software.SoftwareKey,
				software.SoftwareName,
				software.Vendor,
				software.SoftwareGroup,
				software.SoftwareInstallType,
				strconv.Itoa(software.UserTotal),
				strconv.Itoa(software.DeviceTotal),
			})
		}

		table.Render()
		log.Printf("Next Offset: %s", nextOffset)
	}
}

// Helper functions to parse comma-separated input
func parseCommaSeparatedInts(input string) []int {
	if input == "" {
		return nil
	}
	strs := strings.Split(input, ",")
	ints := make([]int, len(strs))
	for i, str := range strs {
		val, err := strconv.Atoi(strings.TrimSpace(str))
		if err != nil {
			log.Fatalf("[ERROR] Invalid integer value: %v\n", err)
		}
		ints[i] = val
	}
	return ints
}

func parseCommaSeparatedStrings(input string) []string {
	if input == "" {
		return nil
	}
	strs := strings.Split(input, ",")
	for i, str := range strs {
		strs[i] = strings.TrimSpace(str)
	}
	return strs
}

// Helper function to convert epoch time to "YYYY-MM-DD" format
func formatEpochToDate(epochStr string) string {
	epoch, err := strconv.ParseInt(epochStr, 10, 64)
	if err != nil {
		return epochStr // Return the original string if parsing fails
	}
	t := time.Unix(epoch, 0)
	return t.Format("2006-01-02")
}
