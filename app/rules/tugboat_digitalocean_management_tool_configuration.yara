rule tugboat_digitalocean_management_tool_configuration
{
	meta:
		author = "Paul Price"
		description = "Finds Tugboat DigitalOcean management tool configuration"

	condition:
		filename matches /^\.?tugboat$/
}