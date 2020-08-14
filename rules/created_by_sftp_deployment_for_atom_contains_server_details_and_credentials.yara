rule created_by_sftp_deployment_for_atom_contains_server_details_and_credentials
{
	meta:
		author = "Paul Price"
		description = "Finds Created by sftp-deployment for Atom, contains server details and credentials"

	condition:
		filename matches /.ftpconfig$/
}