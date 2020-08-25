rule aws_account_id
{
	meta:
		author = "Paul Price"
		description = "Finds AWS Account ID"

	strings:
		$ = /((\"|'|`)?(aws)?_?(account)_?(id)?(\"|'|`)?\\s{0,50}(:|=>|=)\\s{0,50}(\"|'|`)?[0-9]{4}-?[0-9]{4}-?[0-9]{4}(\"|'|`)?)/

	condition:
		any of them
}