rule paypal_braintree_access_token
{
	meta:
		author = "Paul Price"
		description = "Finds PayPal/Braintree Access Token"

	strings:
		$ = /access_token\$production\$[0-9a-z]{16}\$[0-9a-f]{32}/

	condition:
		any of them
}