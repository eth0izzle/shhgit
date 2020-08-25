rule environment_configuration_file
{
	meta:
		author = "Paul Price"
		description = "Finds Environment configuration file"

	condition:
		filename matches /^\.?env$/
}