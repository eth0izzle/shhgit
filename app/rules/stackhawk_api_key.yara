rule stackhawk_api_key
{
	meta:
		author = "Paul Price"
		description = "Finds StackHawk API Key"

	strings:
		$ = /hawk\.[0-9A-Za-z\-_]{20}\.[0-9A-Za-z\-_]{20}/

	condition:
		any of them
}