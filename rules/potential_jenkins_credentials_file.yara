rule potential_jenkins_credentials_file
{
	meta:
		author = "Paul Price"
		description = "Finds Potential Jenkins credentials file"

	condition:
		filename contains "credentials.xml"
}