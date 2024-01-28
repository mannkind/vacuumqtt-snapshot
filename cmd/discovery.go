package cmd

import (
	"os"

	"github.com/mannkind/vacuumqtt-snapshot/logging"
	"github.com/mannkind/vacuumqtt-snapshot/mqtt"
	"github.com/spf13/cobra"
)

// Represents the ability to send a discovery message to mqtt
var discoveryCmd = &cobra.Command{
	Use:    "discovery",
	Short:  "Send the discovery message to MQTT",
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		log := logging.New(rootCmdOpts.Verbosity)

		// Publish the discovery message
		c, err := mqtt.New(rootCmdOpts.Broker, rootCmdOpts.Username, rootCmdOpts.Password)
		if err != nil {
			log.Error(err, "Error creating initial MQTT client", "broker", rootCmdOpts.Broker, "username", rootCmdOpts.Username, "password", "********")
			os.Exit(1)
			return
		}

		// Setup the discovery topic/content
		topic := "homeassistant/camera/DraftyKnottyWasp/DraftyKnottyWasp_camera_Image/config"
		contents := `{
			"topic":"valetudo/DraftyKnottyWasp/Camera/image-data-hass",
			"name":"Camera Image",
			"object_id":"valetudo_DraftyKnottyWasp_camera image",
			"unique_id":"DraftyKnottyWasp_camera_Image",
			"availability_topic":"valetudo/DraftyKnottyWasp/$state",
			"payload_available":"ready",
			"payload_not_available":"lost",
			"device":{"manufacturer":"Dreame","model":"L10S Ultra","name":"Valetudo L10S Ultra DraftyKnottyWasp","identifiers":["DraftyKnottyWasp"],"sw_version":"Valetudo 2023.12.0","configuration_url":"http://valetudo-draftyknottywasp.local"}
		}`

		token := c.Publish(topic, 0, true, contents)
		token.Wait()
		if err := token.Error(); err != nil {
			log.Error(err, "Error publishing discovery", "topic", topic)
			os.Exit(1)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(discoveryCmd)
}
