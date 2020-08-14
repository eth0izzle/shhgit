rule square_oauth_secret
{
	meta:
		author = "Paul Price"
		description = "Finds Square OAuth Secret"

	strings:
		$ = /sq0csp-[0-9A-Za-z\-_]{43}/

	condition:
		any of them
}