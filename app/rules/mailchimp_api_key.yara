rule mailchimp_api_key
{
	meta:
		author = "Paul Price"
		description = "Finds MailChimp API Key"

	strings:
		$ = /[0-9a-f]{32}-us[0-9]{12}/

	condition:
		any of them
}