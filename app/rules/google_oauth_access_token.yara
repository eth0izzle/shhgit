rule google_oauth_access_token
{
	meta:
		author = "Paul Price"
		description = "Finds Google OAuth Access Token"

	strings:
		$ = /ya29\\.[0-9A-Za-z\\-_]+/

	condition:
		any of them
}