rule github_hub_command_line_client_configuration_file
{
	meta:
		author = "Paul Price"
		description = "Finds GitHub Hub command-line client configuration file"

	condition:
		filepath matches /config\/hub$/
}