rule little_snitch_firewall_configuration_file
{
	meta:
		author = "Paul Price"
		description = "Finds Little Snitch firewall configuration file"

	condition:
		filename contains "configuration.user.xpl"
}