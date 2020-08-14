rule nuget_api_key
{
	meta:
		author = "Paul Price"
		description = "Finds NuGet API Key"

	strings:
		$ = /oy2[a-z0-9]{43}/

	condition:
		any of them
}