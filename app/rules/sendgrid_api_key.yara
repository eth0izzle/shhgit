rule sendgrid_api_key
{
	meta:
		author = "Paul Price"
		description = "Finds SendGrid API Key"

	strings:
		$ = /SG\.[0-9A-Za-z\-_]{22}\.[0-9A-Za-z\-_]{43}/

	condition:
		any of them
}