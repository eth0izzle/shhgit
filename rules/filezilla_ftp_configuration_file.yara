rule filezilla_ftp_configuration_file
{
	meta:
		author = "Paul Price"
		description = "Finds FileZilla FTP configuration file"

	condition:
		filename contains "filezilla.xml"
}