rule square_access_token
{
	meta:
		author = "Paul Price"
		description = "Finds Square Access Token"

	strings:
		$ = /sq0atp-[0-9A-Za-z\-_]{22}/

	condition:
		any of them
}