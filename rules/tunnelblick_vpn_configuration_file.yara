rule tunnelblick_vpn_configuration_file
{
	meta:
		author = "Paul Price"
		description = "Finds Tunnelblick VPN configuration file"

	condition:
		extension contains ".tblk"
}