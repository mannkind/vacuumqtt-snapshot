package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// Represents the base command when called without any subcommands
var rootCmdOpts = rootCommandOptions{}
var rootCmd = &cobra.Command{
	Use:   "vacuumqtt-snapshot",
	Short: "The vacuumqtt-snapshot tool",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&rootCmdOpts.Broker, "broker", "mqtt.lan:1883", "The host:port of the broker")
	rootCmd.PersistentFlags().StringVar(&rootCmdOpts.Username, "username", "", "The username for the broker")
	rootCmd.PersistentFlags().StringVar(&rootCmdOpts.Password, "password", "", "The password for the broker")

	rootCmd.PersistentFlags().IntVar(&rootCmdOpts.Verbosity, "verbosity", 1, "The verbosity for logging")

}
