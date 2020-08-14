rule heroku_api_key
{
	meta:
		author = "Paul Price"
		description = "Finds Heroku API key"

	strings:
		$ = /heroku(.{0,20})?['"][0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}['"]/

	condition:
		any of them
}