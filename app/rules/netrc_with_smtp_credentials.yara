rule netrc_with_smtp_credentials
{
	meta:
		author = "Paul Price"
		description = "Finds netrc with SMTP credentials"

	condition:
		extension contains ".netrc"
}