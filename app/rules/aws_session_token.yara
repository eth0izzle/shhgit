rule aws_session_token
{
	meta:
		author = "Paul Price"
		description = "Finds AWS Session Token"

	strings:
		$ = /((\"|'|`)?(aws)?_?(session)?_?(token)?(\"|'|`)?\\s{0,50}(:|=>|=)\\s{0,50}(\"|'|`)?[A-Za-z0-9\/+=]{16,}(\"|'|`)?)/

	condition:
		any of them
}