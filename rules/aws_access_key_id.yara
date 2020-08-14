rule aws_access_key_id
{
	meta:
		author = "Paul Price"
		description = "Finds AWS Access Key ID"

	strings:
		$ = /((\"|'|`)?(aws)?_?(access)_?(key)?_?(id)?(\"|'|`)?\\s{0,50}(:|=>|=)\\s{0,50}(\"|'|`)?(A3T[A-Z0-9]|AKIA|AGPA|AIDA|AROA|AIPA|ANPA|ANVA|ASIA)[A-Z0-9]{16}(\"|'|`)?)/

	condition:
		any of them
}