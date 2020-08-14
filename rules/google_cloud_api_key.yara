rule google_cloud_api_key
{
	meta:
		author = "Paul Price"
		description = "Finds Google Cloud API Key"

	strings:
		$ = /AIza[0-9A-Za-z\\-_]{35}/

	condition:
		any of them
}