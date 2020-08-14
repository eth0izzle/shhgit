rule ssh_configuration_file
{
	meta:
		author = "Paul Price"
		description = "Finds SSH configuration file"

	condition:
		filepath matches /\.?ssh\/config$/
}