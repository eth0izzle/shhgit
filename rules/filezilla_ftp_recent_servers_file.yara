rule filezilla_ftp_recent_servers_file
{
	meta:
		author = "Paul Price"
		description = "Finds FileZilla FTP recent servers file"

	condition:
		filename contains "recentservers.xml"
}