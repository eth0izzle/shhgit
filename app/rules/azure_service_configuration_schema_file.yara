rule azure_service_configuration_schema_file
{
	meta:
		author = "Paul Price"
		description = "Finds Azure service configuration schema file"

	condition:
		extension contains ".cscfg"
}