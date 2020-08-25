rule docker_configuration_file
{
	meta:
		author = "Paul Price"
		description = "Finds Docker configuration file"

	condition:
		filename matches /^\.?dockercfg$/
}