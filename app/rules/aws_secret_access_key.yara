rule aws_secret_access_key
{
	meta:
		author = "Paul Price"
		description = "Finds AWS Secret Access Key"

	strings:
		$ = /((\"|'|`)?(aws)?_?(secret)_?(access)?_?(key)?_?(id)?(\"|'|`)?\\s{0,50}(:|=>|=)\\s{0,50}(\"|'|`)?[A-Za-z0-9\/+=]{40}(\"|'|`)?)/

	condition:
		any of them
}