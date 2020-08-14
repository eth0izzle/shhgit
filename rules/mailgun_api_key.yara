rule mailgun_api_key
{
	meta:
		author = "Paul Price"
		description = "Finds MailGun API Key"

	strings:
		$ = /key-[0-9a-zA-Z]{32}/

	condition:
		any of them
}