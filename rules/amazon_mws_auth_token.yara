rule amazon_mws_auth_token
{
	meta:
		author = "Paul Price"
		description = "Finds Amazon MWS Auth Token"

	strings:
		$ = /amzn\.mws\.[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}/

	condition:
		any of them
}