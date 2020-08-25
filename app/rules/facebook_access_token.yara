rule facebook_access_token
{
	meta:
		author = "Paul Price"
		description = "Finds Facebook access token"

	strings:
		$ = /EAACEdEose0cBA[0-9A-Za-z]+/

	condition:
		any of them
}