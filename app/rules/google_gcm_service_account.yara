rule google_gcm_service_account
{
	meta:
		author = "Paul Price"
		description = "Finds Google (GCM) Service account"

	strings:
		$ = /((\"|'|`)?type(\"|'|`)?\\s{0,50}(:|=>|=)\\s{0,50}(\"|'|`)?service_account(\"|'|`)?,?)/

	condition:
		any of them
}