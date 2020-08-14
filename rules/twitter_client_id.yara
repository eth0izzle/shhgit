rule twitter_client_id
{
	meta:
		author = "Paul Price"
		description = "Finds Twitter Client ID"

	strings:
		$ = /twitter(.{0,20})?['\"][0-9a-z]{18,25}['\"]/

	condition:
		any of them
}