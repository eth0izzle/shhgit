rule cpanel_backup_proftpd_credentials_file
{
	meta:
		author = "Paul Price"
		description = "Finds cPanel backup ProFTPd credentials file"

	condition:
		filename contains "proftpdpasswd"
}