rule stripe_api_key
{
	meta:
		author = "Paul Price"
		description = "Finds Stripe API key"

	strings:
		$ = /(r|s)k_[live|test]_[0-9a-zA-Z]{24}/

	condition:
		any of them
}