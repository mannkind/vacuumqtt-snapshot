package cmd

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"sort"
	"strings"

	"github.com/mannkind/mqtt-snapshot/mqtt"
	"github.com/spf13/cobra"
)

// Represents the ability to send the latest snapshot to mqtt
var sendLatestOpts = sendLatestCommandOptions{}
var sendLatestCmd = &cobra.Command{
	Use:   "send-latest",
	Short: "Send the latest snapshot to MQTT",
	Run: func(cmd *cobra.Command, args []string) {
		contents, err := fetchLatestImageContent(sendLatestOpts.Directory, sendLatestOpts.Extension)
		if err != nil {
			fmt.Printf("Error fetching latest image content; %s\n", err)
			os.Exit(1)
			return
		}

		c, err := mqtt.New(rootCmdOpts.Broker)
		if err != nil {
			fmt.Printf("Error creating MQTT client; %s\n", err)
			os.Exit(1)
			return
		}

		token := c.Publish(sendLatestOpts.Topic, 0, false, contents)
		token.Wait()
		if err := token.Error(); err != nil {
			fmt.Printf("Error publishing image %s\n", err)
			os.Exit(1)
			return
		}
	},
}

func fetchLatestImageContent(directory string, ext string) (string, error) {
	entries, err := os.ReadDir(sendLatestOpts.Directory)
	if err != nil {
		return "", err
	}

	files := make([]fs.FileInfo, 0, len(entries))
	for _, entry := range entries {
		file, err := entry.Info()
		if err != nil {
			return "", err
		}

		if !strings.HasSuffix(file.Name(), ext) {
			continue
		}

		files = append(files, file)
	}

	if err != nil {
		return "", err
	}

	if len(files) == 0 {
		return "", err
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].ModTime().After(files[j].ModTime())
	})

	latestFilename := fmt.Sprintf("%s/%s", directory, files[0].Name())

	fileHandle, err := os.Open(latestFilename)
	if err != nil {
		return "", err
	}
	defer fileHandle.Close()

	contents, err := io.ReadAll(fileHandle)
	if err != nil {
		return "", err
	}

	return string(contents), nil
}

func init() {
	rootCmd.AddCommand(sendLatestCmd)
	sendLatestCmd.Flags().StringVar(&sendLatestOpts.Directory, "directory", "/data/ai_offline_collection", "The directory containing snapshots")
	sendLatestCmd.Flags().StringVar(&sendLatestOpts.Extension, "extension", ".jpg", "The extension of the snapshots")
	sendLatestCmd.Flags().StringVar(&sendLatestOpts.Topic, "topic", "valetudo/DraftyKnottyWasp/Camera/image-data-hass", "The topic to publish the image on")
}
