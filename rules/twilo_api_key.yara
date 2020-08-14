rule twilo_api_key
{
	meta:
		author = "Paul Price"
		description = "Finds Twilo API Key"

	strings:
		$ = /SK[0-9a-fA-F]{32}/

	condition:
		any of them
}