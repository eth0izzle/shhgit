rule digitalocean_doctl_command_line_client_configuration_file
{
	meta:
		author = "Paul Price"
		description = "Finds DigitalOcean doctl command-line client configuration file"

	condition:
		filepath matches /doctl\/config.yaml$/
}