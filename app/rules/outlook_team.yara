rule outlook_team
{
	meta:
		author = "Paul Price"
		description = "Finds Outlook team"

	strings:
		$ = /(https\\:\/\/outlook\\.office.com\/webhook\/[0-9a-f-]{36}\\@)/

	condition:
		any of them
}