rule docker_registry_authentication_file
{
	meta:
		author = "Paul Price"
		description = "Finds Docker registry authentication file"

	condition:
		filepath matches /\.?docker[\\\/]config.json$/
}