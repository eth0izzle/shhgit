rule salesforce_credentials_in_a_nodejs_project
{
	meta:
		author = "Paul Price"
		description = "Finds Salesforce credentials in a nodejs project"

	condition:
		filename contains "salesforce.js"
}