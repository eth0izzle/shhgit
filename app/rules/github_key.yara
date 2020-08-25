rule github_key
{
	meta:
		author = "Paul Price"
		description = "Finds Github Key"

	strings:
		$ = /github(.{0,20})?['\"][0-9a-zA-Z]{35,40}['\"]/

	condition:
		any of them
}