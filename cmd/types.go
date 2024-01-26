package cmd

type rootCommandOptions struct {
	Broker   string
	Username string
	Password string
}

type sendLatestCommandOptions struct {
	Directory string
	Extension string
	Topic     string
}
