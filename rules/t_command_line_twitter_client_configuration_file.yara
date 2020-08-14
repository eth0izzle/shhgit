rule t_command_line_twitter_client_configuration_file
{
	meta:
		author = "Paul Price"
		description = "Finds T command-line Twitter client configuration file"

	condition:
		filename matches /^\.?trc$/
}