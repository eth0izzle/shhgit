rule s3cmd_configuration_file
{
	meta:
		author = "Paul Price"
		description = "Finds S3cmd configuration file"

	condition:
		filename matches /^\.?s3cfg$/
}