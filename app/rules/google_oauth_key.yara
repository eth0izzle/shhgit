rule google_oauth_key
{
	meta:
		author = "Paul Price"
		description = "Finds Google OAuth Key"

	strings:
		$ = /[0-9]+-[0-9A-Za-z_]{32}\.apps\.googleusercontent\.com/

	condition:
		any of them
}