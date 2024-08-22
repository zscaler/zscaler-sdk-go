package zdx

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NewZDXCommand creates the base command for ZDX service
func NewZDXCommand(state *cmd.cliState) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "zdx",
		Short: "Interact with the Zscaler ZDX API",
		Long:  "ZDX command allows you to interact with Zscaler ZDX API to retrieve monitoring and diagnostics data.",
	}

	// Add subcommands
	cmd.AddCommand(NewListDevicesCommand(state))

	return cmd
}

// NewListDevicesCommand returns a command to list ZDX devices
func NewListDevicesCommand(state *cmd.cliState) *cobra.Command {
	return &cobra.Command{
		Use:   "list-devices",
		Short: "List all devices monitored by ZDX",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Ensure the ZDX client is initialized
			if state.ZDXApi == nil {
				if err := state.NewClient(); err != nil {
					return err
				}
			}

			// Call the ZDX API to list devices (example)
			devices, err := state.ZDXApi.ListDevices()
			if err != nil {
				return fmt.Errorf("error fetching devices: %w", err)
			}

			// Print the list of devices
			for _, device := range devices {
				fmt.Printf("Device: %s\n", device.Name)
			}

			return nil
		},
	}
}
