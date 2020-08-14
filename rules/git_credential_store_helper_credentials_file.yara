rule git_credential_store_helper_credentials_file
{
	meta:
		author = "Paul Price"
		description = "Finds git-credential-store helper credentials file"

	condition:
		filename matches /^\.?git-credentials$/
}