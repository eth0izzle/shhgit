rule sauce_token
{
	meta:
		author = "Paul Price"
		description = "Finds Sauce Token"

	strings:
		$ = /sauce.{0,50}(\"|'|`)?[0-9a-f-]{36}(\"|'|`)?/

	condition:
		any of them
}