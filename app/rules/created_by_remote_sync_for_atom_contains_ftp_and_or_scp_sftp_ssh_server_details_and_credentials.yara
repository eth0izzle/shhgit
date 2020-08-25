rule created_by_remote_sync_for_atom_contains_ftp_and_or_scp_sftp_ssh_server_details_and_credentials
{
	meta:
		author = "Paul Price"
		description = "Finds Created by remote-sync for Atom, contains FTP and/or SCP/SFTP/SSH server details and credentials"

	condition:
		filename matches /.remote-sync.json$/
}