package cmd

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/mannkind/vacuumqtt-snapshot/logging"
	"github.com/mannkind/vacuumqtt-snapshot/mqtt"
	"github.com/spf13/cobra"
)

// Represents the ability to send the latest snapshot to mqtt
var sendLatestOpts = sendLatestCommandOptions{}
var sendLatestCmd = &cobra.Command{
	Use:   "send-latest",
	Short: "Send the latest snapshot to MQTT",
	Run: func(cmd *cobra.Command, args []string) {
		log := logging.New(rootCmdOpts.Verbosity)

		// Publish the image content
		c, err := mqtt.New(rootCmdOpts.Broker, rootCmdOpts.Username, rootCmdOpts.Password)
		if err != nil {
			log.Error(err, "Error creating initial MQTT client", "broker", rootCmdOpts.Broker, "username", rootCmdOpts.Username, "password", "********")

			connect := time.NewTicker(7 * time.Second)
			for range connect.C {
				c, err = mqtt.New(rootCmdOpts.Broker, rootCmdOpts.Username, rootCmdOpts.Password)

				if err != nil {
					log.Error(err, "Error creating MQTT client", "broker", rootCmdOpts.Broker, "username", rootCmdOpts.Username, "password", "********")
					continue
				}

				connect.Stop()
				break
			}
		}

		lastImage := ""
		ticker := time.NewTicker(7 * time.Second)
		for range ticker.C {
			// Fetch the latest image content
			file, contents, err := latestImageContent(sendLatestOpts.Directory, sendLatestOpts.Extension, lastImage)
			if err != nil {
				log.Error(err, "Error fetching latest image content")
				continue
			}

			// Don't publish duplicate images
			if lastImage == file {
				log.Info("Won't publish duplicate image", "file", file)
				continue
			}

			// Publish the image content
			token := c.Publish(sendLatestOpts.Topic, 0, true, contents)
			token.Wait()
			if err := token.Error(); err != nil {
				log.Error(err, "Error publishing image", "topic", sendLatestOpts.Topic, "file", file)
				continue
			}

			lastImage = file
		}
	},
}

func latestImageContent(directory string, ext string, last string) (string, string, error) {
	// Read all files in the directory
	// Filter out only the files with the correct extension
	files, err := readDir(directory, ext)
	if err != nil {
		return "", "", err
	}

	// Sort the files by modification time, descending
	sort.Slice(files, func(i, j int) bool {
		return files[i].ModTime().After(files[j].ModTime())
	})

	// Obtain the contents of the latest file
	file := fmt.Sprintf("%s/%s", directory, files[0].Name())
	contents, err := readFile(file)
	if err != nil {
		return "", "", err
	}

	// Return the contents
	return file, contents, nil
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

		if filepath.Ext(file.Name()) != ext {
			continue
		}

		if file.Size() == 0 {
			continue
		}

		files = append(files, file)
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
