package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zdx/services/troubleshooting/deeptrace"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	clientID := os.Getenv("ZSCALER_CLIENT_ID")
	clientSecret := os.Getenv("ZSCALER_CLIENT_SECRET")
	vanityDomain := os.Getenv("ZSCALER_VANITY_DOMAIN")

	// Initialize ZDX configuration
	zdxCfg, err := zscaler.NewConfiguration(
		zscaler.WithClientID(clientID),
		zscaler.WithClientSecret(clientSecret),
		zscaler.WithVanityDomain(vanityDomain),
		zscaler.WithDebug(true),
	)
	if err != nil {
		log.Fatalf("Error creating ZDX configuration: %v", err)
	}

	service, err := zscaler.NewOneAPIClient(zdxCfg)
	if err != nil {
		log.Fatalf("Error creating OneAPI client: %v", err)
	}

	ctx := context.Background()

	// Prompt the user for required IDs
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

	fmt.Print("Enter web probe ID: ")
	webProbeIDInput, _ := reader.ReadString('\n')
	webProbeIDInput = strings.TrimSpace(webProbeIDInput)
	webProbeID, err := strconv.Atoi(webProbeIDInput)
	if err != nil {
		log.Fatalf("[ERROR] Invalid web probe ID: %v\n", err)
	}

	fmt.Print("Enter cloud path probe ID: ")
	cloudPathProbeIDInput, _ := reader.ReadString('\n')
	cloudPathProbeIDInput = strings.TrimSpace(cloudPathProbeIDInput)
	cloudPathProbeID, err := strconv.Atoi(cloudPathProbeIDInput)
	if err != nil {
		log.Fatalf("[ERROR] Invalid cloud path probe ID: %v\n", err)
	}

	// Prompt for the session name
	fmt.Print("Enter a name for the deep trace session: ")
	sessionName, _ := reader.ReadString('\n')
	sessionName = strings.TrimSpace(sessionName)

	// Prompt for the session length in minutes
	fmt.Print("Enter session length in minutes: ")
	sessionLengthInput, _ := reader.ReadString('\n')
	sessionLengthInput = strings.TrimSpace(sessionLengthInput)
	sessionLength, err := strconv.Atoi(sessionLengthInput)
	if err != nil {
		log.Fatalf("[ERROR] Invalid session length: %v\n", err)
	}

	// Create a DeepTrace session
	payload := deeptrace.DeepTraceSessionPayload{
		SessionName:          sessionName,
		AppID:                appID,
		WebProbeID:           webProbeID,
		CloudPathProbeID:     cloudPathProbeID,
		SessionLengthMinutes: sessionLength,
		ProbeDevice:          true,
	}

	createdSession, resp, err := deeptrace.CreateDeepTraceSession(ctx, service, deviceID, payload)
	if err != nil {
		log.Fatalf("Error creating deep trace session: %v", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		log.Fatalf("Expected status code 200 or 201, got %d", resp.StatusCode)
	}

	traceID := createdSession.TraceID
	log.Printf("Created Deep Trace Session: %s\n", traceID)

	// Pause for 30 seconds
	log.Println("Pausing for 30 seconds to allow the session to start...")
	time.Sleep(30 * time.Second)

	// Get deep trace session again to update the status
	traceSessionResp, err := deeptrace.GetDeepTraceSession(ctx, service, deviceID, traceID)
	if err != nil {
		log.Fatalf("Error getting deep trace session: %v", err)
	}

	var updatedSession deeptrace.DeepTraceSession
	err = json.NewDecoder(traceSessionResp.Body).Decode(&updatedSession)
	if err != nil {
		log.Fatalf("Error decoding deep trace session response: %v", err)
	}

	if traceSessionResp.StatusCode != http.StatusOK {
		log.Fatalf("Expected status code 200, got %d", traceSessionResp.StatusCode)
	}

	// Display the deep trace session details
	displayTraceSession(&updatedSession)

	// Ask the user if they want to delete the deep trace session
	fmt.Print("Do you want to delete the deep trace session? (yes/no): ")
	deleteSession, _ := reader.ReadString('\n')
	deleteSession = strings.TrimSpace(deleteSession)

	if strings.ToLower(deleteSession) == "yes" {
		// Delete the deep trace session
		deleteResp, err := deeptrace.DeleteDeepTraceSession(ctx, service, deviceID, traceID)
		if err != nil {
			log.Fatalf("Error deleting deep trace session: %v", err)
		}

		if deleteResp.StatusCode != http.StatusOK && deleteResp.StatusCode != http.StatusNoContent {
			log.Fatalf("Expected status code 200 or 204, got %d", deleteResp.StatusCode)
		}

		log.Printf("Deleted deep trace session: %s\n", traceID)

		// Display the updated deep trace session details
		displayTraceSession(&updatedSession)
	}
}

func displayTraceSession(session *deeptrace.DeepTraceSession) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Trace ID", "Session Name", "Status", "Created At", "Started At", "Ended At"})

	formatTimestamp := func(ts int) string {
		if ts == 0 {
			return "N/A"
		}
		return time.Unix(int64(ts), 0).Format("2006-01-02")
	}

	table.Append([]string{
		session.TraceID,
		session.TraceDetails.SessionName,
		session.Status,
		formatTimestamp(session.CreatedAt),
		formatTimestamp(session.StartedAt),
		formatTimestamp(session.EndedAt),
	})
	table.Render()
}
