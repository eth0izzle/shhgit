rule django_configuration_file
{
	meta:
		author = "Paul Price"
		description = "Finds Django configuration file"

	condition:
		filename contains "settings.py"
}