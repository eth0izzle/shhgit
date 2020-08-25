rule php_configuration_file
{
	meta:
		author = "Paul Price"
		description = "Finds PHP configuration file"

	condition:
		filename matches /config(\.inc)?\.php$/
}