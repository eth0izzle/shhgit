rule created_by_vscode_sftp_for_vscode_contains_sftp_ssh_server_details_and_credentials
{
	meta:
		author = "Paul Price"
		description = "Finds Created by vscode-sftp for VSCode, contains SFTP/SSH server details and credentials"

	condition:
		filepath matches /\.?vscode[\\\/]sftp.json$/
}