package cmd

type rootCommandOptions struct {
	Broker string
}

type sendLatestCommandOptions struct {
	Directory string
	Extension string
	Topic     string
}
