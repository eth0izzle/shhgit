rule linkedin_secret_key
{
	meta:
		author = "Paul Price"
		description = "Finds LinkedIn Secret Key"

	strings:
		$ = /linkedin(.{0,20})?['\"][0-9a-z]{16}['\"]/

	condition:
		any of them
}