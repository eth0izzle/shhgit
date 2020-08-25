rule openvpn_client_configuration_file
{
	meta:
		author = "Paul Price"
		description = "Finds OpenVPN client configuration file"

	condition:
		extension contains ".ovpn"
}