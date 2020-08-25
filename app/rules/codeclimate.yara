rule codeclimate
{
	meta:
		author = "Paul Price"
		description = "Finds CodeClimate"

	strings:
		$ = /codeclima.{0,50}(\"|'|`)?[0-9a-f]{64}(\"|'|`)?/

	condition:
		any of them
}