package cmd

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/mannkind/mqtt-snapshot/mqtt"
	"github.com/spf13/cobra"
)

// Represents the ability to send the latest snapshot to mqtt
var sendLatestOpts = sendLatestCommandOptions{}
var sendLatestCmd = &cobra.Command{
	Use:   "send-latest",
	Short: "Send the latest snapshot to MQTT",
	Run: func(cmd *cobra.Command, args []string) {
		// Publish the image content
		c, err := mqtt.New(rootCmdOpts.Broker)
		if err != nil {
			fmt.Printf("Error creating initial MQTT client; %s\n", err)

			connect := time.NewTicker(7 * time.Second)
			for range connect.C {
				c, err = mqtt.New(rootCmdOpts.Broker)

				if err != nil {
					fmt.Printf("Error creating MQTT client; %s\n", err)
					continue
				}

				connect.Stop()
			}
		}

		ticker := time.NewTicker(7 * time.Second)
		for range ticker.C {
			// Fetch the latest image content
			contents, err := latestImageContent(sendLatestOpts.Directory, sendLatestOpts.Extension)
			if err != nil {
				fmt.Printf("Error fetching latest image content; %s\n", err)
				continue
			}

			token := c.Publish(sendLatestOpts.Topic, 0, true, contents)
			token.Wait()
			if err := token.Error(); err != nil {
				fmt.Printf("Error publishing image; %s\n", err)
				continue
			}
		}
	},
}

func latestImageContent(directory string, ext string) (string, error) {
	// Read all files in the directory
	// Filter out only the files with the correct extension
	files, err := readDir(directory, ext)
	if err != nil {
		return "", err
	}

	// Sort the files by modification time, descending
	sort.Slice(files, func(i, j int) bool {
		return files[i].ModTime().After(files[j].ModTime())
	})

	// Obtain the contents of the latest file
	file := fmt.Sprintf("%s/%s", directory, files[0].Name())
	contents, err := readFile(file)
	if err != nil {
		return "", err
	}

	// Return the contents
	return contents, nil
}

func readDir(directory string, ext string) ([]fs.FileInfo, error) {
	entries, err := os.ReadDir(sendLatestOpts.Directory)
	if err != nil {
		return nil, err
	}

	files := make([]fs.FileInfo, 0, len(entries))
	for _, entry := range entries {
		file, err := entry.Info()
		if err != nil {
			return nil, err
		}

		if !strings.HasSuffix(file.Name(), ext) {
			continue
		}

		files = append(files, file)
	}

	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, fmt.Errorf("no %s files found in %s", ext, directory)
	}

	return files, nil
}

func readFile(file string) (string, error) {
	// Open the file
	fileHandle, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer fileHandle.Close()

	// Read the contents
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
