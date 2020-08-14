rule ventrilo_server_configuration_file
{
	meta:
		author = "Paul Price"
		description = "Finds Ventrilo server configuration file"

	condition:
		filename contains "ventrilo_srv.ini"
}