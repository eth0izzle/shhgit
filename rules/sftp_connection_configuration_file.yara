rule sftp_connection_configuration_file
{
	meta:
		author = "Paul Price"
		description = "Finds SFTP connection configuration file"

	condition:
		filename matches /^sftp-config(\.json)?$/
}