rule linkedin_client_id
{
	meta:
		author = "Paul Price"
		description = "Finds Linkedin Client ID"

	strings:
		$ = /linkedin(.{0,20})?['\"][0-9a-z]{12}['\"]/

	condition:
		any of them
}