package cmd

type rootCommandOptions struct {
	Broker   string
	Username string
	Password string

	Verbosity int
}

type sendLatestCommandOptions struct {
	Directory string
	Extension string
	Topic     string
}
