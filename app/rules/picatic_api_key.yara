rule picatic_api_key
{
	meta:
		author = "Paul Price"
		description = "Finds Picatic API key"

	strings:
		$ = /sk_[live|test]_[0-9a-z]{32}/

	condition:
		any of them
}