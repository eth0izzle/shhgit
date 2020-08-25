rule facebook_secret_key
{
	meta:
		author = "Paul Price"
		description = "Finds Facebook Secret Key"

	strings:
		$ = /(facebook|fb)(.{0,20})?['\"][0-9a-f]{32}['\"]/

	condition:
		any of them
}