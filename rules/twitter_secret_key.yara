rule twitter_secret_key
{
	meta:
		author = "Paul Price"
		description = "Finds Twitter Secret Key"

	strings:
		$ = /twitter(.{0,20})?['\"][0-9a-z]{35,44}['\"]/

	condition:
		any of them
}